#include <obs-module.h>
#include <obs.h>
#include <obs-source.h>
#include <graphics/image-file.h>
#include <util/platform.h>
#include <util/dstr.h>
#include <sys/stat.h>

using namespace std;

#define blog(log_level, format, ...)                    \
	blog(log_level, "[image_source: '%s'] " format, \
	     obs_source_get_name(context->source), ##__VA_ARGS__)

#define debug(format, ...) blog(LOG_DEBUG, format, ##__VA_ARGS__)
#define info(format, ...) blog(LOG_INFO, format, ##__VA_ARGS__)
#define warn(format, ...) blog(LOG_WARNING, format, ##__VA_ARGS__)

struct image_source {
	obs_source_t* source;

	char* file;
	bool persistent;
	bool linear_alpha;
	time_t file_timestamp;
	float update_time_elapsed;
	uint64_t last_time;
	bool active;

	gs_image_file3_t if3;
};

static time_t get_modified_timestamp(const char* filename)
{
	struct stat stats;
	if (os_stat(filename, &stats) != 0)
		return -1;
	return stats.st_mtime;
}

static const char* image_source_get_name(void* unused)
{
	UNUSED_PARAMETER(unused);
	return obs_module_text("ImageInput");
}

static void image_source_load(struct image_source* context)
{
	char* file = context->file;

	obs_enter_graphics();
	gs_image_file3_free(&context->if3);
	obs_leave_graphics();

	if (file && *file) {
		debug("loading texture '%s'", file);
		context->file_timestamp = get_modified_timestamp(file);
		gs_image_file3_init(&context->if3, file,
			context->linear_alpha
			? GS_IMAGE_ALPHA_PREMULTIPLY_SRGB
			: GS_IMAGE_ALPHA_PREMULTIPLY);
		context->update_time_elapsed = 0;

		obs_enter_graphics();
		gs_image_file3_init_texture(&context->if3);
		obs_leave_graphics();

		if (!context->if3.image2.image.loaded)
			warn("failed to load texture '%s'", file);
	}
}

static void image_source_unload(struct image_source* context)
{
	obs_enter_graphics();
	gs_image_file3_free(&context->if3);
	obs_leave_graphics();
}

static void image_source_update(void* data, obs_data_t* settings)
{
	struct image_source* context = (image_source*)data;
	const char* file = obs_data_get_string(settings, "file");
	const bool unload = obs_data_get_bool(settings, "unload");
	const bool linear_alpha = obs_data_get_bool(settings, "linear_alpha");

	if (context->file)
		bfree(context->file);
	context->file = bstrdup(file);
	context->persistent = !unload;
	context->linear_alpha = linear_alpha;

	/* Load the image if the source is persistent or showing */
	if (context->persistent || obs_source_showing(context->source))
		image_source_load((image_source*)data);
	else
		image_source_unload((image_source*)data);
}

static void image_source_defaults(obs_data_t* settings)
{
	obs_data_set_default_bool(settings, "unload", false);
	obs_data_set_default_bool(settings, "linear_alpha", false);
}

static void image_source_show(void* data)
{
	struct image_source* context = (image_source*)data;

	if (!context->persistent)
		image_source_load(context);
}

static void image_source_hide(void* data)
{
	struct image_source* context = (image_source*)data;

	if (!context->persistent)
		image_source_unload(context);
}

static void* image_source_create(obs_data_t* settings, obs_source_t* source)
{
	struct image_source* context = (image_source*)bzalloc(sizeof(struct image_source));
	context->source = source;

	image_source_update(context, settings);
	return context;
}

static void image_source_destroy(void* data)
{
	struct image_source* context = (image_source*)data;

	image_source_unload(context);

	if (context->file)
		bfree(context->file);
	bfree(context);
}

static uint32_t image_source_getwidth(void* data)
{
	struct image_source* context = (image_source*)data;
	return context->if3.image2.image.cx;
}

static uint32_t image_source_getheight(void* data)
{
	struct image_source* context = (image_source*)data;
	return context->if3.image2.image.cy;
}

static void image_source_render(void* data, gs_effect_t* effect)
{
	struct image_source* context = (image_source*)data;

	if (!context->if3.image2.image.texture)
		return;

	const bool previous = gs_framebuffer_srgb_enabled();
	gs_enable_framebuffer_srgb(true);

	gs_blend_state_push();
	gs_blend_function(GS_BLEND_ONE, GS_BLEND_INVSRCALPHA);

	gs_eparam_t* const param = gs_effect_get_param_by_name(effect, "image");
	gs_effect_set_texture_srgb(param, context->if3.image2.image.texture);

	gs_draw_sprite(context->if3.image2.image.texture, 0,
		context->if3.image2.image.cx,
		context->if3.image2.image.cy);

	gs_blend_state_pop();

	gs_enable_framebuffer_srgb(previous);
}

static void image_source_tick(void* data, float seconds)
{
	struct image_source* context = (image_source*)data;
	uint64_t frame_time = obs_get_video_frame_time();

	context->update_time_elapsed += seconds;

	if (obs_source_showing(context->source)) {
		if (context->update_time_elapsed >= 1.0f) {
			time_t t = get_modified_timestamp(context->file);
			context->update_time_elapsed = 0.0f;

			if (context->file_timestamp != t) {
				image_source_load(context);
			}
		}
	}

	if (obs_source_active(context->source)) {
		if (!context->active) {
			if (context->if3.image2.image.is_animated_gif)
				context->last_time = frame_time;
			context->active = true;
		}

	}
	else {
		if (context->active) {
			if (context->if3.image2.image.is_animated_gif) {
				context->if3.image2.image.cur_frame = 0;
				context->if3.image2.image.cur_loop = 0;
				context->if3.image2.image.cur_time = 0;

				obs_enter_graphics();
				gs_image_file3_update_texture(&context->if3);
				obs_leave_graphics();
			}

			context->active = false;
		}

		return;
	}

	if (context->last_time && context->if3.image2.image.is_animated_gif) {
		uint64_t elapsed = frame_time - context->last_time;
		bool updated = gs_image_file3_tick(&context->if3, elapsed);

		if (updated) {
			obs_enter_graphics();
			gs_image_file3_update_texture(&context->if3);
			obs_leave_graphics();
		}
	}

	context->last_time = frame_time;
}

static const char* image_filter =
"All formats (*.bmp *.tga *.png *.jpeg *.jpg *.gif *.psd *.webp);;"
"BMP Files (*.bmp);;"
"Targa Files (*.tga);;"
"PNG Files (*.png);;"
"JPEG Files (*.jpeg *.jpg);;"
"GIF Files (*.gif);;"
"PSD Files (*.psd);;"
"WebP Files (*.webp);;"
"All Files (*.*)";

static obs_properties_t* image_source_properties(void* data)
{
	struct image_source* s = (image_source*)data;
	struct dstr path = { 0 };

	obs_properties_t* props = obs_properties_create();

	if (s && s->file && *s->file) {
		const char* slash;

		dstr_copy(&path, s->file);
		dstr_replace(&path, "\\", "/");
		slash = strrchr(path.array, '/');
		if (slash)
			dstr_resize(&path, slash - path.array + 1);
	}

	obs_properties_add_path(props, "file", obs_module_text("File"),
		OBS_PATH_FILE, image_filter, path.array);
	obs_properties_add_bool(props, "unload",
		obs_module_text("UnloadWhenNotShowing"));
	obs_properties_add_bool(props, "linear_alpha",
		obs_module_text("LinearAlpha"));
	dstr_free(&path);

	return props;
}

uint64_t image_source_get_memory_usage(void* data)
{
	struct image_source* s = (image_source*)data;
	return s->if3.image2.mem_usage;
}

static void missing_file_callback(void* src, const char* new_path, void* data)
{
	struct image_source* s = (image_source*)src;

	obs_source_t* source = s->source;
	obs_data_t* settings = obs_source_get_settings(source);
	obs_data_set_string(settings, "file", new_path);
	obs_source_update(source, settings);
	obs_data_release(settings);

	UNUSED_PARAMETER(data);
}

static obs_missing_files_t* image_source_missingfiles(void* data)
{
	struct image_source* s = (image_source*)data;
	obs_missing_files_t* files = obs_missing_files_create();

	if (strcmp(s->file, "") != 0) {
		if (!os_file_exists(s->file)) {
			obs_missing_file_t* file = obs_missing_file_create(
				s->file, missing_file_callback,
				OBS_MISSING_FILE_SOURCE, s->source, NULL);

			obs_missing_files_add_file(files, file);
		}
	}

	return files;
}

struct obs_source_info create_source_info()
{
	struct obs_source_info image_source_info = {};
	image_source_info.id = "image_source";
	image_source_info.type = OBS_SOURCE_TYPE_INPUT;
	image_source_info.output_flags = OBS_SOURCE_VIDEO | OBS_SOURCE_SRGB | OBS_SOURCE_CUSTOM_DRAW;
	image_source_info.get_name = image_source_get_name;
	image_source_info.create = image_source_create;
	image_source_info.destroy = image_source_destroy;
	image_source_info.update = image_source_update;
	image_source_info.get_defaults = image_source_defaults;
	image_source_info.show = image_source_show;
	image_source_info.hide = image_source_hide;
	image_source_info.get_width = image_source_getwidth;
	image_source_info.get_height = image_source_getheight;
	image_source_info.video_render = image_source_render;
	image_source_info.video_tick = image_source_tick;
	image_source_info.missing_files = image_source_missingfiles;
	image_source_info.get_properties = image_source_properties;
	image_source_info.icon_type = OBS_ICON_TYPE_IMAGE;
	return image_source_info;
};

OBS_DECLARE_MODULE()
OBS_MODULE_USE_DEFAULT_LOCALE("CMakeProject2", "en-US")
MODULE_EXPORT const char* obs_module_description(void)
{
	return "Image/color/slideshow sources";
}

bool obs_module_load(void)
{
	obs_startup("en-US", NULL, NULL);
	obs_source_info img_source_info = create_source_info();
	obs_register_source(&img_source_info);
	return true;
}

//void /*obs_source_video_render_pure(obs_source_t* source)
//{
//	source->info.video_render(source->context.data, NULL);
//}*/

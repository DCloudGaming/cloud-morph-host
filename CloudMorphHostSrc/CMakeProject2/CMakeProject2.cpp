// CMakeProject2.cpp : Defines the entry point for the application.
//
#include <windows.h>
#include <stdlib.h>
#include <string.h>
#include <tchar.h>
#include "CMakeProject2.h"
#include "MyPlugin.h"
#include "MyPlugin2.h"

using namespace std;

class SourceContext {

public:
	obs_source_t* source;

	inline SourceContext(obs_source_t* source) : source(source) {}
	inline ~SourceContext() { obs_source_release(source); }
	inline operator obs_source_t* () { return source; }
};

class SceneContext {
	obs_scene_t* scene;

public:
	inline SceneContext(obs_scene_t* scene) : scene(scene) {}
	inline ~SceneContext() { obs_scene_release(scene); }
	inline operator obs_scene_t* () { return scene; }
};

int main()
{
	obs_module_load();
	//obs_load_all_modules();
	//obs_module_load_2();
	SourceContext source = obs_source_create(
		"image_source", "some randon source", NULL, nullptr);

	//obs_source_video_render(source);
	obs_source_video_render(source);

	SceneContext scene = obs_scene_create("test scene");
	if (!scene)
		throw "Couldn't create scene";

	cout << "Hello CMaksafasfgge." << endl;
	MessageBox(NULL,
		_T("Call to Create Window failed!"),
		_T("Windows Desktop Guided Tour"),
		NULL);
	return 0;
}

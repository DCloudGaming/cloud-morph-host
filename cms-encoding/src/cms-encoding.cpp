// cms-encoding.cpp : Defines the entry point for the application.
//
#include <iostream>
// #include <boost/stacktrace.hpp>
#include <stdexcept>
#include <exception>

#include <libavformat/avformat.h>
#include <libavcodec/avcodec.h>
#include <ScreenRecorder.h>

using namespace std;

void print_exception(const std::exception& e, int level =  0)
{
    std::cerr << std::string(level, ' ') << "exception: " << e.what() << '\n';
    try {
        std::rethrow_if_nested(e);
    } catch(const std::exception& e) {
        print_exception(e, level+1);
    } catch(...) {}
}

int main()
{
	try {
		char* capture_method = "gdigrab";
		char* video_stream_url = "desktop";
		// char* capture_method = "dshow";
		// char* video_stream_url = "video=screen-capture-recorder";
		char* mp4_out_file = "./outputHieu.mp4";

		ScreenRecorder screen_record;
		screen_record.openCamera(capture_method, video_stream_url);
		screen_record.init_outputfile(mp4_out_file);
		//screen_record.init_outstream();
		screen_record.CaptureVideoFrames();
		return 0;
    } catch(const std::exception& e) {
        print_exception(e);
    }
}

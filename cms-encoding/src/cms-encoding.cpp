// cms-encoding.cpp : Defines the entry point for the application.
//
#include <iostream>

#include <libavformat/avformat.h>
#include <libavcodec/avcodec.h>
#include <ScreenRecorder.h>

using namespace std;

int main()
{
	std::cout << "Hello CMake21." << std::endl;
	ScreenRecorder screen_record;
	screen_record.openCamera();
	screen_record.init_outputfile();
	screen_record.CaptureVideoFrames();
	return 0;
}

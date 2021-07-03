# CMakeList.txt : CMake project for CMakeProject2, include source and define
# project specific logic here.
#
cmake_minimum_required (VERSION 3.8)

project(CMakeProject2)
set (CMAKE_CXX_STANDARD 11)
set(CMAKE_CXX_FLAGS "-D__STDC_CONSTANT_MACROS")
set(CMAKE_MODULE_PATH ${CMAKE_CURRENT_SOURCE_DIR}/../CMakeModules)
#set(FFMPEGEXAMPLE_PROTO_VER 0)
#set(FFMPEGEXAMPLE_VERSION "${CPACK_PACKAGE_VERSION_MAJOR}.\
#${CPACK_PACKAGE_VERSION_MINOR}.${CPACK_PACKAGE_VERSION_PATCH}")
#find_package(LIBAV REQUIRED)
#
#if(LIBAV_FOUND)
#  message("Found FFMPEG/LibAV libraries")
#  include_directories(${LIBAV_INCLUDE_DIR})
#  target_link_libraries(ffmpeg-example-bin ${LIBAV_LIBRARIES})
#else()
#  message("Can't find libavcodec, libavformat or libavutil. Add them!")
#endif()

find_path(AVCODEC_INCLUDE_DIR libavcodec/avcodec.h)
find_library(AVCODEC_LIBRARY avcodec)

find_path(AVFORMAT_INCLUDE_DIR libavformat/avformat.h)
find_library(AVFORMAT_LIBRARY avformat)

find_path(AVUTIL_INCLUDE_DIR libavutil/avutil.h)
find_library(AVUTIL_LIBRARY avutil)

find_path(AVDEVICE_INCLUDE_DIR libavdevice/avdevice.h)
find_library(AVDEVICE_LIBRARY avdevice)

set(myplugin_SOURCES
	MyPlugin.cpp
	MySource.cpp
)
set(myplugin_HEADERS
	MyPlugin.h
	MySource.h
	 MyPlugin2.h)
add_library(my-plugin ${myplugin_SOURCES} ${myplugin_HEADERS})
target_link_libraries(my-plugin libobs w32-pthreads avlformat avcodec)
# Add source to this project's executable.
add_executable (my-server "CMakeProject2.cpp" "CMakeProject2.h")
target_include_directories(my-server PRIVATE ${AVCODEC_INCLUDE_DIR} ${AVFORMAT_INCLUDE_DIR} ${AVUTIL_INCLUDE_DIR} ${AVDEVICE_INCLUDE_DIR})
target_link_libraries(my-server my-plugin libobs ${AVCODEC_LIBRARY} ${AVFORMAT_LIBRARY} ${AVUTIL_LIBRARY} ${AVDEVICE_LIBRARY})
install_obs_plugin_with_data(my-plugin data)

set(my-gui_SOURCES GUI.cpp)
add_executable (my-gui WIN32 ${my-gui_SOURCES})
target_link_libraries(my-gui libobs)
set_target_properties(my-gui PROPERTIES FOLDER "tests and examples")
define_graphic_modules(my-gui)



#add_executable (CMakeProject2 "ClientApp2.cpp" "ClientApp2.h")

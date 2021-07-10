# CMakeList.txt : CMake project for cms-encoding, include source and define
# project specific logic here.
#
cmake_minimum_required (VERSION 3.8)

project(src)
include_directories(${CMAKE_CURRENT_SOURCE_DIR})

set (CMAKE_CXX_STANDARD 11)
set(CMAKE_CXX_FLAGS "-D__STDC_CONSTANT_MACROS")
#set(CMAKE_MODULE_PATH ${CMAKE_CURRENT_SOURCE_DIR}\\..\\CMakeModules)
set (DCMAKE_TOOLCHAIN_FILE "D:\\repos\\vcpkg\\scripts\\buildsystems\\vcpkg.cmake")
set (CMAKE_TOOLCHAIN_FILE "D:\\repos\\vcpkg\\scripts\\buildsystems\\vcpkg.cmake")
#set (VCPKG_TARGET_TRIPLET "x86-windows")

FIND_PACKAGE(Boost 1.61 COMPONENTS filesystem program_options iostreams timer REQUIRED)
INCLUDE_DIRECTORIES ( ${Boost_INCLUDE_DIR} )

find_path(AVCODEC_INCLUDE_DIR libavcodec/avcodec.h)
find_library(AVCODEC_LIBRARY avcodec)
#
find_path(AVFORMAT_INCLUDE_DIR libavformat/avformat.h)
find_library(AVFORMAT_LIBRARY avformat)
#
find_path(AVUTIL_INCLUDE_DIR libavutil/avutil.h)
find_library(AVUTIL_LIBRARY avutil)
#
find_path(AVDEVICE_INCLUDE_DIR libavdevice/avdevice.h)
find_library(AVDEVICE_LIBRARY avdevice)

find_path(AVFILTER_INCLUDE_DIR libavfilter/avfilter.h)
find_library(AVFILTER_LIBRARY avfilter)

find_path(SWSCALE_INCLUDE_DIR libswscale/swscale.h)
find_library(SWSCALE_LIBRARY swscale)

set(cmsEncoding_SOURCES 
	${CMAKE_CURRENT_SOURCE_DIR}/ScreenRecorder.cpp
)

set(cmsEncoding_HEADERS
	${CMAKE_CURRENT_SOURCE_DIR}/ScreenRecorder.h
)
add_library(cmsEncodingLib  ${cmsEncoding_SOURCES} ${cmsEncoding_HEADERS} )
target_include_directories(cmsEncodingLib PUBLIC ${SWSCALE_INCLUDE_DIR} ${AVCODEC_INCLUDE_DIR} ${AVFORMAT_INCLUDE_DIR} ${AVUTIL_INCLUDE_DIR} ${AVDEVICE_INCLUDE_DIR} ${AVFILTER_INCLUDE_DIR})
target_link_libraries(cmsEncodingLib PUBLIC ${SWSCALE_LIBRARY} ${AVCODEC_LIBRARY} ${AVFORMAT_LIBRARY} ${AVUTIL_LIBRARY} ${AVDEVICE_LIBRARY} ${AVFILTER_LIBRARY})

# Add source to this project's executable.
add_executable (cms-encoding cms-encoding.cpp ${cmsEncoding_HEADERS})
target_include_directories(cms-encoding PUBLIC ${SWSCALE_INCLUDE_DIR} ${AVCODEC_INCLUDE_DIR} ${AVFORMAT_INCLUDE_DIR} ${AVUTIL_INCLUDE_DIR} ${AVDEVICE_INCLUDE_DIR} ${AVFILTER_INCLUDE_DIR} )
target_link_libraries(cms-encoding PUBLIC cmsEncodingLib ${SWSCALE_LIBRARY} ${AVCODEC_LIBRARY} ${AVFORMAT_LIBRARY} ${AVUTIL_LIBRARY} ${AVDEVICE_LIBRARY} ${AVFILTER_LIBRARY})
# TODO: Add tests and install targets if needed.

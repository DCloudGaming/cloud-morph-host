# CMakeList.txt : Top-level CMake project file, do global configuration
# and include sub-projects here.
#
cmake_minimum_required (VERSION 3.8)

project (CloudMorphHostSrc C CXX)

if(MSVC)
	set(test-input_PLATFORM_DEPS
		w32-pthreads)
endif()

# Include sub-projects.
add_subdirectory (contrib EXCLUDE_FROM_ALL)
add_subdirectory (CMakeProject2)

# CMakeList.txt : Top-level CMake project file, do global configuration
# and include sub-projects here.
#
cmake_minimum_required (VERSION 3.8)

project ("cms-encoding")

set(CMAKE_VERBOSE_MAKEFILE ON)
set(CMAKE_GENERATOR "Ninja" CACHE INTERNAL "" FORCE)
message("generator is set to ${CMAKE_GENERATOR}")

# Include sub-projects.
add_subdirectory ("cms-encoding")

#cmake ..  -DCMAKE_TOOLCHAIN_FILE=D:/repos/vcpkg/scripts/buildsystems/vcpkg.cmake -DCMAKE_GENERATOR:INTERNAL=Ninja
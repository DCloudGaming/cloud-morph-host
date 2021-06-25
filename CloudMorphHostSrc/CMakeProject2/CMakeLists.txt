# CMakeList.txt : CMake project for CMakeProject2, include source and define
# project specific logic here.
#
cmake_minimum_required (VERSION 3.8)

project(CMakeProject2)
set (CMAKE_CXX_STANDARD 11)


set(myplugin_SOURCES
	MyPlugin.cpp
	MySource.cpp
)
set(myplugin_HEADERS
	MyPlugin.h
	MySource.h
	 MyPlugin2.h)
add_library(my-plugin ${myplugin_SOURCES} ${myplugin_HEADERS})
target_link_libraries(my-plugin libobs w32-pthreads)
# Add source to this project's executable.
add_executable (my-server "CMakeProject2.cpp" "CMakeProject2.h")
target_link_libraries(my-server my-plugin libobs)
install_obs_plugin_with_data(my-plugin data)

set(my-gui_SOURCES GUI.cpp)
add_executable (my-gui WIN32 ${my-gui_SOURCES})
target_link_libraries(my-gui libobs)
set_target_properties(my-gui PROPERTIES FOLDER "tests and examples")
define_graphic_modules(my-gui)



#add_executable (CMakeProject2 "ClientApp2.cpp" "ClientApp2.h")

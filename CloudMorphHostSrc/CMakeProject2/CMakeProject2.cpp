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

int main()
{
	obs_module_load();
	//obs_module_load_2();
	cout << "Hello CMaksafasfgge." << endl;
	MessageBox(NULL,
		_T("Call to Create Window failed!"),
		_T("Windows Desktop Guided Tour"),
		NULL);
	return 0;
}

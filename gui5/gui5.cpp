// gui5.cpp : Defines the entry point for the application.
//

#include "pch.h"
#include "framework.h"
#include "gui5.h"
#include "string"
#include "map"
#include "vector"
#include "stdlib.h"
#include <cstdio>
#include <iostream>
#include <memory>
#include <stdexcept>
#include <string>
#include <array>
//#include <http_client.h>
//#include<filestream.h>
//#include <uri.h>
#include <iostream>       // std::cout, std::hex
#include <string>         // std::string, std::u32string
#include <locale>         // std::wstring_convert
#include <codecvt>        // std::codecvt_utf8
#include <cstdint> 
#include <wininet.h>  
#include <tchar.h>
//#include <curl/curl.h>

#define MAX_LOADSTRING 100
#define WM_LBUTTONDBLCLK 0x0203

using namespace std;


std::string exec(const char* cmd) {
    std::array<char, 128> buffer;
    std::string result;
    std::unique_ptr<FILE, decltype(&_pclose)> pipe(_popen(cmd, "r"), _pclose);
    if (!pipe) {
        throw std::runtime_error("popen() failed!");
    }
    while (fgets(buffer.data(), buffer.size(), pipe.get()) != nullptr) {
        result += buffer.data();
    }
    return result;
}

// Global Variables:
HINSTANCE hInst;                                // current instance
WCHAR szTitle[MAX_LOADSTRING];                  // The title bar text
WCHAR szWindowClass[MAX_LOADSTRING];            // the main window class name

bool g_MovingMainWnd = FALSE;
POINT g_OrigCursorPos;
POINT g_OrigWndPos;
WNDPROC prevWndProc;
std::map<wstring, vector<string>> exec_map = {
    {L"Note Pad", {"notepad", "notepad", "Untitled - Notepad"}},
    {L"Chrome Session", {"chrome", "start chrome", "Chrome Session"}}
};

std::map<string, vector<string>> chosen_apps = {};

// Forward declarations of functions included in this code module:
ATOM                MyRegisterClass(HINSTANCE hInstance);
BOOL                InitInstance(HINSTANCE, int);
LRESULT CALLBACK    WndProc(HWND, UINT, WPARAM, LPARAM);
LRESULT CALLBACK    ButtonWndProc(HWND, UINT, WPARAM, LPARAM);
INT_PTR CALLBACK    About(HWND, UINT, WPARAM, LPARAM);



int APIENTRY wWinMain(_In_ HINSTANCE hInstance,
                     _In_opt_ HINSTANCE hPrevInstance,
                     _In_ LPWSTR    lpCmdLine,
                     _In_ int       nCmdShow)
{
    exec("main.exe");
    UNREFERENCED_PARAMETER(hPrevInstance);
    UNREFERENCED_PARAMETER(lpCmdLine);

    // TODO: Place code here.

    // Initialize global strings
    LoadStringW(hInstance, IDS_APP_TITLE, szTitle, MAX_LOADSTRING);
    LoadStringW(hInstance, IDC_GUI5, szWindowClass, MAX_LOADSTRING);
    MyRegisterClass(hInstance);

    // Perform application initialization:
    if (!InitInstance (hInstance, nCmdShow))
    {
        return FALSE;
    }

    HACCEL hAccelTable = LoadAccelerators(hInstance, MAKEINTRESOURCE(IDC_GUI5));

    MSG msg;

    // Main message loop:
    while (GetMessage(&msg, nullptr, 0, 0))
    {
        if (!TranslateAccelerator(msg.hwnd, hAccelTable, &msg))
        {
            TranslateMessage(&msg);
            DispatchMessage(&msg);
        }
    }

    return (int) msg.wParam;
}



//
//  FUNCTION: MyRegisterClass()
//
//  PURPOSE: Registers the window class.
//
ATOM MyRegisterClass(HINSTANCE hInstance)
{
    WNDCLASSEXW wcex;

    wcex.cbSize = sizeof(WNDCLASSEX);

    wcex.style          = CS_HREDRAW | CS_VREDRAW;
    wcex.lpfnWndProc    = WndProc;
    wcex.cbClsExtra     = 0;
    wcex.cbWndExtra     = 0;
    wcex.hInstance      = hInstance;
    wcex.hIcon          = LoadIcon(hInstance, MAKEINTRESOURCE(IDI_GUI5));
    wcex.hCursor        = LoadCursor(nullptr, IDC_ARROW);
    wcex.hbrBackground  = (HBRUSH)(COLOR_WINDOW+1);
    wcex.lpszMenuName   = MAKEINTRESOURCEW(IDC_GUI5);
    wcex.lpszClassName  = szWindowClass;
    wcex.hIconSm        = LoadIcon(wcex.hInstance, MAKEINTRESOURCE(IDI_SMALL));

    return RegisterClassExW(&wcex);
}

//
//   FUNCTION: InitInstance(HINSTANCE, int)
//
//   PURPOSE: Saves instance handle and creates main window
//
//   COMMENTS:
//
//        In this function, we save the instance handle in a global variable and
//        create and display the main program window.
//
BOOL InitInstance(HINSTANCE hInstance, int nCmdShow)
{
   hInst = hInstance; // Store instance handle in our global variable

   HWND hWnd = CreateWindowW(szWindowClass, szTitle, WS_OVERLAPPEDWINDOW,
      CW_USEDEFAULT, 0, CW_USEDEFAULT, 0, nullptr, nullptr, hInstance, nullptr);

   if (!hWnd)
   {
      return FALSE;
   }


   // For now, we don't run run-app.bat yet, need to refactor

   map<wstring, vector<string>>::iterator it;
   int i = 0;
   for (auto const& x : exec_map) {
       auto key = x.first;
       const wchar_t* key_name = key.c_str();
       auto value = x.second;
       // TODO: Change y-position too when out of frame
       HWND hwndButton = CreateWindow( 
           L"BUTTON",  // Predefined class; Unicode assumed  TEXT("button")
           key_name,      // Button text 
           //WS_OVERLAPPEDWINDOW,
           WS_VISIBLE | WS_CHILD | BS_AUTOCHECKBOX,  // Styles 
           10 + 200 * i,         // x position 
           10,         // y position // TODO: Change y-position too when out of frame
           150,        // Button width
           100,        // Button height
           hWnd,     // Parent window
           (HMENU)IDC_SELECT_APP,       // No menu.
           hInstance,
           //(HINSTANCE)GetWindowLongPtr(hWnd, GWLP_HINSTANCE),
           NULL);      // Pointer not needed.

       //prevWndProc = (WNDPROC) SetWindowLongPtr(hwndButton, GWL_WNDPROC, (LONG_PTR)&ButtonWndProc);
       i += 1;
   }

   std::wstring register_key = L"Register Apps";
   const wchar_t* register_key_name = register_key.c_str();
   HWND hwndButton = CreateWindow(
       L"BUTTON",  // Predefined class; Unicode assumed 
       register_key_name,      // Button text 
       //WS_OVERLAPPEDWINDOW,
       WS_TABSTOP | WS_VISIBLE | WS_CHILD | BS_DEFPUSHBUTTON,  // Styles 
       10 + 500,         // x position 
       400,         // y position // TODO: Change y-position too when out of frame
       350,        // Button width
       100,        // Button height
       hWnd,     // Parent window
       (HMENU)IDC_REGISTER_APPS,       // No menu.
       //hInstance,
       (HINSTANCE)GetWindowLongPtr(hWnd, GWLP_HINSTANCE),
       NULL);      // Pointer not needed.

   ShowWindow(hWnd, nCmdShow);
   UpdateWindow(hWnd);

   return TRUE;
}

LRESULT CALLBACK ButtonWndProc(HWND hWnd, UINT message, WPARAM wParam, LPARAM lParam)
{
    switch (message)
    {
    case WM_LBUTTONDOWN:
        if (GetCursorPos(&g_OrigCursorPos))
        {
            HWND click_window = WindowFromPoint(g_OrigCursorPos);
            DialogBox(hInst, MAKEINTRESOURCE(IDD_ABOUTBOX), hWnd, About);
        }
        break;
    default:
        return DefWindowProc(hWnd, message, wParam, lParam);
    }
    return CallWindowProc(prevWndProc, hWnd, message, wParam, lParam);
}

//void sendRegistrationRequest() {
//    LPCWSTR 
//}

//
//  FUNCTION: WndProc(HWND, UINT, WPARAM, LPARAM)
//
//  PURPOSE: Processes messages for the main window.
//
//  WM_COMMAND  - process the application menu
//  WM_PAINT    - Paint the main window
//  WM_DESTROY  - post a quit message and return
//
//
LRESULT CALLBACK WndProc(HWND hWnd, UINT message, WPARAM wParam, LPARAM lParam)
{
    switch (message)
    {
    case WM_COMMAND:
        {
            int wmId = LOWORD(wParam);
            // Parse the menu selections:
            switch (wmId)
            {
            case IDM_ABOUT:
                DialogBox(hInst, MAKEINTRESOURCE(IDD_ABOUTBOX), hWnd, About);
                break;
            case IDM_EXIT:
                DestroyWindow(hWnd);
                break;
            case IDC_SELECT_APP:
                std::cout << "Select App " << std::endl;
                if (GetCursorPos(&g_OrigCursorPos))
                {
                    using convert_type = std::codecvt_utf8<wchar_t>;
                    std::wstring_convert<convert_type, wchar_t> converter;

                    HWND click_window = WindowFromPoint(g_OrigCursorPos);
                    //BOOL checked;
                    //checked = IsDlgButtonChecked(hWnd, IDC_SELECT_APP);
                    //chekced = CheckDlgButton(click_window, IDC_SELECT_APP, BST_CHECKED);
                    int len = GetWindowTextLength(click_window) + 1;
                    vector<wchar_t> buf(len);
                    GetWindowText(click_window, &buf[0], len);
                    wstring stxt = &buf[0];
                    std::string converted_str = converter.to_bytes(stxt);
 
                    //LONG style = GetWindowLong(click_window, GWL_STYLE);
                    //style = (style & ~BS_BOTTOM) | BS_TOP;
                    //style = WS_TABSTOP | WS_VISIBLE | WS_CHILD | BS_AUTOCHECKBOX;
                    //SetWindowLong(click_window, GWL_STYLE, style);

                    //if (SendDlgItemMessage(click_window, IDC_SELECT_APP, BM_GETCHECK, 0, 0)) {

                    /*if (checked == BST_CHECKED) {
                        CheckDlgButton(click_window, IDC_SELECT_APP, BST_UNCHECKED);
                        std::cout << "Check" << std::endl;
                    } else {
                        CheckDlgButton(click_window, IDC_SELECT_APP, BST_CHECKED);
                        std::cout << "Uncheck" << std::endl;
                    }*/
                    //exec("main.exe");

                    string app_name = exec_map[stxt][0];
                    string app_exec = exec_map[stxt][1];
                    string app_title = exec_map[stxt][2];

                    std::map<string, vector<string>>::iterator it;
                    it = chosen_apps.find(app_name);
                    if (it != chosen_apps.end()) {
                        chosen_apps.erase(it);
                    }
                    else {
                        chosen_apps[app_name] = { app_exec, app_title };
                    }

                    std::cout << converted_str << std::endl;

                }
                break;
            case IDC_REGISTER_APPS:
            {
                std::map<string, vector<string>>::iterator it2;
                bool is_begin = true;
                std::string param_str;
                for (it2 = chosen_apps.begin(); it2 != chosen_apps.end(); it2++) {
                    if (!is_begin) {
                        param_str.append(",");
                    }
                    else {
                        is_begin = false;
                    }
                    string app_name = it2->first;
                    string app_exec = it2->second[0];
                    string app_title = it2->second[1];
                    param_str.append(app_name);
                }

                //HINTERNET hIntSession =
                //    ::InternetOpen(_T("MyApp"), INTERNET_OPEN_TYPE_DIRECT, NULL, NULL, 0);

                //HINTERNET hHttpSession =
                //    InternetConnect(hIntSession, _T("localhost"), 8081, 0, 0, INTERNET_SERVICE_HTTP, 0, NULL);

                //HINTERNET hHttpRequest = HttpOpenRequest(
                //    hHttpSession,
                //    _T("GET"),
                //    _T("registerApp"),
                //    0, 0, 0, INTERNET_FLAG_RELOAD, 0);
                //TCHAR* szHeaders = _T("Content-Type: text/html\nMySpecialHeder: whatever");
                //CHAR szReq[1024] = "";
                //if (!HttpSendRequest(hHttpRequest, szHeaders, _tcslen(szHeaders), szReq, strlen(szReq))) {
                //    DWORD dwErr = GetLastError();
                //    /// handle error
                //}

                //CHAR szBuffer[1025];
                //DWORD dwRead = 0;
                //while (::InternetReadFile(hHttpRequest, szBuffer, sizeof(szBuffer) - 1, &dwRead) && dwRead) {
                //    szBuffer[dwRead] = 0;
                //    OutputDebugStringA(szBuffer);
                //    dwRead = 0;
                //}

                //::InternetCloseHandle(hHttpRequest);
                //::InternetCloseHandle(hHttpSession);
                //::InternetCloseHandle(hIntSession);

                // TODO: This is temporary stupid solution, since curl/curl.h import not working.

                string query_str = "curl http://localhost:8082/registerApp?data=";
                query_str.append(param_str);
                exec(query_str.c_str());

                /*CURL* curl;
                CURcode res;

                curl = curl_easy_init()
                if (curl) {
                    curl_easy_setopt(curl, CURLOPT_URL, "http://localhost:8082/registerApp");
                    curl_easy_setopt(curl, CURLOPT_POST, 1);
                    curl_easy_setopt(curl, CUROPT_POSTFIELDS, "name=Hieu&comment=A");
                }

                res = curl_easy_perform(curl);
                curl_easy_cleanup(curl);*/

                MessageBox(
                    NULL,
                    (LPCWSTR)L"Apps Registered",
                    (LPCWSTR)L"Account Details",
                    MB_DEFBUTTON2
                );
                break;
            }
            default:
                return DefWindowProc(hWnd, message, wParam, lParam);
            }
        }
        break;
    case WM_PAINT:
        {
            PAINTSTRUCT ps;
            HDC hdc = BeginPaint(hWnd, &ps);
            // TODO: Add any drawing code that uses hdc here...
            EndPaint(hWnd, &ps);
        }
        break;
    case WM_DESTROY:
        PostQuitMessage(0);
        break;
    case WM_LBUTTONDOWN:
        if (GetCursorPos(&g_OrigCursorPos))
        {
            HWND click_window = WindowFromPoint(g_OrigCursorPos);
            std::cout << "CLICK" << std::endl;
            //exec("main.exe");
            //DialogBox(hInst, MAKEINTRESOURCE(IDD_ABOUTBOX), hWnd, About);
        }
        break;
    default:
        return DefWindowProc(hWnd, message, wParam, lParam);
    }
    return 0;
}

// Message handler for about box.
INT_PTR CALLBACK About(HWND hDlg, UINT message, WPARAM wParam, LPARAM lParam)
{
    UNREFERENCED_PARAMETER(lParam);
    switch (message)
    {
    case WM_INITDIALOG:
        return (INT_PTR)TRUE;

    case WM_COMMAND:
        if (LOWORD(wParam) == IDOK || LOWORD(wParam) == IDCANCEL)
        {
            EndDialog(hDlg, LOWORD(wParam));
            return (INT_PTR)TRUE;
        }
        break;
    }
    return (INT_PTR)FALSE;
}

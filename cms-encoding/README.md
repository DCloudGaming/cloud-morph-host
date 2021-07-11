# Description:
- Captures the video and audio of a specified Window application, encode and stream through RTP.

# How-to-use:
Currently, there are two ways to perform this task. With the first basic MVP, it's fine to use approach 2 directly. But we need better flexibility for customization as it scales, so approach 1 is needed too.

## 1/ Through ffmpeg libraries
+ Note: Currently, this approach only supports saving to mp4 files, and also record only entire desktop screen. I am now working on supporting streaming to rtp protocol , as well as only recording an application screen instead of desktop screen

- Step 1: Compile the code based on Development Guide below. Or use the executable in cms-encoding/bin/*

## 2/ Through ffmpeg executables:
+ Note: Already can select which Windows application to record the screen, and stream through RTP

- Step 1: Enquire the window title of the application you are interested in
```
tasklist /v  /fo list
```

- Step 2: Records the application screen through gdigrab (there's another way through dshow, will check which one is better). Stream to rtp at port 5004. The "-pix_fmt" may differ based on different machines.
```
ffmpeg -f gdigrab  -framerate 30 -i title="<window_title_name_from_step-1>" -pix_fmt yuv420p -vf scale=1280:-2 -c:v libvpx -f rtp rtp://127.0.0.2:5004
```

# Development Guide
(For Windows Machine)
Working Directory: cms-encoding
- Step 1: Install vcpkg on windows. vcpkg is a tool to help compile and integrate popular C++ libraries easily
- Step 2: Install necessary libaries pre-built and make them available globally that CMake can later use easily
```
vcpkg install ffmpeg
vcpkg integrate install
```
- Step 3: Run CMake with Ninja as default generator instead of MSBuild, specify it to use vcpkg toolchain as well. Also build in debug mode so gdb can debug well.
```
mkdir build
cd build
cmake ..  -DCMAKE_TOOLCHAIN_FILE=D:/repos/vcpkg/scripts/buildsystems/vcpkg.cmake -DCMAKE_GENERATOR:INTERNAL=Ninja -DCMAKE_BUILD_TYPE=Debug
ninja
```

- Step 4: Download this to enable screen-record-capture as one input stream for ffmpeg: https://github.com/rdp/screen-capture-recorder-to-video-windows-free/releases



# Helper functions (for self-note)
- '.\ffmpeg.exe -list_devices true -f dshow -i dummy' 
- .\ffmpeg.exe -f dshow -fflags nobuffer -pix_fmt yuv420p -rtbufsize 300M -i video="screen-capture-recorder" output2.mp4
- .\ffmpeg.exe -f gdigrab -framerate 30 -i desktop -pix_fmt yuv420p  output11.mp4
- .\ffmpeg.exe -f gdigrab  -framerate 30 -i title="Windows PowerShell" -pix_fmt yuv420p -vf scale=1280:-2 output12.mp4
- tasklist /v  /fo list (This one is to get window title, ffmpeg can stream from a specified window title only)
  
# TODO: 
-  Compare gdigrab vs dshow performance, why one vs another?
-  window title (to stream from specific app) is currently gotten through terminal. Find a good UI/UX flow for that.
-  Stream to RTP through code.
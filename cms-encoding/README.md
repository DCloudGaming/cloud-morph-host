# Development Guide
(For Windows Machine)
Working Directory: cms-encoding
- Step 1: Install vcpkg on windows. vcpkg helps compile and integrate popular C++ libraries easily
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



# Other notes
- '.\ffmpeg.exe -list_devices true -f dshow -i dummy' 
- .\ffmpeg.exe -f dshow -fflags nobuffer -pix_fmt yuv420p -rtbufsize 300M -i video="screen-capture-recorder" output2.mp4
- .\ffmpeg.exe -f gdigrab -framerate 30 -i desktop -pix_fmt yuv420p  output11.mp4
- .\ffmpeg.exe -f gdigrab  -framerate 30 -i title="Windows PowerShell" -pix_fmt yuv420p -vf scale=1280:-2 output12.mp4
- tasklist /v  /fo list (This one is to get window title, ffmpeg can stream from a specified window title only)
  
# TODO: 
-  Compare gdigrab vs dshow performance, why one vs another?
-  window title (to stream from specific app) is currently gotten through terminal. Find a good UI/UX flow for that.
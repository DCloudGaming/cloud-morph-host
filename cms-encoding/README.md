# Development Guide
(For Windows Machine)
Working Directory: cms-encoding
- Step 1: Install vcpkg on windows
- Step 2: Install necessary libaries pre-built and make them available globally that CMake can later use easily
```
vcpkg install ffmpeg
vcpkg integrate install
```
- Step 3: Run CMake with Ninja as default generator instead of MSBuild
```
mkdir build
cd build
cmake ..  -DCMAKE_TOOLCHAIN_FILE=D:/repos/vcpkg/scripts/buildsystems/vcpkg.cmake -DCMAKE_GENERATOR:INTERNAL=Ninja
ninja
```


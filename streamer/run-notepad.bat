set /P PATH=PATH
START /b notepad
START /b D:\ffmpeg-2021-07-06-git-758e2da289-full_build\bin\ffmpeg.exe -f gdigrab  -framerate 30 -i title="Untitled - Notepad" -pix_fmt yuv420p -vf scale=1280:-2 -c:v libvpx -f rtp rtp://127.0.0.2:5006

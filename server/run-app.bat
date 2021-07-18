set /P PATH=PATH
taskkill /IM "ffmpeg.exe" /F
taskkill /IM "notepad.exe" /F
START /b notepad
@echo off 
g++ ../syncinput/syncinput.cpp -o ../syncinput/syncinput.exe -lws2_32 -lpthread -static
START /b ../syncinput/syncinput.exe "Notepad"
START /b ffmpeg -f gdigrab  -framerate 30 -i title="Untitled - Notepad" -pix_fmt yuv420p -vf scale=1280:-2 -c:v libvpx -f rtp rtp://127.0.0.2:5004
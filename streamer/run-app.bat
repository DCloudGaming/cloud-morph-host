set /P PATH=%PATH%
echo %1
TASKKILL /FI "ImageName eq %1" /F
TASKKILL /FI "ImageName eq ffmpeg.exe" /F
START /B %1
@REM Investigate running FFMPEG in background
START ffmpeg -f gdigrab -framerate 30 -i title=%2 -pix_fmt yuv420p -vf scale=1280:-2 -c:v libvpx -f rtp rtp://127.0.0.2:5006
echo "Done"
START /b notepad
START "" ffmpeg -f gdigrab  -framerate 30 -i title="Untitled - Notepad" -pix_fmt yuv420p -vf scale=1280:-2 -c:v libvpx -f rtp rtp://127.0.0.2:5004
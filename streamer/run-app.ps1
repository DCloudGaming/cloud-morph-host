# set /P PATH=%PATH%
$path=$args[0]
$title=$args[1]

echo $path
taskkill /FI "ImageName eq $path" /F
taskkill /FI "ImageName eq ffmpeg.exe" /F
$app = Start-Process $path -passthru
$app.id
# START /B %1
# Investigate running FFMPEG in background
Start-Process ffmpeg -ArgumentList "-f gdigrab -framerate 30 -i title=$title -pix_fmt yuv420p -vf scale=1280:-2 -c:v libvpx -f rtp rtp://127.0.0.2:5006"
echo "Done"
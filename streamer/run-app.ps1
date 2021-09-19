# set /P PATH=%PATH%
$path = $args[0]
$filename = $args[1]

echo $path
echo $filename
taskkill /FI "ImageName eq $filename" /F
taskkill /FI "ImageName eq ffmpeg.exe" /F
$app = Start-Process $path -PassThru
$app.id
sleep 2
$title = (Get-Process -Id $app.id).mainWindowTitle
# START /B %1
# Investigate running FFMPEG in background
sleep 2
echo "Title"$title
Start-Process ffmpeg -PassThru -ArgumentList "-f gdigrab -framerate 30 -i title=`"$title`" -pix_fmt yuv420p -vf scale=1280:-2 -c:v libvpx -f rtp rtp://127.0.0.2:5006"
echo "Done"
sleep 2
echo "Run Syncinput"
Start-Process syncinput.exe -PassThru -ArgumentList $title
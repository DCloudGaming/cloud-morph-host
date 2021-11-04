$path = $args[0]
$filename = $args[1]
$isSandbox = $args[2]
$hostIP = $args[3]

if ([string]::IsNullOrEmpty($hostIP)) {
    $hostIP = '127.0.0.1';
}
# Split-Path $outputPath -leaf
echo "running $path"

taskkill /FI "ImageName eq $filename" /F
taskkill /FI "ImageName eq ffmpeg.exe" /F
taskkill /FI "ImageName eq syncinput.exe" /F
# $app = Start-Process $path -PassThru
$app = Start-Process "C:\Program Files\Sandboxie\Start.exe" -ArgumentList $path -PassThru
sleep 2
$title = ((Get-Process -Id $app.id).mainWindowTitle)
echo "Computed Title $title"
# $title = "[#] Deverse (64-bit Development PCD3D_SM5) [#]"
$title = "[#] 7554 [#]"
echo "Title $title"
sleep 2
x86_64-w64-mingw32-g++ $PSScriptRoot\syncinput.cpp -o $PSScriptRoot\syncinput.exe -lws2_32 -lpthread -static
if ($isSandbox -eq "sandbox") {
    Start-Process $PSScriptRoot/pkg/ffmpeg/ffmpeg.exe -PassThru -NoNewWindow -ArgumentList "-f gdigrab -framerate 30 -i title=`"$title`" -pix_fmt yuv420p -vf scale=1280:-2 -tune zerolatency -c:v libx264 -f rtp rtp://$hostIP`:5006"
    sleep 2
    while ($true) {
        Start-Process -Wait $PSScriptRoot/syncinput.exe -PassThru -NoNewWindow -ArgumentList "$title", ".", "windows", $hostIP
    }
    # Restart on failure. Using service to restart on failure, not working now
    # $syncinput = New-Service -Name "Syncinput" -BinaryPathName "$PSScriptRoot\syncinput.exe $title . windows $hostIP"
    # sc failure Syncinput reset= 30 actions= restart/5000
    # $syncinput.Start()
}
else {
    Start-Process ffmpeg -PassThru -ArgumentList "-f gdigrab -framerate 30 -i title=`"$title`" -pix_fmt yuv420p -vf scale=1280:-2 -tune zerolatency -c:v libx264 -f rtp rtp://127.0.0.1:5006"
    sleep 2
    Start-Process -PassThru $PSScriptRoot/syncinput.exe -ArgumentList "$title", ".", "windows"
}
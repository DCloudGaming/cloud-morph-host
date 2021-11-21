# 0: absolute (path + filename) on host to folder
# 1: filename
# 2: relative dir path on host
# 3: sandbox flag 
# 4: local Ethernet IP (Will be different if it's in virtual network like sandbox/Docker)

$template = @'
<Configuration>
    <vGPU>Enable</vGPU>
    <Networking>Default</Networking>
    <MappedFolders>
        <MappedFolder>
            <HostFolder>{2}</HostFolder>
            <SandboxFolder>C:\Users\declo</SandboxFolder>
        </MappedFolder>
        <MappedFolder>
            <HostFolder>{0}</HostFolder>
            <SandboxFolder>C:\Users\declo\apps</SandboxFolder>
        </MappedFolder>
    </MappedFolders>
    <LogonCommand>
        <Command>C:\\Windows\\System32\\WindowsPowerShell\\v1.0\\powershell.exe -ExecutionPolicy Bypass -F C:\Users\declo\run-app.ps1 C:\Users\declo\apps\ {1} C:\Users\declo\apps\{1} sandbox {3}</Command>
    </LogonCommand>
</Configuration>
'@

# To install Virtual Box Image. Copy FFMPEG to VM

# Create Sandbox Config

$localEthernetIP = (Get-NetIPAddress -AddressFamily IPv4 -InterfaceAlias ethernet).IPAddress
# pass variables in orders to template
$template -f $args[0], $args[1], "$PWD", $localEthernetIP  | Out-File -FilePath .\run-sandbox.wsb
# x86_64-w64-mingw32-g++ $PSScriptRoot\winvm\syncinput.cpp -o $PSScriptRoot\winvm\syncinput.exe -lws2_32 -lpthread -static

powershell -ExecutionPolicy Bypass -F "setup-sandbox.ps1"
# Run Sandbox
.\run-sandbox.wsb
do {
    $rsp = Read-Host "install extras? (y/n)"
} while ($rsp -notmatch '^[YyNn]$')

if ($rsp -match '^[Yy]$') {
    $extra = $true
} else {
    $extra = $false
}

Set-ExecutionPolicy Bypass -Scope Process -Force
[System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072
iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))

choco install notepadplusplus -y
choco install 7zip -y
choco install golang -y

if ($extra) {
    Write-Host "installing extras..."
    choco install python -y
    choco install git -y
    go install -v github.com/projectdiscovery/nuclei/v3/cmd/nuclei@latest
}
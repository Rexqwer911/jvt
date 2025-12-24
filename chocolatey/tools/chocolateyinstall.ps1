$ErrorActionPreference = 'Stop'

$packageName = 'jvt'
$toolsDir = "$(Split-Path -parent $MyInvocation.MyCommand.Definition)"
$url64 = 'https://github.com/rexqwer911/jvt/releases/download/v1.2.0/jvt-windows-amd64.zip'

$packageArgs = @{
  packageName    = $packageName
  unzipLocation  = $toolsDir
  url64bit       = $url64
  checksum64     = '3E193B663B83D1464BA67A715BD25C2CA159193B5092016A934A478422DF37E0'
  checksumType64 = 'sha256'
}

Install-ChocolateyZipPackage @packageArgs


Write-Host "JVT has been installed successfully!"
Write-Host "Run 'jvt --help' to get started."

$ErrorActionPreference = 'Stop'

$packageName = 'jvt'
$toolsDir = "$(Split-Path -parent $MyInvocation.MyCommand.Definition)"
$url64 = 'https://github.com/rexqwer911/jvt/releases/download/v1.0.0/jvt-windows-amd64.zip'

$packageArgs = @{
  packageName    = $packageName
  unzipLocation  = $toolsDir
  url64bit       = $url64
  checksum64     = ''  # Will be filled during release
  checksumType64 = 'sha256'
}

Install-ChocolateyZipPackage @packageArgs

# Add to PATH
Install-ChocolateyPath -PathToInstall $toolsDir -PathType 'User'

Write-Host "JVT has been installed successfully!"
Write-Host "Run 'jvt --help' to get started."

$ErrorActionPreference = 'Stop'

$packageName = 'jvt'
$toolsDir = "$(Split-Path -parent $MyInvocation.MyCommand.Definition)"
$url64 = 'https://github.com/rexqwer911/jvt/releases/download/v1.3.0/jvt-windows-amd64.zip'

$packageArgs = @{
  packageName    = $packageName
  unzipLocation  = $toolsDir
  url64bit       = $url64
  checksum64     = '55F71B88BCC4A7934055206220046C7CDA9A19BC8AE5C48273AB6FD17F0772EB'
  checksumType64 = 'sha256'
}

Install-ChocolateyZipPackage @packageArgs


Write-Host "JVT has been installed successfully!"
Write-Host "Run 'jvt --help' to get started."

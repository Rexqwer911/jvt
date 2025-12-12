$ErrorActionPreference = 'Stop'

$toolsDir = "$(Split-Path -parent $MyInvocation.MyCommand.Definition)"

# Remove from PATH
Uninstall-ChocolateyPath -PathToUninstall $toolsDir -PathType 'User'

Write-Host "JVT has been uninstalled."
Write-Host "Note: Installed Java versions in ~/.jvt were not removed."
Write-Host "To remove them manually, delete the folder: $env:USERPROFILE\.jvt"

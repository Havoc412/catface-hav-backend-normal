$CurrentPath = Split-Path -Parent $MyInvocation.MyCommand.Definition

# Define paths to check
$NginxConf = Join-Path $CurrentPath "nginx.conf"
$DefaultConf = Join-Path $CurrentPath "conf.d\default.conf"

$DataVolume = Join-Path $CurrentPath "..\..\public\nginx"
$LogVolume = Join-Path $CurrentPath "..\..\store\logs\nginx"

$pathsToCheck = @(
    $NginxConf,
    $DataVolume,
    $LogVolume,
    $DefaultConf
)

# INFO Check if each path exists
$allPathsExist = $true
foreach ($path in $pathsToCheck) {
    if (-not (Test-Path -Path $path)) {
        Write-Host "Path does not exist1: $path"
        $allPathsExist = $false
    }
    & cmd /c echo Path Check: $path
}

# If all paths exist, start the Docker container
if ($allPathsExist) {
    docker run -d --name nginx1 -p 80:80 `
    -v "${NginxConf}:/etc/nginx/nginx.conf" `
    -v "${DataVolume}:/mnt/data" `
    -v "${LogVolume}:/var/log/nginx" `
    -v "${DefaultConf}:/etc/nginx/conf.d/default.conf" nginx
} else {
    Write-Host "One or more paths do not exist, cannot start the Docker container."
}

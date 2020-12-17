$savedLocation = Get-Location
if ($args.Count -eq 0) {
    Write-Output "Building..."
    # Ensure that the dist directory exists and that it contains 
    # at least the sample configuration, or my personal configuration,
    # if my personal repo exists.
    if (!(Test-Path -Path "dist")) {
        New-Item -ItemType directory dist
        if (Test-Path -Path "..\..\bc\dotfiles\photo\config.json") {
            Copy-Item ..\..\bc\dotfiles\photo\config.json dist\config.json
        }
        else {
            Copy-Item config-sample.json dist\config.json
        }
        Copy-Item -Recurse exiftool dist\exiftool
    }
    $Env:GOPATH = $PSScriptRoot
    # Ensure that the third-party libraries are present
    if (!(Test-Path -Path "src\github.com\rwcarlsen\goexif")) {
        go get github.com/rwcarlsen/goexif/exif
    }
    if (!(Test-Path -Path "src\golang.org\x\crypto")) {
        go get golang.org/x/crypto/ssh
    }
    if (!(Test-Path -Path "src\github.com\tmc\scp")) {
        go get github.com/tmc/scp
    }
    # Linux/amd64 build
    $Env:GOOS = "linux"
    $Env:GOARCH = "amd64"
    go build github.com/bernarpa/photo
    Move-Item -Force photo dist\photo-linux
    # Windows/amd64 build
    $Env:GOOS = "windows"
    $Env:GOARCH = "amd64"
    go build github.com/bernarpa/photo
    Move-Item -Force photo.exe dist\photo-win.exe
    # Mac/amd64 build
    $Env:GOOS = "darwin"
    $Env:GOARCH = "amd64"
    go build github.com/bernarpa/photo
    Move-Item -Force photo dist\photo-mac
}
elseif ($args[0] -eq "clean") {
    Write-Output "Cleaning..."
    Remove-Item -Force -Recurse "bin"
    Remove-Item -Force -Recurse "pkg"
    Remove-Item -Force -Recurse "dist"  
    Remove-Item -Force -Recurse "src\golang.org"
    Remove-Item -Force -Recurse "src\github.com\tmc"
    Remove-Item -Force -Recurse "src\github.com\kballard"
    Remove-Item -Force -Recurse "src\github.com\rwcarlsen"
}
elseif ($args[0] -eq "install") {
    # This bit is veeery customized for my system
    Write-Output "Installing..."
    $paolosCustomPath = "..\..\Nextcloud\App\photo"
    if (!(Test-Path -Path $paolosCustomPath) -and ($args.Count -eq 1)) {
        Write-Output "Please specify the installation directory"
    }
    else {
        if ($args.Count -eq 1) {
            $installPath = $paolosCustomPath
        }
        else {
            $installPath = $args[1]
        }
        if (!(Test-Path -Path $installPath)) {
            New-Item -ItemType directory "$installPath"
        }
        Copy-Item -Recurse ".\dist\*" "$installPath"
    }
}
Set-Location $savedLocation
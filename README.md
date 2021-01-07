<a href="https://www.bernardi.cloud/">
    <img src=".readme-files/photo-logo-72.png" alt="Photo logo" title="Photo" align="right" height="72" />
</a>

# Photo
> A command line tool to manage your (my!😉) photo library

[![Go](https://img.shields.io/badge/Go-v1.15-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/github/license/bernarpa/photo.svg)](https://opensource.org/licenses/GPL-3.0)
[![GitHub issues](https://img.shields.io/github/issues/bernarpa/photo.svg)](https://github.com/bernarpa/photo/issues)

## Table of contents

- [What is Photo](#what-is-photo)
- [Installation](#installation)
- [License](#license)

## What is Photo

Photo is a command line tool that I created to manage my photo library according to my established *modus operandi*.

Basically with Photo you can manage a collection of JPEG photographs stored locally or on a computer accessible via SSH (such as my NAS, what a coincidence huh?).

Currently Photo supports the following operations:

1. **stats**: prints statistics about the photo collection, such as the most recent photos uploaded for each camera.
2. **filter**: filters the photos contained in a local directory by separating these already in the collection from the new ones, which are neatly renamed and organized in "daily" folders.
3. **update**: manually update the collection index cache (please note that the *stats* and *filter* operations will automatically performe an update if the collection index cache is not present of if it is older than one day).
4. **fix**: renames JPEG files accordingly to their Exif timestamp and converts HEIC files to the JPEG format. This command doesn't require a target.
5. **info**: shows Exif metadata for a supported image file. This command doesn't require a target.
6. **ignore**: creates a `photoignore` file, which can be uploaded to the photo collection, which marks all the photos in the specified local directory as ignored with respect to the *filter* command.

Please note that Photo is a multi-platform tool. It supports any combination of Linux, Windows and Mac (currently Intel only, as it's what I own) systems. Depending on your system, you should use one of the following executables to run Photo:

* `photo-windows.exe`
* `photo-linux`
* `photo-mac`

## Installation

Photo is written in Go; I've tested it with Go 1.15 and I suggest you to use at least that version to compile it.

To compile the program you should use `.\make.ps1` on Windows (it requires a recent PowerShell - I use 7.1 - with script execution enabled) or `make` on Linux and Mac systems. In any case, a `dist` directory will be created, together with a `config.json` file that must be customized to match your system parameters.

Please note that the **filter** command converts HEIC files in JPEG format (I know, I know...), but only if [ImageMagick](https://imagemagick.org/) is installed in the local system.

### config.json

* **workers**: number of parallel "goroutines" used by parallel operations (e.g. *update*)
* **targets**: remote or local photo library.
* **target.name**: name of the photo library, to be used in the photo command line.
* **target.target_type**: `local` or `ssh`.
* **target.work_dir**: local or remote working directory; Photo actually copies its executable (see *target.ssh_exe*) to this directory, in order to run on the remote system.
* **target.ssh_\***: SSH configuration parameters (currently only password authentication is supported). Please note that *ssh_exe* is the name of the Photo executable file to be used on the remote platform (e.g. for a Linux NAS you should use `photo-linux`).
* **target.collections**: list of the directories that contain the photo library. Photo analyizes each of them recursively, so only the root directories should be specified.
* **target.cameras**: camera models of interest, used by the *stat* operation (unless `--all` is specified).

# License

Photo is licensed under the terms of the GNU General Public License version 3.

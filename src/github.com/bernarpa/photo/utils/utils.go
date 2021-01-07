package utils

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// GetExePath returns the path of the executable file or Exit(1) if it cannot read it.
func GetExePath() string {
	exe, err := os.Executable()
	if err != nil {
		log.Fatal("Unable to get the path of the executable file")
	}
	return filepath.Dir(exe)
}

// EnsureDir ensures that the specified directory exists
func EnsureDir(directory string) string {
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		os.MkdirAll(directory, 0755)
	}
	return directory
}

// MD5 computes the MD5 hash of a file.
func MD5(path string) (string, error) {
	var md5Hash string
	file, err := os.Open(path)
	if err != nil {
		return md5Hash, err
	}
	defer file.Close()
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return md5Hash, err
	}
	hashInBytes := hash.Sum(nil)[:16]
	md5Hash = hex.EncodeToString(hashInBytes)
	return md5Hash, nil
}

// HeicToJPEG converts an HEIC image to a JPEG image.
// It requires ImageMagick in the PATH (convert for Unix platforms, magick.exe for Windows).
func HeicToJPEG(heicFile, jpegFile string) error {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("magick", "convert", heicFile, jpegFile)
	} else {
		cmd = exec.Command("convert", heicFile, jpegFile)
	}
	return cmd.Run()
}

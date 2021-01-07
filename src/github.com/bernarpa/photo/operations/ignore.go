package operations

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/bernarpa/photo/cache"
	"github.com/bernarpa/photo/config"
	"github.com/bernarpa/photo/exiftool"
)

// ShowHelpIgnore prints the help for the info operation.
func ShowHelpIgnore() {
	fmt.Println()
	fmt.Println("Usage: photo ignore [directory]")
	fmt.Println()
	fmt.Println("   directory       directory containing the files to ignore (recursive),")
	fmt.Println("                   by default it's the current directory")
	fmt.Println()
}

// Ignore creates a photoignore file with the files in the current directory.
// It process all files, recursively.
func Ignore(conf *config.Config, target *config.Target) {
	var targetDir string
	if len(os.Args) == 3 {
		targetDir = os.Args[2]
	} else {
		targetDir = "."
	}
	et := exiftool.Create(conf.Perl)
	log.Printf("exiftool created: %s\n", et.Perl)
	myCache := cache.Create(target)
	err := myCache.AnalyzeDir(targetDir, conf.Workers, et, []string{})
	if err != nil {
		log.Fatal("Cache update failure: " + err.Error())
	}
	jsonContent, err := json.Marshal(myCache)
	now := time.Now()
	nowStr := now.Format("2006-01-02_15-04-05")
	photoIgnoreFileName := fmt.Sprintf("photoignore_%s.json.gz", nowStr)
	photoIgnorePath := filepath.Join(targetDir, photoIgnoreFileName)
	f, err := os.Create(photoIgnorePath)
	if err != nil {
		log.Fatal("Photoignore file creation error: " + err.Error())
	}
	w := gzip.NewWriter(f)
	defer w.Close()
	_, err = w.Write(jsonContent)
	if err != nil {
		log.Fatal("Photoignore file writing error: " + err.Error())
	}
}

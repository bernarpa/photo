package operations

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/bernarpa/photo/cache"
	"github.com/bernarpa/photo/config"
	"github.com/bernarpa/photo/exiftool"
	"github.com/bernarpa/photo/utils"
)

// ShowHelpFilter prints the help for the stats operation.
func ShowHelpFilter() {
	fmt.Println()
	fmt.Println("Usage: photo filter <TARGET> [directory]")
	fmt.Println()
	fmt.Println("   TARGET     one of the targets defined in config.json")
	fmt.Println("   directory  local directory with the photos to be filtered")
	fmt.Println()
}

// Filter analyzes the photos in the current local directory, puts
// these that are already present in the target in the "Trash" directory
// and reorganizes the new ones in daily folders.
func Filter(conf *config.Config, target *config.Target) {
	var localDir string
	if len(os.Args) == 4 {
		localDir = os.Args[3]
	} else {
		localDir = "."
	}
	et := exiftool.Create(conf.Perl)
	duplicatesDir := utils.EnsureDir(filepath.Join(localDir, "AlreadyImported"))
	noExifDir := utils.EnsureDir(filepath.Join(localDir, "NoExif"))
	newDir := utils.EnsureDir(filepath.Join(localDir, "ToBeImported"))
	myCache := loadLocalCache(conf, target)
	localCache := cache.Create(target)
	localCache.AnalyzeDir(localDir, conf.Workers, et, target.Ignore)
	// Create an hash map of the target cache
	hashMap := make(map[string]cache.Photo)
	for _, targetPhoto := range myCache.Photos {
		hashMap[targetPhoto.Hash] = targetPhoto
	}
	// I've loaded both caches, now I should find
	// photos that are on localCache but NOT on myCache
	for _, localPhoto := range localCache.Photos {
		fmt.Printf("Filtering %s\n", localPhoto.Path)
		localPhoto.HeicToJPEG(et)
		if !localPhoto.HasExif() {
			newPath := filepath.Join(noExifDir, filepath.Base(localPhoto.Path))
			err := os.Rename(localPhoto.Path, newPath)
			if err != nil {
				log.Printf("Warning: unable to move photo %s to %s: %s\n", localPhoto.Path, newPath, err.Error())
			}
		} else {
			targetPhoto, exists := hashMap[localPhoto.Hash]
			if exists {
				log.Printf("Photo already exists in the target:\n  (%s) %s\n  (%s) %s\n", localPhoto.Hash, localPhoto.Path, targetPhoto.Hash, targetPhoto.Path)
				newPath := filepath.Join(duplicatesDir, filepath.Base(localPhoto.Path))
				err := os.Rename(localPhoto.Path, newPath)
				if err != nil {
					log.Printf("Warning: unable to move photo %s to %s\n", localPhoto.Path, newPath)
				}
			} else {
				// Rename the JPEG file according to its Exif timestamp
				err := localPhoto.RenameToExif()
				if err != nil {
					log.Printf("Warning: unable to rename photo %s according to Exif: %s\n", localPhoto.Path, err.Error())
					continue
				}
				// Ensure that the daily directory yyyy-mm-dd exists
				t := time.Unix(localPhoto.Timestamp, 0)
				dailyDir := filepath.Join(newDir, t.Format("2006-01-02"))
				if _, err := os.Stat(dailyDir); os.IsNotExist(err) {
					os.Mkdir(dailyDir, 0755)
				}
				newPath := filepath.Join(dailyDir, filepath.Base(localPhoto.Path))
				err = os.Rename(localPhoto.Path, newPath)
				if err != nil {
					log.Printf("Warning: unable to move photo %s to %s\n", localPhoto.Path, newPath)
				}
			}
		}
	}
}

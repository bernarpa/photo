package operations

import (
	"fmt"
	"log"
	"os"

	"github.com/bernarpa/photo/cache"
	"github.com/bernarpa/photo/config"
	"github.com/bernarpa/photo/exiftool"
)

// ShowHelpFix prints the help for the info operation.
func ShowHelpFix() {
	fmt.Println()
	fmt.Println("Usage: photo fix [directory]")
	fmt.Println()
	fmt.Println("   directory  local directory with the photos to be fixed")
	fmt.Println()
}

// Fix renames the photo in the specified directory according to
// their Exif timestamps. HEIC photos are converted to JPEG.
func Fix(conf *config.Config, target *config.Target) {
	var localDir string
	if len(os.Args) == 3 {
		localDir = os.Args[2]
	} else {
		localDir = "."
	}
	localCache := cache.Create(target)
	et, err := exiftool.Create(conf, target)
	if err != nil {
		log.Printf("exiftool instantation error: %s\n", err.Error())
		return
	}
	localCache.AnalyzeDir(localDir, conf.Workers, et, target.Ignore)
	for _, localPhoto := range localCache.Photos {
		fmt.Printf("Fixing %s\n", localPhoto.Path)
		localPhoto.HeicToJPEG(et)
		if localPhoto.Timestamp == 0 {
			fmt.Println("no timestamp")
			continue
		}
		err := localPhoto.RenameToExif()
		if err != nil {
			log.Printf("Warning: unable to rename photo %s according to Exif: %s\n", localPhoto.Path, err.Error())
			continue
		}
	}
}

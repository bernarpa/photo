package operations

import (
	"fmt"
	"log"
	"os"

	"github.com/bernarpa/photo/cache"
	"github.com/bernarpa/photo/config"
)

// ShowHelpFix prints the help for the stats operation.
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
	localCache.AnalyzeDir(localDir, conf.Workers)
	for _, localPhoto := range localCache.Photos {
		fmt.Printf("Fixing %s\n", localPhoto.Path)
		localPhoto.HeicToJPEG()
		if !localPhoto.HasExif() {
			continue
		}
		err := localPhoto.RenameToExif()
		if err != nil {
			log.Printf("Warning: unable to rename photo %s according to Exif: %s\n", localPhoto.Path, err.Error())
			continue
		}
	}
}

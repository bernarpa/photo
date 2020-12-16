package operations

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/bernarpa/photo/cache"
	"github.com/bernarpa/photo/config"
)

// ShowHelpStats prints the help for the stats operation.
func ShowHelpStats() {
	fmt.Println()
	fmt.Println("Usage: photo stats <TARGET> [--all]")
	fmt.Println()
	fmt.Println("   TARGET     one of the targets defined in config.json")
	fmt.Println("   --all      show statistics for all cameras;")
	fmt.Println("              if not specified, use the cameras defined in config.json")
	fmt.Println()
}

// Stats shows interesting information and statistics about the
// specified target. The information is inferred from the cache file,
// which will be created if it doesn't exist or it will be updated if
// it is too old.
func Stats(conf *config.Config, target *config.Target) {
	allCameras := len(os.Args) == 4 && os.Args[3] == "--all"
	myCache := loadLocalCache(conf, target)
	// Provide the user with a summary of the most recent photo timestamps
	// for each camera model
	lastPhoto := make(map[string]cache.Photo)
	for _, photo := range myCache.Photos {
		last, exists := lastPhoto[photo.Camera]
		if !exists || last.Timestamp < photo.Timestamp {
			lastPhoto[photo.Camera] = photo
		}
	}
	var cameras []string
	title := "Latest photo per camera"
	if allCameras {
		for camera := range lastPhoto {
			cameras = append(cameras, camera)
		}
		title += " (all cameras)"
	} else {
		cameras = target.Cameras
	}
	sort.Strings(cameras)
	fmt.Printf("%s\n%s\n", title, strings.Repeat("=", len(title)))
	maxCameraLen := 0
	for _, camera := range cameras {
		if len(camera) > maxCameraLen {
			maxCameraLen = len(camera)
		}
	}
	maxTimestampLen := len("2008-06-01 22:11:04 +0200 CEST")
	for _, camera := range cameras {
		photo, exists := lastPhoto[camera]
		var strTime string
		var path string
		if !exists {
			strTime = "-"
			path = "-"
		} else {
			strTime = time.Unix(photo.Timestamp, 0).String()
			path = photo.Path
		}
		spaces1 := strings.Repeat(" ", maxCameraLen-len(camera))
		spaces2 := strings.Repeat(" ", maxTimestampLen-len(strTime))
		fmt.Printf("%s %s %s %s %s\n", camera, spaces1, strTime, spaces2, path)
	}
}

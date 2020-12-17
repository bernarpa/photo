package operations

import (
	"fmt"
	"log"
	"os"

	"github.com/bernarpa/photo/config"
	"github.com/bernarpa/photo/exiftool"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/tiff"
)

// ShowHelpInfo prints the help for the info operation.
func ShowHelpInfo() {
	fmt.Println()
	fmt.Println("Usage: photo info [file]")
	fmt.Println()
	fmt.Println("   file       photo or video file")
	fmt.Println()
}

type infoWalker struct{}

func (infoWalker) Walk(name exif.FieldName, tag *tiff.Tag) error {
	data, _ := tag.MarshalJSON()
	fmt.Printf("    %s: %s\n", name, string(data))
	return nil
}

// Info tries to print the metadata of the photo or video file.
func Info(conf *config.Config, target *config.Target) {
	var fileName string
	if len(os.Args) == 3 {
		fileName = os.Args[2]
	} else {
		ShowHelpInfo()
		return
	}
	et, err := exiftool.Create(conf, target)
	if err != nil {
		log.Printf("exiftool initialization error: %s\n", err.Error())
		return
	}
	et.Dump(fileName)
}

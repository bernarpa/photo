package operations

import (
	"fmt"
	"log"
	"os"

	"github.com/bernarpa/photo/config"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
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
	exif.RegisterParsers(mknote.All...)
	f, err := os.Open(fileName)
	if err != nil {
		log.Printf("Error while opening %s: %s\n", fileName, err.Error())
		return
	}
	x, err := exif.Decode(f)
	if err != nil {
		log.Printf("Error while decoding metadata for %s: %s\n", fileName, err.Error())
		return
	}
	x.Walk(infoWalker{})
}

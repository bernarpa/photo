package exiftool

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/bernarpa/photo/utils"
)

// Output describes part of the exiftool -json output.
type Output struct {
	Timestamp        int64
	DateTimeOriginal string `json:"DateTimeOriginal"`
	MediaCreateDate  string `json:"MediaCreateDate"`
	Make             string `json:"Make"`
	Model            string `json:"Model"`
}

// Exiftool is a wrapper around the exiftool Perl program
type Exiftool struct {
	Perl string
}

// Create creates a new Exiftool wrapper instance.
func Create(perl string) *Exiftool {
	return &Exiftool{Perl: perl}
}

func parseExifTstamp(exifTstamp string) int64 {
	tm, err := time.Parse("2006:01:02 15:04:05-0700", exifTstamp)
	if err != nil {
		tm, err = time.Parse("2006:01:02 15:04:05", exifTstamp)
	}
	if err != nil {
		return 0
	}
	return tm.Unix()
}

// Parse parses the tags for the specified file by using exiftool.
func (et *Exiftool) Parse(fileName string) (*Output, error) {
	exePath := utils.GetExePath()
	exiftoolExe := filepath.Join(exePath, "exiftool", "exiftool")
	cmd := exec.Command(et.Perl, exiftoolExe, "-json", fileName)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	var outputs []Output
	err = json.Unmarshal(out, &outputs)
	if err != nil {
		return nil, err
	}
	for i := range outputs {
		tstamp := parseExifTstamp(outputs[i].DateTimeOriginal)
		if tstamp == 0 {
			tstamp = parseExifTstamp(outputs[i].MediaCreateDate)
		}
		if err == nil {
			outputs[i].Timestamp = tstamp
		}
	}
	return &outputs[0], nil
}

// Dump prints the tags for the specified file by using exiftool.
func (et *Exiftool) Dump(fileName string) {
	exePath := utils.GetExePath()
	exiftoolExe := filepath.Join(exePath, "exiftool", "exiftool")
	cmd := exec.Command(et.Perl, exiftoolExe, fileName)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Parsing error: %s\n", err.Error())
	}
	fmt.Printf("%s\n", out)
}

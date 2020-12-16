package cache

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/bernarpa/photo/config"
	"github.com/bernarpa/photo/utils"
	"github.com/rwcarlsen/goexif/exif"
)

// Cache is the struct that represents a Photo cache JSON file.
type Cache struct {
	Target     string  `json:"target"`
	LastUpdate int64   `json:"last_update"`
	Photos     []Photo `json:"photos"`
}

// Photo represents a JPEG file entry of a JSON cache "photos" property.
type Photo struct {
	Path      string `json:"path"`
	Size      int64  `json:"size"`
	Timestamp int64  `json:"tstamp"`
	Camera    string `json:"camera"`
	Hash      string `json:"hash"`
}

// HasExif checks whether the photo has Exif metadata.
func (photo *Photo) HasExif() bool {
	return photo.Timestamp != 0 && photo.Camera != ""
}

// HeicToJPEG converts an HEIC photo to the JPEG format.
// If the photo is not an HEIC file or if there is already a
// file with the same name but .jpg extension this function
// does nothing.
func (photo *Photo) HeicToJPEG() error {
	ext := filepath.Ext(photo.Path)
	if strings.ToLower(ext) == ".heic" {
		jpg := strings.TrimSuffix(photo.Path, ext) + ".jpg"
		if _, err := os.Stat(jpg); !os.IsNotExist(err) {
			return nil
		}
		err := utils.HeicToJPEG(photo.Path, jpg)
		if err != nil {
			return err
		}
		// If the conversion was successful, analyze the newly created JPEG
		jpgInfo, err := os.Stat(jpg)
		if err != nil {
			return err
		}
		jpgPhoto, err := AnalyzePhoto(jpg, jpgInfo)
		os.Remove(photo.Path)
		if err != nil {
			log.Printf("Warning: unable to analyze %s: %s\n", jpg, err.Error())
			photo.Path = jpg
		} else {
			photo.Path = jpgPhoto.Path
			photo.Size = jpgPhoto.Size
			photo.Timestamp = jpgPhoto.Timestamp
			photo.Camera = jpgPhoto.Camera
			photo.Hash = jpgPhoto.Hash
		}
	}
	return nil
}

// RenameToExif renames the photo according to the Exif timestamp
// in the YYYY-MM-DD_HH-MM-SS.jpg format.
func (photo *Photo) RenameToExif() error {
	if photo.Timestamp != 0 {
		t := time.Unix(photo.Timestamp, 0)
		timeStr := t.Format("2006-01-02_15-04-05")
		newFileName := timeStr + ".jpg"
		newPath := filepath.Join(filepath.Dir(photo.Path), newFileName)
		err := os.Rename(photo.Path, newPath)
		if err != nil {
			fmt.Printf("Warning: error while renaming %s to %s: %s\n", photo.Path, newPath, err.Error())
			return err
		}
		photo.Path = newPath
	}
	return nil
}

// Create returns an empty Cache.
func Create(target *config.Target) *Cache {
	if target != nil {
		return &Cache{Target: target.Name, LastUpdate: time.Now().Unix()}
	} else {
		return &Cache{Target: "", LastUpdate: time.Now().Unix()}
	}
}

// Load loads a cache from a cache file. If the cache file doesn't exist
// exist or if the cache is too old, the cache will be updated.
func Load(conf *config.Config, target *config.Target) (*Cache, error) {
	filename := target.GetLocalCachePath()
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	gz, err := gzip.NewReader(file)
	if err != nil {
		return nil, err
	}
	defer gz.Close()
	content, err := ioutil.ReadAll(gz)
	if err != nil {
		return nil, err
	}
	var c Cache
	err = json.Unmarshal(content, &c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// AnalyzePhoto analyizes a JPEG files, including the Exif metadata.
func AnalyzePhoto(path string, info os.FileInfo) (Photo, error) {
	photo := Photo{Path: path, Size: info.Size()}
	f, err := os.Open(path)
	if err != nil {
		//log.Printf("Error opening %s: %s", path, err.Error())
		return photo, err
	}
	defer f.Close()
	x, err := exif.Decode(f)
	if err == nil {
		tm, err := x.DateTime()
		if err == nil {
			photo.Timestamp = tm.Unix()
		}
		camModel, err := x.Get(exif.Model)
		if err == nil && camModel != nil {
			model, modelErr := camModel.StringVal()
			if modelErr == nil {
				photo.Camera = model
			}
		}
	}
	// The ideal hash is camera + timestamp
	if photo.Timestamp != 0 && photo.Camera != "" {
		photo.Hash = strconv.FormatInt(photo.Timestamp, 10) + "|" + photo.Camera
	} else {
		// If that doesn't work, try with the file MD5
		photo.Hash, err = utils.MD5(photo.Path)
		if err != nil {
			// In case of MD5 error, use the file name as hash
			photo.Hash = filepath.Base(photo.Path)
		}
	}
	return photo, nil
}

type workerInput struct {
	path string
	info os.FileInfo
}

type workerOutput struct {
	photo Photo
	err   error
}

func workerAnalyzePhoto(id int, jobs <-chan workerInput, results chan<- workerOutput) {
	for j := range jobs {
		photo, err := AnalyzePhoto(j.path, j.info)
		results <- workerOutput{photo, err}
	}
}

// AnalyzeDir fills the cache with data about the JPEG images contained in the
// specified directory.
func (myCache *Cache) AnalyzeDir(dir string, numWorkers int) error {
	var inputs []workerInput
	err := filepath.Walk(dir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			lowerPath := strings.ToLower(path)
			if strings.HasSuffix(lowerPath, ".jpg") || strings.HasSuffix(lowerPath, ".jpeg") || strings.HasSuffix(lowerPath, ".heic") {
				inputs = append(inputs, workerInput{path, info})
			}
			return nil
		})
	if err != nil {
		return err
	}
	numJobs := len(inputs)
	jobs := make(chan workerInput, numJobs)
	results := make(chan workerOutput, numJobs)
	for w := 0; w < numWorkers; w++ {
		go workerAnalyzePhoto(w, jobs, results)
	}
	for j := 0; j < numJobs; j++ {
		jobs <- inputs[j]
	}
	close(jobs)
	for a := 0; a < numJobs; a++ {
		output := <-results
		if output.err != nil {
			//return output.err
			continue
		}
		myCache.Photos = append(myCache.Photos, output.photo)
	}
	return nil
}

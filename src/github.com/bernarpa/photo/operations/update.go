package operations

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/bernarpa/photo/cache"
	"github.com/bernarpa/photo/config"
	"github.com/bernarpa/photo/exiftool"
	"github.com/bernarpa/photo/ssh"
	"github.com/bernarpa/photo/utils"
)

// ShowHelpUpdate prints the help for the update operation.
func ShowHelpUpdate() {
	fmt.Println()
	fmt.Println("Usage: photo update <TARGET>")
	fmt.Println()
	fmt.Println("   TARGET     one of the targets defined in config.json")
	fmt.Println()
}

func sshUpdate(conf *config.Config, target *config.Target) {
	// SSH connection
	client, _, err := ssh.Connect(target)
	if err != nil {
		log.Fatal("SSH connection error: " + err.Error())
	}
	// Ensures that the remote working dir exists
	cmdEnsureWorkDir := fmt.Sprintf("test -d '%s' || mkdir -p '%s'", target.WorkDir, target.WorkDir)
	ssh.Exec(client, cmdEnsureWorkDir)
	// Copies config.json to the remote work dir
	exePath := utils.GetExePath()
	localConfig := filepath.Join(exePath, "config.json")
	remoteConfig := target.WorkDir + "config.json"
	ssh.Copy(client, localConfig, remoteConfig)
	// Copies the exe file to the remote work dir
	localExe := filepath.Join(exePath, target.SSHExe)
	remoteExe := target.WorkDir + target.SSHExe
	ssh.Copy(client, localExe, remoteExe)
	// Ensures that the exe file is executable
	ssh.Exec(client, fmt.Sprintf("chmod +x '%s'", remoteExe))
	// Create the exiftool directory structure
	ssh.Exec(client, fmt.Sprintf("mkdir -p '%s'", strings.ReplaceAll(filepath.Join(target.WorkDir, "exiftool", "lib", "File"), conf.PathSeparator, target.SSHPathSeparator)))
	ssh.Exec(client, fmt.Sprintf("mkdir -p '%s'", strings.ReplaceAll(filepath.Join(target.WorkDir, "exiftool", "lib", "Image", "ExifTool", "Charset"), conf.PathSeparator, target.SSHPathSeparator)))
	ssh.Exec(client, fmt.Sprintf("mkdir -p '%s'", strings.ReplaceAll(filepath.Join(target.WorkDir, "exiftool", "lib", "Image", "ExifTool", "Lang"), conf.PathSeparator, target.SSHPathSeparator)))
	localExiftoolDir := filepath.Join(exePath, "exiftool")
	err = filepath.Walk(localExiftoolDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			remotePath := strings.TrimPrefix(path, exePath)
			remotePath = strings.TrimPrefix(remotePath, conf.PathSeparator)
			remotePath = strings.ReplaceAll(remotePath, conf.PathSeparator, target.SSHPathSeparator)
			remotePath = target.WorkDir + remotePath
			ssh.Copy(client, path, remotePath)
		}
		return nil
	})
	if err != nil {
		log.Fatal(fmt.Sprintf("error walking the path %s: %s\n", localExiftoolDir, err.Error()))
	}
	// Runs photo localupdate TARGET on the SSH server
	ssh.Exec(client, fmt.Sprintf("'%s' localupdate %s", remoteExe, target.Name))
	// Downloads the newly generated cache
	out := ssh.Exec(client, fmt.Sprintf("cat '%s'", target.GetRemoteCachePath()))
	localCache := target.GetLocalCachePath()
	err = ioutil.WriteFile(localCache, out, 0644)
	if err != nil {
		log.Fatal("Remote cache download error: " + err.Error())
	}
}

// LocalUpdate updates the cache for a local target.
func LocalUpdate(conf *config.Config, target *config.Target) {
	et, err := exiftool.Create(conf, target)
	if err != nil {
		log.Printf("exiftool instantation error: %s\n", err.Error())
		return
	}
	log.Printf("exiftool created: %s\n", et.Perl)
	myCache := cache.Create(target)
	for _, targetDir := range target.Collections {
		err := myCache.AnalyzeDir(targetDir, conf.Workers, et, target.Ignore)
		if err != nil {
			log.Fatal("Cache update failure: " + err.Error())
		}
	}
	jsonContent, err := json.Marshal(myCache)
	localCacheFileName := target.GetLocalCachePath()
	f, err := os.Create(localCacheFileName)
	if err != nil {
		log.Fatal("Cache file creation error: " + err.Error())
	}
	w := gzip.NewWriter(f)
	defer w.Close()
	_, err = w.Write(jsonContent)
	if err != nil {
		log.Fatal("Cache file writing error: " + err.Error())
	}
}

// Update the cache for the target specified on the command line.
func Update(conf *config.Config, target *config.Target) {
	if target.TargetType == "local" {
		LocalUpdate(conf, target)
	} else if target.TargetType == "ssh" {
		sshUpdate(conf, target)
	} else {
		log.Fatal("Unsupported target type: " + target.TargetType)
	}
}

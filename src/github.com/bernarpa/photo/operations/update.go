package operations

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/bernarpa/photo/cache"
	"github.com/bernarpa/photo/config"
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

func sshUpdate(target *config.Target) {
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
	cmdMakePhotoExecutable := fmt.Sprintf("chmod +x '%s'", remoteExe)
	ssh.Exec(client, cmdMakePhotoExecutable)
	// Runs photo localupdate TARGET on the SSH server
	cmdRemoteUpdate := fmt.Sprintf("'%s' localupdate %s", remoteExe, target.Name)
	ssh.Exec(client, cmdRemoteUpdate)
	// Downloads the newly generated cache
	remoteCache := target.GetRemoteCachePath()
	cmdCatRemoteCache := fmt.Sprintf("cat '%s'", remoteCache)
	out := ssh.Exec(client, cmdCatRemoteCache)
	localCache := target.GetLocalCachePath()
	err = ioutil.WriteFile(localCache, out, 0644)
	if err != nil {
		log.Fatal("Remote cache download error: " + err.Error())
	}
}

// LocalUpdate updates the cache for a local target.
func LocalUpdate(config *config.Config, target *config.Target) {
	myCache := cache.Create(target)
	for _, targetDir := range target.Collections {
		err := myCache.AnalyzeDir(targetDir, config.Workers)
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
		sshUpdate(target)
	} else {
		log.Fatal("Unsupported target type: " + target.TargetType)
	}
}

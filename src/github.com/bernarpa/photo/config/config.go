package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/bernarpa/photo/utils"
)

// Config is the struct that represents a Photo config.json file.
type Config struct {
	Workers       int      `json:"workers"`
	Targets       []Target `json:"targets"`
	Perl          string   `json:"perl"`
	PathSeparator string   `json:"path_separator"`
}

// Target is a photo collection to be manage through Photo. it can be local or accessible via SSH.
type Target struct {
	Name             string   `json:"name"`
	TargetType       string   `json:"target_type"`
	WorkDir          string   `json:"work_dir"`
	Perl             string   `json:"perl"`
	SSHPathSeparator string   `json:"ssh_path_separator"`
	SSHExe           string   `json:"ssh_exe"`
	SSHHost          string   `json:"ssh_host"`
	SSHPort          string   `json:"ssh_port"`
	SSHUser          string   `json:"ssh_user"`
	SSHPassword      string   `json:"ssh_password"`
	Collections      []string `json:"collections"`
	Cameras          []string `json:"cameras"`
	Ignore           []string `json:"ignore"`
}

// Load reads the content of the config.json file that should be in the same directory
// than the executable file, or Exit(1) if something goes wrong.
func Load() (*Config, error) {
	exePath := utils.GetExePath()
	configFile := filepath.Join(exePath, "config.json")
	f, err := os.Open(configFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	byteValue, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	var c Config
	err = json.Unmarshal(byteValue, &c)
	if err != nil {
		return nil, err
	}
	// Set the default Perl interpreter for local targets for
	// which a specific interpreter isn't configured.
	for i := range c.Targets {
		if c.Targets[i].Perl == "" {
			c.Targets[i].Perl = c.Perl
		}
	}
	return &c, nil
}

// GetTarget returns the specified target configuration or nil if it doesn't exist.
func (c *Config) GetTarget(name string) *Target {
	for _, t := range c.Targets {
		if t.Name == name {
			return &t
		}
	}
	return nil
}

// GetRemoteCachePath returns the cache filename on the filesystem of the target.
func (t *Target) GetRemoteCachePath() string {
	// ACHTUNG! This means that it's mandatory to have a trailing \ or / in work_dir
	return t.WorkDir + t.Name + "_cache.json.gz"
}

// GetLocalCachePath returns the cache filename on the executable filesystem.
func (t *Target) GetLocalCachePath() string {
	exePath := utils.GetExePath()
	return filepath.Join(exePath, t.Name+"_cache.json.gz")
}

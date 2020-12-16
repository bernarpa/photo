package operations

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/bernarpa/photo/cache"
	"github.com/bernarpa/photo/config"
)

type helpFunction func()
type commandFunction func(*config.Config, *config.Target)

// RunCommandFunction executes the given function by passing the
// configuration and the target specified in the command line as
// parameters.
func RunCommandFunction(cmd commandFunction, help helpFunction, requiresTarget bool) {
	if requiresTarget && len(os.Args) < 3 {
		help()
		os.Exit(1)
	}
	start := time.Now()
	config, err := config.Load()
	if err != nil {
		log.Fatal("Error while loading the configuration file")
	}
	if requiresTarget {
		targetName := os.Args[2]
		target := config.GetTarget(targetName)
		if target == nil {
			log.Fatal("Target not found: " + targetName)
		}
		cmd(config, target)
	} else {
		cmd(config, nil)
	}
	duration := time.Since(start)
	fmt.Printf("%f minutes elapsed\n", duration.Minutes())
}

func loadLocalCache(conf *config.Config, target *config.Target) *cache.Cache {
	myCache, err := cache.Load(conf, target)
	if err != nil {
		fmt.Println("Cannot load local cache, performing update...")
		Update(conf, target)
		myCache, err = cache.Load(conf, target)
		if err != nil {
			log.Fatal("Error while updating cache: " + err.Error())
		}
		return myCache
	}
	now := time.Now().Unix()
	if now-myCache.LastUpdate > 86400 {
		fmt.Println("Local cache is older than 1 day, performing update...")
		Update(conf, target)
		myCache, err = cache.Load(conf, target)
		if err != nil {
			log.Fatal("Error while updating cache: " + err.Error())
		}
	}
	return myCache
}

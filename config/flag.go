package config

import (
	"flag"
)

var (
	configFileFlagName = "cf"
	configFileFromFlag string
)

func init() {
	flag.StringVar(&configFileFromFlag, configFileFlagName, "/test/Config.yaml", "Config file")
	// flag.Set(configFileFlagName, "/anywhere/Config.yaml")
}

package config

import (
	"flag"
)

var (
	configFileFlagName = "cf"
	configFileFromFlag string
)

func init() {
	flag.StringVar(&configFileFromFlag, configFileFlagName, "/test/config.yaml", "config file")
	// flag.Set(configFileFlagName, "/anywhere/config.yaml")
}

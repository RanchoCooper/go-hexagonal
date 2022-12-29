package config

import (
	"flag"
)

/**
 * @author Rancho
 * @date 2021/12/25
 */

var (
	configFileFlagName = "cf"
	configFileFromFlag string
)

func init() {
	flag.StringVar(&configFileFromFlag, configFileFlagName, "/test/config.yaml", "config file")
	// flag.Set(configFileFlagName, "/anywhere/config.yaml")
}

package config

import (
    "flag"
)

/**
 * @author Rancho
 * @date 2021/12/25
 */

var (
    configPathFromFlag string
)

func init() {
    flag.StringVar(&configPathFromFlag, "cf", "./config/config.yaml", "path of config file")
}

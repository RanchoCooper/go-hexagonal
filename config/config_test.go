package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	// just run and see check the output information
	configPath := flag.String("config-path", "./", "path to configuration path")
	configFile := flag.String("config-file", "config.yaml", "name of configuration file (without extension)")
	flag.Parse()

	conf, err := Load(*configPath, *configFile)
	assert.NoError(t, err)
	assert.NotNil(t, conf)
	assert.True(t, conf.App.Debug)
	assert.False(t, conf.HTTPServer.Pprof)
	bytes, err := json.MarshalIndent(conf, "", "  ")
	assert.NoError(t, err)
	fmt.Println(string(bytes))
}

package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

// TestConfigEnvOverrides tests that environment variables correctly override config values
func TestConfigEnvOverrides(t *testing.T) {
	// Setup test environment variables
	os.Setenv("APP_ENV", "prod")
	os.Setenv("APP_APP_NAME", "test-app")
	os.Setenv("APP_APP_DEBUG", "false")
	os.Setenv("APP_HTTP_SERVER_ADDR", ":4000")
	os.Setenv("APP_MYSQL_HOST", "test-mysql-host")
	os.Setenv("APP_MYSQL_PORT", "3307")
	os.Setenv("APP_REDIS_HOST", "test-redis-host")
	os.Setenv("APP_LOG_COMPRESS", "true")

	// Load config
	conf, err := Load("./", "config.yaml")
	assert.NoError(t, err)

	// Clean up environment variables after test
	defer func() {
		os.Unsetenv("APP_ENV")
		os.Unsetenv("APP_APP_NAME")
		os.Unsetenv("APP_APP_DEBUG")
		os.Unsetenv("APP_HTTP_SERVER_ADDR")
		os.Unsetenv("APP_MYSQL_HOST")
		os.Unsetenv("APP_MYSQL_PORT")
		os.Unsetenv("APP_REDIS_HOST")
		os.Unsetenv("APP_LOG_COMPRESS")
	}()

	// Verify environment variables were applied correctly
	assert.Equal(t, Env("prod"), conf.Env)
	assert.Equal(t, "test-app", conf.App.Name)
	assert.False(t, conf.App.Debug)
	assert.Equal(t, ":4000", conf.HTTPServer.Addr)
	assert.Equal(t, "test-mysql-host", conf.MySQL.Host)
	assert.Equal(t, 3307, conf.MySQL.Port)
	assert.Equal(t, "test-redis-host", conf.Redis.Host)
	assert.True(t, conf.Log.Compress)
}

// TestConfigWatchChanges tests the config file change monitoring feature
func TestConfigWatchChanges(t *testing.T) {
	// Create a temporary directory for the test config file
	tempDir, err := os.MkdirTemp("", "config-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a temporary config file
	configContent := `
env: "test"
app:
  name: "test-app"
  version: "v0.1.0"
  debug: true
http_server:
  addr: ":3000"
  pprof: false
  default_page_size: 10
  max_page_size: 100
mysql:
  host: "localhost"
  port: 3306
`
	configFile := filepath.Join(tempDir, "config.yaml")
	err = os.WriteFile(configFile, []byte(configContent), 0644)
	require.NoError(t, err)

	// Load the config
	GlobalConfig = nil // Reset global config
	conf, err := Load(tempDir, "config")
	require.NoError(t, err)
	require.NotNil(t, conf)
	GlobalConfig = conf

	// Verify initial config values
	assert.Equal(t, Env("test"), conf.Env)
	assert.Equal(t, "test-app", conf.App.Name)
	assert.True(t, conf.App.Debug)
	assert.Equal(t, ":3000", conf.HTTPServer.Addr)

	// Modify the config file
	updatedConfigContent := `
env: "prod"
app:
  name: "updated-app"
  version: "v0.2.0"
  debug: false
http_server:
  addr: ":4000"
  pprof: true
  default_page_size: 20
  max_page_size: 200
mysql:
  host: "updated-host"
  port: 3307
`
	// Write the updated config to the file
	err = os.WriteFile(configFile, []byte(updatedConfigContent), 0644)
	require.NoError(t, err)

	// Wait for file system events to propagate
	time.Sleep(500 * time.Millisecond)

	// Safely access the config with the read lock
	configMutex.RLock()
	updatedConfig := GlobalConfig
	configMutex.RUnlock()

	// Verify the config was updated
	assert.Equal(t, Env("prod"), updatedConfig.Env)
	assert.Equal(t, "updated-app", updatedConfig.App.Name)
	assert.False(t, updatedConfig.App.Debug)
	assert.Equal(t, ":4000", updatedConfig.HTTPServer.Addr)
	assert.Equal(t, "updated-host", updatedConfig.MySQL.Host)
	assert.Equal(t, 3307, updatedConfig.MySQL.Port)

	// Check that the last change time was updated
	assert.NotEqual(t, time.Time{}, GetLastConfigChangeTime())
}

// TestGetDuration tests the GetDuration function
func TestGetDuration(t *testing.T) {
	tests := []struct {
		input    string
		expected time.Duration
	}{
		{"10s", 10 * time.Second},
		{"5m", 5 * time.Minute},
		{"2h", 2 * time.Hour},
		{"1h30m", 90 * time.Minute},
		{"0", 0},
		{"", 0},
		{"invalid", 0},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result := GetDuration(test.input)
			assert.Equal(t, test.expected, result)
		})
	}
}

// TestEnvIsProd tests the IsProd method of the Env type
func TestEnvIsProd(t *testing.T) {
	tests := []struct {
		env      Env
		expected bool
	}{
		{"prod", true},
		{"production", false},
		{"dev", false},
		{"test", false},
		{"", false},
	}

	for _, test := range tests {
		t.Run(string(test.env), func(t *testing.T) {
			result := test.env.IsProd()
			assert.Equal(t, test.expected, result)
		})
	}
}

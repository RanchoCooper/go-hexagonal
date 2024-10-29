package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"go-hexagonal/util"
)

const (
	configFilePath        = "config.yaml"
	privateConfigFilePath = "config.private.yaml"
)

var Config = &config{}

type config struct {
	Env        Env               `yaml:"env"`
	App        *appConfig        `yaml:"app"`
	HTTPServer *httpServerConfig `yaml:"http_server"`
	Log        *logConfig        `yaml:"log"`
	MySQL      *MySQLConfig      `yaml:"mysql"`
	Redis      *RedisConfig      `yaml:"redis"`
}

type appConfig struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
	Debug   bool   `yaml:"debug"`
}

type httpServerConfig struct {
	Addr            string `yaml:"addr"`
	Pprof           bool   `yaml:"pprof"`
	DefaultPageSize int    `yaml:"default_page_size"`
	MaxPageSize     int    `yaml:"max_page_size"`
	ReadTimeout     string `yaml:"read_timeout"`
	WriteTimeout    string `yaml:"write_timeout"`
}

type logConfig struct {
	SavePath  string `yaml:"save_path"`
	FileName  string `yaml:"file_name"`
	MaxSize   int    `yaml:"max_size"`
	MaxAge    int    `yaml:"max_age"`
	LocalTime bool   `yaml:"local_time"`
	Compress  bool   `yaml:"compress"`
}

type MySQLConfig struct {
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	Database     string `yaml:"database"`
	MaxIdleConns int    `yaml:"max_idle_conns"`
	MaxOpenConns int    `yaml:"max_open_conns"`
	MaxLifeTime  string `yaml:"max_life_time"`
	MaxIdleTime  string `yaml:"max_idle_time"`
	CharSet      string `yaml:"char_set"`
	ParseTime    bool   `yaml:"parse_time"`
	TimeZone     string `yaml:"time_zone"`
}

type PostgresDBConf struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DbName   string `yaml:"bbName"`
	SSLMode  string `yaml:"ssl_mode"`
	TimeZone string `yaml:"time_zone"`
}

type RedisConfig struct {
	Host         string `yaml:"host"`
	UserName     string `yaml:"user_name"`
	Password     string `yaml:"password"`
	DB           int    `yaml:"db"`
	PoolSize     int    `yaml:"pool_size"`
	IdleTimeout  int    `yaml:"idle_timeout"`
	MinIdleConns int    `yaml:"min_idle_conns"`
}

func readYamlConfig(configPath string) {
	yamlFile, err := filepath.Abs(configPath)
	if err != nil {
		log.Fatalf("invalid config file path, err: %v", err)
	}
	content, err := os.ReadFile(yamlFile)
	if err != nil {
		log.Fatalf("read config file fail, err: %v", err)
	}
	err = yaml.Unmarshal(content, Config)
	if err != nil {
		log.Fatalf("config file unmarshal fail, err: %v", err)
	}
}

func Init() {
	flag.Parse()

	configPath := ""
	configFileFlagSet := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == configFileFlagName {
			configFileFlagSet = true
		}
	})
	if configFileFlagSet {
		configPath = configFileFromFlag
	} else {
		configPath = util.GetCurrentPath()
	}

	readYamlConfig(filepath.Join(configPath, configFilePath))
	if !Config.Env.IsGithub() {
		// read private sensitive configs
		readYamlConfig(filepath.Join(configPath, privateConfigFilePath))
	}

	bf, _ := json.MarshalIndent(Config, "", "    ")
	fmt.Printf("Config:\n%s\n", string(bf))
}

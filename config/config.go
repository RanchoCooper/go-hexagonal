package config

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "path/filepath"

    "gopkg.in/yaml.v3"

    "go-hexagonal/util"
)

const (
    configFilePath        = "/config.yaml"
    privateConfigFilePath = "/config.private.yaml"
    EnvLocal              = "local"
    EnvGithub             = "github"
)

var Config = &config{}

type appConfig struct {
    Name    string `yaml:"name"`
    Version string `yaml:"version"`
    Debug   bool   `yaml:"debug"`
    Env     string `yaml:"env"`
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
    SavePath string `yaml:"save_path"`
    FileName string `yaml:"file_name"`
}

type mysqlConfig struct {
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

type redisConfig struct {
    Addr         string
    UserName     string
    Password     string
    DB           int
    PoolSize     int
    IdleTimeout  int
    MinIdleConns int
}

type config struct {
    App        *appConfig        `yaml:"app"`
    HTTPServer *httpServerConfig `yaml:"http_server"`
    Log        *logConfig        `yaml:"log"`
    MySQL      *mysqlConfig      `yaml:"mysql"`
    Redis      *redisConfig      `yaml:"redis"`
}

func readYamlConfig(configPath string) {
    yamlFile, err := filepath.Abs(configPath)
    if err != nil {
        log.Fatalf("invalid config file path, err: %v", err)
    }
    content, err := ioutil.ReadFile(yamlFile)
    if err != nil {
        log.Printf("read config file fail, err: %v", err)
    }
    err = yaml.Unmarshal(content, Config)
    if err != nil {
        log.Printf("config file unmarshal fail, err: %v", err)
    }
}

func Init() {
    configPath := util.GetCurrentPath()

    readYamlConfig(configPath + configFilePath)
    // read private sensitive configs
    readYamlConfig(configPath + privateConfigFilePath)

    bf, _ := json.MarshalIndent(Config, "", "    ")
    fmt.Printf("Config:\n%s\n", string(bf))
}

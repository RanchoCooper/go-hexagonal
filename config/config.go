package config

import (
	"flag"

	"github.com/spf13/viper"
)

type Env string

var GlobalConfig *Config

type Config struct {
	Env          Env               `yaml:"env"`
	App          *AppConfig        `yaml:"app"`
	HTTPServer   *HttpServerConfig `yaml:"http_server"`
	Log          *LogConfig        `yaml:"log"`
	MySQL        *MySQLConfig      `yaml:"mysql"`
	Redis        *RedisConfig      `yaml:"redis"`
	Postgre      *PostgreSQLConfig `yaml:"postgres"`
	MigrationDir string            `yaml:"migration_dir"`
}

type AppConfig struct {
	Name    string `yaml:"name"`
	Debug   bool   `yaml:"debug"`
	Version string `yaml:"version"`
}

type HttpServerConfig struct {
	Addr            string `yaml:"addr"`
	Pprof           bool   `yaml:"pprof"`
	DefaultPageSize int    `yaml:"default_page_size"`
	MaxPageSize     int    `yaml:"max_page_size"`
	ReadTimeout     string `yaml:"read_timeout"`
	WriteTimeout    string `yaml:"write_timeout"`
}

type LogConfig struct {
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

type PostgreSQLConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DbName   string `yaml:"dbName"`
	SSLMode  string `yaml:"ssl_mode"`
	TimeZone string `yaml:"time_zone"`
}

type RedisConfig struct {
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	Password     string `yaml:"password"`
	DB           int    `yaml:"db"`
	PoolSize     int    `yaml:"poolSize"`
	IdleTimeout  int    `yaml:"idleTimeout"`
	MinIdleConns int    `yaml:"minIdleConns"`
}

func Load(configPath string, configFile string) (*Config, error) {
	var conf *Config
	vip := viper.New()
	vip.AddConfigPath(configPath)
	vip.SetConfigName(configFile)

	vip.SetConfigType("yaml")
	if err := vip.ReadInConfig(); err != nil {
		return nil, err
	}

	err := vip.Unmarshal(&conf)
	if err != nil {
		return nil, err
	}

	return conf, nil
}

func Init(path, file string) {
	configPath := flag.String("config-path", path, "path to configuration path")
	configFile := flag.String("config-file", file, "name of configuration file (without extension)")
	flag.Parse()

	conf, err := Load(*configPath, *configFile)
	if err != nil {
		panic("Load config fail : " + err.Error())
	}
	GlobalConfig = conf
}

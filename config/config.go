package config

import (
	"flag"

	"github.com/spf13/viper"
)

type Env string

func (e Env) IsProd() bool {
	return e == "prod"
}

var GlobalConfig *Config

type Config struct {
	Env          Env               `yaml:"env" mapstructure:"env"`
	App          *AppConfig        `yaml:"app" mapstructure:"app"`
	HTTPServer   *HttpServerConfig `yaml:"http_server" mapstructure:"http_server"`
	Log          *LogConfig        `yaml:"log" mapstructure:"log"`
	MySQL        *MySQLConfig      `yaml:"mysql" mapstructure:"mysql"`
	Redis        *RedisConfig      `yaml:"redis" mapstructure:"redis"`
	Postgre      *PostgreSQLConfig `yaml:"postgres" mapstructure:"postgres"`
	MigrationDir string            `yaml:"migration_dir" mapstructure:"migration_dir"`
}

type AppConfig struct {
	Name    string `yaml:"name" mapstructure:"name"`
	Debug   bool   `yaml:"debug" mapstructure:"debug"`
	Version string `yaml:"version" mapstructure:"version"`
}

type HttpServerConfig struct {
	Addr            string `yaml:"addr" mapstructure:"addr"`
	Pprof           bool   `yaml:"pprof" mapstructure:"pprof"`
	DefaultPageSize int    `yaml:"default_page_size" mapstructure:"default_page_size"`
	MaxPageSize     int    `yaml:"max_page_size" mapstructure:"max_page_size"`
	ReadTimeout     string `yaml:"read_timeout" mapstructure:"read_timeout"`
	WriteTimeout    string `yaml:"write_timeout" mapstructure:"write_timeout"`
}

type LogConfig struct {
	SavePath  string `yaml:"save_path" mapstructure:"save_path"`
	FileName  string `yaml:"file_name" mapstructure:"file_name"`
	MaxSize   int    `yaml:"max_size" mapstructure:"max_size"`
	MaxAge    int    `yaml:"max_age" mapstructure:"max_age"`
	LocalTime bool   `yaml:"local_time" mapstructure:"local_time"`
	Compress  bool   `yaml:"compress" mapstructure:"compress"`
}

type MySQLConfig struct {
	User         string `yaml:"user" mapstructure:"user"`
	Password     string `yaml:"password" mapstructure:"password"`
	Host         string `yaml:"host" mapstructure:"host"`
	Port         int    `yaml:"port" mapstructure:"port"`
	Database     string `yaml:"database" mapstructure:"database"`
	MaxIdleConns int    `yaml:"max_idle_conns" mapstructure:"max_idle_conns"`
	MaxOpenConns int    `yaml:"max_open_conns" mapstructure:"max_open_conns"`
	MaxLifeTime  string `yaml:"max_life_time" mapstructure:"max_life_time"`
	MaxIdleTime  string `yaml:"max_idle_time" mapstructure:"max_idle_time"`
	CharSet      string `yaml:"char_set" mapstructure:"char_set"`
	ParseTime    bool   `yaml:"parse_time" mapstructure:"parse_time"`
	TimeZone     string `yaml:"time_zone" mapstructure:"time_zone"`
}

type PostgreSQLConfig struct {
	Host     string `yaml:"host" mapstructure:"host"`
	Port     int    `yaml:"port" mapstructure:"port"`
	Username string `yaml:"username" mapstructure:"username"`
	Password string `yaml:"password" mapstructure:"password"`
	DbName   string `yaml:"dbName" mapstructure:"db_name"`
	SSLMode  string `yaml:"ssl_mode" mapstructure:"ssl_mode"`
	TimeZone string `yaml:"time_zone" mapstructure:"time_zone"`
}

type RedisConfig struct {
	Host         string `yaml:"host" mapstructure:"host"`
	Port         int    `yaml:"port" mapstructure:"port"`
	Password     string `yaml:"password" mapstructure:"password"`
	DB           int    `yaml:"db" mapstructure:"db"`
	PoolSize     int    `yaml:"poolSize" mapstructure:"poolSize"`
	IdleTimeout  int    `yaml:"idleTimeout" mapstructure:"idleTimeout"`
	MinIdleConns int    `yaml:"minIdleConns" mapstructure:"minIdleConns"`
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

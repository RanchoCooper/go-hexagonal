package config

import (
	"flag"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

// Constants
const (
	TrueStr = "true" // String representation of boolean true value
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

	// Enable environment variables to override config
	vip.SetEnvPrefix("APP")
	vip.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	vip.AutomaticEnv()

	err := vip.Unmarshal(&conf)
	if err != nil {
		return nil, err
	}

	// Apply environment variable overrides
	applyEnvOverrides(conf)

	return conf, nil
}

// applyEnvOverrides applies environment variable overrides to the configuration
func applyEnvOverrides(conf *Config) {
	// Apply config overrides by category
	applyAppEnvOverrides(conf)
	applyHTTPServerEnvOverrides(conf)
	applyMySQLEnvOverrides(conf)
	applyPostgresEnvOverrides(conf)
	applyRedisEnvOverrides(conf)
	applyLogEnvOverrides(conf)

	// Migration directory
	if migrationDir := os.Getenv("APP_MIGRATION_DIR"); migrationDir != "" {
		conf.MigrationDir = migrationDir
	}
}

// applyAppEnvOverrides applies App related environment variables
func applyAppEnvOverrides(conf *Config) {
	// Environment
	if env := os.Getenv("APP_ENV"); env != "" {
		conf.Env = Env(env)
	}

	// App config
	if name := os.Getenv("APP_APP_NAME"); name != "" {
		conf.App.Name = name
	}
	if debug := os.Getenv("APP_APP_DEBUG"); debug != "" {
		conf.App.Debug = debug == TrueStr
	}
	if version := os.Getenv("APP_APP_VERSION"); version != "" {
		conf.App.Version = version
	}
}

// applyHTTPServerEnvOverrides applies HTTP server related environment variables
func applyHTTPServerEnvOverrides(conf *Config) {
	if addr := os.Getenv("APP_HTTP_SERVER_ADDR"); addr != "" {
		conf.HTTPServer.Addr = addr
	}
	if pprof := os.Getenv("APP_HTTP_SERVER_PPROF"); pprof != "" {
		conf.HTTPServer.Pprof = pprof == TrueStr
	}
	if pageSize := os.Getenv("APP_HTTP_SERVER_DEFAULT_PAGE_SIZE"); pageSize != "" {
		if val, err := strconv.Atoi(pageSize); err == nil {
			conf.HTTPServer.DefaultPageSize = val
		}
	}
	if maxPageSize := os.Getenv("APP_HTTP_SERVER_MAX_PAGE_SIZE"); maxPageSize != "" {
		if val, err := strconv.Atoi(maxPageSize); err == nil {
			conf.HTTPServer.MaxPageSize = val
		}
	}
	if readTimeout := os.Getenv("APP_HTTP_SERVER_READ_TIMEOUT"); readTimeout != "" {
		conf.HTTPServer.ReadTimeout = readTimeout
	}
	if writeTimeout := os.Getenv("APP_HTTP_SERVER_WRITE_TIMEOUT"); writeTimeout != "" {
		conf.HTTPServer.WriteTimeout = writeTimeout
	}
}

// applyMySQLEnvOverrides applies MySQL related environment variables
func applyMySQLEnvOverrides(conf *Config) {
	if host := os.Getenv("APP_MYSQL_HOST"); host != "" {
		conf.MySQL.Host = host
	}
	if port := os.Getenv("APP_MYSQL_PORT"); port != "" {
		if val, err := strconv.Atoi(port); err == nil {
			conf.MySQL.Port = val
		}
	}
	if user := os.Getenv("APP_MYSQL_USER"); user != "" {
		conf.MySQL.User = user
	}
	if password := os.Getenv("APP_MYSQL_PASSWORD"); password != "" {
		conf.MySQL.Password = password
	}
	if database := os.Getenv("APP_MYSQL_DATABASE"); database != "" {
		conf.MySQL.Database = database
	}
	if maxIdleConns := os.Getenv("APP_MYSQL_MAX_IDLE_CONNS"); maxIdleConns != "" {
		if val, err := strconv.Atoi(maxIdleConns); err == nil {
			conf.MySQL.MaxIdleConns = val
		}
	}
	if maxOpenConns := os.Getenv("APP_MYSQL_MAX_OPEN_CONNS"); maxOpenConns != "" {
		if val, err := strconv.Atoi(maxOpenConns); err == nil {
			conf.MySQL.MaxOpenConns = val
		}
	}
	if maxLifeTime := os.Getenv("APP_MYSQL_MAX_LIFE_TIME"); maxLifeTime != "" {
		conf.MySQL.MaxLifeTime = maxLifeTime
	}
	if maxIdleTime := os.Getenv("APP_MYSQL_MAX_IDLE_TIME"); maxIdleTime != "" {
		conf.MySQL.MaxIdleTime = maxIdleTime
	}
	if charSet := os.Getenv("APP_MYSQL_CHAR_SET"); charSet != "" {
		conf.MySQL.CharSet = charSet
	}
	if parseTime := os.Getenv("APP_MYSQL_PARSE_TIME"); parseTime != "" {
		conf.MySQL.ParseTime = parseTime == TrueStr
	}
	if timeZone := os.Getenv("APP_MYSQL_TIME_ZONE"); timeZone != "" {
		conf.MySQL.TimeZone = timeZone
	}
}

// applyPostgresEnvOverrides applies PostgreSQL related environment variables
func applyPostgresEnvOverrides(conf *Config) {
	if host := os.Getenv("APP_POSTGRES_HOST"); host != "" {
		conf.Postgre.Host = host
	}
	if port := os.Getenv("APP_POSTGRES_PORT"); port != "" {
		if val, err := strconv.Atoi(port); err == nil {
			conf.Postgre.Port = val
		}
	}
	if username := os.Getenv("APP_POSTGRES_USERNAME"); username != "" {
		conf.Postgre.Username = username
	}
	if password := os.Getenv("APP_POSTGRES_PASSWORD"); password != "" {
		conf.Postgre.Password = password
	}
	if dbName := os.Getenv("APP_POSTGRES_DB_NAME"); dbName != "" {
		conf.Postgre.DbName = dbName
	}
	if sslMode := os.Getenv("APP_POSTGRES_SSL_MODE"); sslMode != "" {
		conf.Postgre.SSLMode = sslMode
	}
	if timeZone := os.Getenv("APP_POSTGRES_TIME_ZONE"); timeZone != "" {
		conf.Postgre.TimeZone = timeZone
	}
}

// applyRedisEnvOverrides applies Redis related environment variables
func applyRedisEnvOverrides(conf *Config) {
	if host := os.Getenv("APP_REDIS_HOST"); host != "" {
		conf.Redis.Host = host
	}
	if port := os.Getenv("APP_REDIS_PORT"); port != "" {
		if val, err := strconv.Atoi(port); err == nil {
			conf.Redis.Port = val
		}
	}
	if password := os.Getenv("APP_REDIS_PASSWORD"); password != "" {
		conf.Redis.Password = password
	}
	if db := os.Getenv("APP_REDIS_DB"); db != "" {
		if val, err := strconv.Atoi(db); err == nil {
			conf.Redis.DB = val
		}
	}
	if poolSize := os.Getenv("APP_REDIS_POOL_SIZE"); poolSize != "" {
		if val, err := strconv.Atoi(poolSize); err == nil {
			conf.Redis.PoolSize = val
		}
	}
	if idleTimeout := os.Getenv("APP_REDIS_IDLE_TIMEOUT"); idleTimeout != "" {
		if val, err := strconv.Atoi(idleTimeout); err == nil {
			conf.Redis.IdleTimeout = val
		}
	}
	if minIdleConns := os.Getenv("APP_REDIS_MIN_IDLE_CONNS"); minIdleConns != "" {
		if val, err := strconv.Atoi(minIdleConns); err == nil {
			conf.Redis.MinIdleConns = val
		}
	}
}

// applyLogEnvOverrides applies Log related environment variables
func applyLogEnvOverrides(conf *Config) {
	if savePath := os.Getenv("APP_LOG_SAVE_PATH"); savePath != "" {
		conf.Log.SavePath = savePath
	}
	if fileName := os.Getenv("APP_LOG_FILE_NAME"); fileName != "" {
		conf.Log.FileName = fileName
	}
	if maxSize := os.Getenv("APP_LOG_MAX_SIZE"); maxSize != "" {
		if val, err := strconv.Atoi(maxSize); err == nil {
			conf.Log.MaxSize = val
		}
	}
	if maxAge := os.Getenv("APP_LOG_MAX_AGE"); maxAge != "" {
		if val, err := strconv.Atoi(maxAge); err == nil {
			conf.Log.MaxAge = val
		}
	}
	if localTime := os.Getenv("APP_LOG_LOCAL_TIME"); localTime != "" {
		conf.Log.LocalTime = localTime == TrueStr
	}
	if compress := os.Getenv("APP_LOG_COMPRESS"); compress != "" {
		conf.Log.Compress = compress == TrueStr
	}
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

// GetDuration converts a duration string to time.Duration
func GetDuration(durationStr string) time.Duration {
	return cast.ToDuration(durationStr)
}

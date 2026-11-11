package config

import (
	"log/slog"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

type Config struct {
	DataSource string `yaml:"data_source"`
	Server     Server `yaml:"server"`
	// MySQL      MySQL      `yaml:"mysql"`
	Redis Redis `yaml:"redis"`
	CSV   CSV   `yaml:"csv"`
	// GRPC       GRPC       `yaml:"grpc"`
	Middleware Middleware `yaml:"middleware"`
}

type Server struct {
	// BasePath is prefix path for routing (e.g. /v1, /api, or empty for root)
	BasePath string `yaml:"base_path"`

	// Port is the local machine TCP Port to bind the HTTP Server to
	Port string `yaml:"port"`

	// Prefork will spawn multiple Go processes listening on the same port
	Prefork bool `yaml:"prefork"`

	// StrictRouting
	// When enabled, the router treats /foo and /foo/ as different.
	// Otherwise, the router treats /foo and /foo/ as the same.
	StrictRouting bool `yaml:"strict_routing"`

	// CaseSensitive
	// When enabled, /Foo and /foo are different routes.
	// When disabled, /Foo and /foo are treated the same.
	CaseSensitive bool `yaml:"case_sensitive"`

	// BodyLimit
	// Sets the maximum allowed size for a request body, if the size exceeds
	// the configured limit, it sends 413 - Request Entity Too Large response.
	BodyLimit int `yaml:"body_limit"`

	// Concurrency maximum number of concurrent connections
	Concurrency int `yaml:"concurrency"`

	Timeout Timeout `yaml:"timeout"`

	// LogLevel is log level, available value: error, warning, info, debug
	LogLevel string `yaml:"log_level" mapstructure:"log_level"`

	// LogType is log format, available value: text, json (default: json)
	LogType string `yaml:"log_type" mapstructure:"log_type"`

	// GRPCPort is the local machine TCP port to bind the gRPC server to
	GRPCPort string `yaml:"grpc_port"`
}

type Timeout struct {
	// Read is the amount of time to wait until an HTTP server
	// read operation is cancelled
	Read time.Duration `yaml:"read"`

	// Write is the amount of time to wait until an HTTP server
	// write opperation is cancelled
	Write time.Duration `yaml:"write"`

	// Read is the amount of time to wait
	// until an IDLE HTTP session is closed
	Idle time.Duration `yaml:"idle"`
}

// type MySQL struct {
// 	// Host is the MySQL IP Address to connect to
// 	Host string `yaml:"host,omitempty"`

// 	// Port is the MySQL Port to connect to
// 	Port string `yaml:"port,omitempty"`

// 	// Database is MySQL database name
// 	Database string `yaml:"database"`

// 	// User is MySQL username
// 	User string `yaml:"user"`

// 	// Password is MySQL password
// 	Password string `yaml:"password"`

// 	// PathMigrate is directory for migration file
// 	PathMigrate string `yaml:"path_migrate"`
// }

type Redis struct {
	// Host is the Redis IP Address to connect to
	Host string `yaml:"host,omitempty"`

	// Port is the Redis Port to connect to
	Port string `yaml:"port,omitempty"`

	// MaxActive is Redis maximum connection
	MaxConnection uint64 `yaml:"max_connection"`

	// Username
	Username string `yaml:"username"`

	// Password
	Password string `yaml:"password"`

	// Database
	Database uint64 `yaml:"database"`
}

// type GRPC struct {
// 	Host            string        `yaml:"host"`
// 	Port            string        `yaml:"port"`
// 	Init            int           `yaml:"init"`
// 	Capacity        int           `yaml:"capacity"`
// 	IdleDuration    time.Duration `yaml:"idle_duration"`
// 	MaxLifeDuration time.Duration `yaml:"max_life_duration"`
// }

type CSV struct {
	// Path is the directory location containing CSV files (cities.csv, praytimes/{city_code}.csv)
	Path string `yaml:"path"`
}

type Middleware struct {
	AllowsOrigin string `yaml:"allows_origin"`
}

// Default config
var defaultConfig = &Config{
	Server: Server{
		BasePath:      "",
		Port:          "4321",
		Prefork:       false,
		StrictRouting: false,
		CaseSensitive: false,
		BodyLimit:     4 * 1024 * 1024,
		Concurrency:   256 * 1024,
		Timeout: Timeout{
			Read:  5,
			Write: 10,
			Idle:  120,
		},
		LogLevel: "debug",
		LogType:  "json",
		GRPCPort: "54321",
	},
	// MySQL: MySQL{
	// 	Host:        "localhost",
	// 	Port:        "3306",
	// 	Database:    "test",
	// 	User:        "root",
	// 	Password:    "root",
	// 	PathMigrate: "file:///migration",
	// },
	Redis: Redis{
		Host:          "localhost",
		Port:          "6379",
		MaxConnection: 80,
		Username:      "",
		Password:      "",
		Database:      0,
	},
	CSV: CSV{
		Path: "./data",
	},
	// GRPC: GRPC{
	// 	Host:            "localhost",
	// 	Port:            "59090",
	// 	Init:            5,
	// 	Capacity:        50,
	// 	IdleDuration:    60,
	// 	MaxLifeDuration: 60,
	// },
	Middleware: Middleware{
		AllowsOrigin: "*",
	},
	DataSource: "redis",
}

// registerDefaults recursively registers all leaf fields from defaultConfig into viper
func registerDefaults(prefix string, rv reflect.Value) {
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	rt := rv.Type()
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		tag := strings.Split(field.Tag.Get("yaml"), ",")[0]
		if tag == "" || tag == "-" {
			continue
		}
		key := tag
		if prefix != "" {
			key = prefix + "." + tag
		}

		fieldVal := rv.Field(i)
		if fieldVal.Kind() == reflect.Struct {
			registerDefaults(key, fieldVal)
		} else {
			viper.SetDefault(key, fieldVal.Interface())
		}
	}
}

func ReadConfig(configFile string) {
	// 1. Register defaults for all leaf keys
	registerDefaults("", reflect.ValueOf(defaultConfig))

	// 2. Enable environment variable mapping (e.g. SERVER_LOG_LEVEL -> server.log_level)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// 3. Resolve config file
	if configFile != "" {
		if _, err := os.Stat(configFile); err == nil {
			viper.SetConfigFile(configFile)
		} else {
			viper.SetConfigName(configFile)
			viper.SetConfigType("yaml")
			viper.AddConfigPath(".")
			viper.AddConfigPath("/usr/local/etc/")
			viper.AddConfigPath("/")
		}

		if err := viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok || os.IsNotExist(err) {
				slog.Info("Config file not found, using default configuration with environment overrides", "config_file", configFile)
			} else {
				slog.Error("Failed to read config file", "error", err, "config_file", configFile)
				os.Exit(1)
			}
		} else {
			slog.Info("Config loaded successfully", "config_file", viper.ConfigFileUsed())
		}
	}

	c := viper.AllSettings()
	bs, err := yaml.Marshal(c)
	if err != nil {
		slog.Error("Unable to marshal config to YAML", "error", err)
	}
	slog.Debug("Resolved configuration settings", "config_yaml", string(bs))
}

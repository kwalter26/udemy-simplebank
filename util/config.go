package util

import (
	"fmt"
	"github.com/spf13/viper"
	"path/filepath"
	"reflect"
	"time"
)

type Config struct {
	DBDriver                          string        `mapstructure:"DB_DRIVER"`
	DBSource                          string        `mapstructure:"DB_SOURCE"`
	HttpServerAddress                 string        `mapstructure:"HTTP_SERVER_ADDRESS"`
	GrpcServerAddress                 string        `mapstructure:"GRPC_SERVER_ADDRESS"`
	TokenSymmetricKey                 string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration               time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	NewRelicAppName                   string        `mapstructure:"NEWRELIC_APP_NAME"`
	NewRelicLicenseKey                string        `mapstructure:"NEWRELIC_LICENSE_KEY"`
	NewRelicLogForwardingEnabled      bool          `mapstructure:"NEWRELIC_LOG_FORWARDING_ENABLED"`
	NewRelicDistributedTracingEnabled bool          `mapstructure:"NEWRELIC_DIST_TRACING_ENABLED"`
	NewRelicAppEnabled                bool          `mapstructure:"NEWRELIC_APP_ENABLED"`
	RefreshTokenDuration              time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
}

type Environment string

const (
	Local Environment = "local"
	Prod  Environment = "prod"
)

func LoadConfig(path string, env Environment) (config Config, err error) {
	v := viper.New()
	c, err2 := BindEnv(config, v)
	if err2 != nil {
		return c, err2
	}
	v.AutomaticEnv()
	absPath, err := filepath.Abs(path)
	filename := string(env) + ".env"

	v.AddConfigPath(absPath)
	v.SetConfigName(filename)
	v.SetConfigType("env")

	err = v.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Printf("Config file not found: %s/%s\n", absPath, filename)
		} else {
			fmt.Printf("Config file was found but another error was produced: %s/%s\n", absPath, filename)
			return Config{}, err
		}
	}

	err = v.Unmarshal(&config)
	return config, err
}

func BindEnv(config Config, v *viper.Viper) (Config, error) {
	keys := reflect.ValueOf(&config).Elem()
	for i := 0; i < keys.NumField(); i++ {
		key := keys.Type().Field(i).Tag.Get("mapstructure")
		err := v.BindEnv(key)
		if err != nil {
			return config, err
		}
	}
	return Config{}, nil
}

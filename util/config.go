package util

import (
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	DBDriver                          string        `mapstructure:"DB_DRIVER"`
	DBSource                          string        `mapstructure:"DB_SOURCE"`
	ServerAddress                     string        `mapstructure:"SERVER_ADDRESS"`
	TokenSymmetricKey                 string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration               time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	NewRelicAppName                   string        `mapstructure:"NEWRELIC_APP_NAME"`
	NewRelicLicenseKey                string        `mapstructure:"NEWRELIC_LICENSE_KEY"`
	NewRelicLogForwardingEnabled      bool          `mapstructure:"NEWRELIC_LOG_FORWARDING_ENABLED"`
	NewRelicDistributedTracingEnabled bool          `mapstructure:"NEWRELIC_DIST_TRACING_ENABLED"`
	NewRelicAppEnabled                bool          `mapstructure:"NEWRELIC_APP_ENABLED"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}

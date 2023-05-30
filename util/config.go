package util

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
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

type Environment string

const (
	Local Environment = "local"
	Prod  Environment = "prod"
)

func LoadConfig(path string, env Environment) (config Config, err error) {
	viper.AutomaticEnv()
	absPath, err := filepath.Abs(path)
	filename := string(env) + ".env"

	filePath := absPath + "/" + filename
	if _, err := os.Stat(filePath); err == nil {
		fmt.Printf("Found file '%s'. Loading variables...", filePath)
		viper.AddConfigPath(absPath)
		viper.SetConfigName(filename)
		viper.SetConfigType("env")
		err = viper.ReadInConfig()
		if err != nil {
			return Config{}, err
		}
	} else {
		fmt.Printf("%s does not exist\n", filePath)
	}

	err = viper.Unmarshal(&config)
	return config, err
}

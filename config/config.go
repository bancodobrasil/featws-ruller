package config

import (
	"os"
	"strings"

	"github.com/spf13/viper"
)

//Config ...
type Config struct {
	ResourceLoaderType       string `mapstructure:"FEATWS_RULLER_RESOURCE_LOADER_TYPE"`
	ResourceLoaderURL        string `mapstructure:"FEATWS_RULLER_RESOURCE_LOADER_URL"`
	ResourceLoaderHeaders    map[string]string
	ResourceLoaderHeadersStr string `mapstructure:"FEATWS_RULLER_RESOURCE_LOADER_HEADERS"`
	Port                     string `mapstructure:"FEATWS_RULLER_PORT"`
	DefaultRules             string `mapstructure:"FEATWS_RULLER_DEFAULT_RULES"`
}

func LoadConfig(config *Config) (err error) {
	viper.AddConfigPath("./")
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	viper.SetDefault("FEATWS_RULLER_RESOURCE_LOADER_TYPE", "http")
	viper.SetDefault("FEATWS_RULLER_RESOURCE_LOADER_URL", "")
	viper.SetDefault("FEATWS_RULLER_RESOURCE_LOADER_HEADERS", "map[string]string{}")
	viper.SetDefault("FEATWS_RULLER_DEFAULT_RULES", "")
	viper.SetDefault("FEATWS_RULLER_PORT", "8000")

	err = viper.ReadInConfig()
	if err != nil {
		if err2, ok := err.(*os.PathError); !ok {
			err = err2
			return
		}
	}

	err = viper.Unmarshal(config)

	config.ResourceLoaderHeaders = map[string]string{}

	headers := strings.Split(config.ResourceLoaderHeadersStr, ",")

	for _, value := range headers {
		entries := strings.Split(value, ":")
		if len(entries) == 2 {
			config.ResourceLoaderHeaders[entries[0]] = entries[1]
		}

	}

	return
}

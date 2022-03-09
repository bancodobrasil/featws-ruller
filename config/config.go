package config

import (
	"net/http"
	"os"
	"strings"

	"github.com/spf13/viper"
)

//Config ...
type Config struct {
	ResourceLoaderType       string `mapstructure:"FEATWS_RULLER_RESOURCE_LOADER_TYPE"`
	ResourceLoaderURL        string `mapstructure:"FEATWS_RULLER_RESOURCE_LOADER_URL"`
	ResourceLoaderHeaders    http.Header
	ResourceLoaderHeadersStr string `mapstructure:"FEATWS_RULLER_RESOURCE_LOADER_HEADERS"`

	Port         string `mapstructure:"PORT"`
	DefaultRules string `mapstructure:"FEATWS_RULLER_DEFAULT_RULES"`

	ResolverBridgeURL        string `mapstructure:"FEATWS_RULLER_RESOLVER_BRIDGE_URL"`
	ResolverBridgeHeaders    http.Header
	ResolverBridgeHeadersStr string `mapstructure:"FEATWS_RULLER_RESOLVER_BRIDGE_HEADERS"`
}

var config = &Config{}

//LoadConfig ...
func LoadConfig() (err error) {
	viper.AddConfigPath("./")
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	viper.SetDefault("FEATWS_RULLER_RESOURCE_LOADER_TYPE", "http")
	viper.SetDefault("FEATWS_RULLER_RESOURCE_LOADER_URL", "")
	viper.SetDefault("FEATWS_RULLER_RESOURCE_LOADER_HEADERS", "")
	viper.SetDefault("FEATWS_RULLER_RESOLVER_BRIDGE_URL", "")
	viper.SetDefault("FEATWS_RULLER_RESOLVER_BRIDGE_HEADERS", "")
	viper.SetDefault("FEATWS_RULLER_DEFAULT_RULES", "")
	viper.SetDefault("PORT", "8000")

	err = viper.ReadInConfig()
	if err != nil {
		if err2, ok := err.(*os.PathError); !ok {
			err = err2
			return
		}
	}

	err = viper.Unmarshal(config)

	config.ResourceLoaderHeaders = make(http.Header)
	resourceLoaderHeaders := strings.Split(config.ResourceLoaderHeadersStr, ",")
	for _, value := range resourceLoaderHeaders {
		entries := strings.Split(value, ":")
		if len(entries) == 2 {
			config.ResourceLoaderHeaders.Set(entries[0], entries[1])
		}
	}

	config.ResolverBridgeHeaders = make(http.Header)
	resolverBridgeHeaders := strings.Split(config.ResolverBridgeHeadersStr, ",")
	for _, value := range resolverBridgeHeaders {
		entries := strings.Split(value, ":")
		if len(entries) == 2 {
			config.ResolverBridgeHeaders.Set(entries[0], entries[1])
		}
	}

	return
}

//GetConfig ...
func GetConfig() *Config {
	return config
}

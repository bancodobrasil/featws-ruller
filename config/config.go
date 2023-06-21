package config

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

<<<<<<< HEAD
// Config type contains various configuration options for a program, including resource loader
// settings, port, default rules, SSL verification, resolver bridge settings, external host, and API
// key.
//
// Property:
//   - ResourceLoaderType: This property specifies the type of resource loader to be used by the application.
//   - ResourceLoaderURL: The URL of the resource loader used by the feature flag ruller to load feature flag rules.
//   - ResourceLoaderHeaders: This is a field of type http.Header which represents the headers to be sent with the HTTP requests made by the resource loader. It can be used to set custom headers such as authentication tokens or user agents.
//   - ResourceLoaderHeadersStr: This is a string representation of the HTTP headers that will be sent with requests made by the resource loader. These headers can be used to provide additional information or authentication credentials to the server being accessed. The headers will be parsed into an http.Header object before being used.
//   - Port: The port number on which the application will listen for incoming requests.
//   - DefaultRules: This property is used to specify the default rules that should be loaded by the resource loader. It is specified in the configuration file using the key "FEATWS_RULLER_DEFAULT_RULES".
//   - DisableSSLVerify: A boolean flag that indicates whether SSL verification should be disabled or not. If set to true, SSL verification will be disabled.
//   - ResolverBridgeURL: This property is a string that represents the URL of the resolver bridge. The resolver bridge is a service that is responsible for resolving feature flags and rules.
//   - ResolverBridgeHeaders: This field will be used to store HTTP headers that will be sent along with requests to the resolver bridge URL. The `http.Header` type is a map of strings to slices of strings, representing the headers and their values.
//   - ResolverBridgeHeadersStr: This property is a string representation of the HTTP headers that will be sent to the resolver bridge. It is used in conjunction with the ResolverBridgeHeaders property to set the headers for requests made to the resolver bridge. The headers can be specified as a JSON object in string format.
//   - ExternalHost: This property represents the external host name or IP address of the server where the application is running. It is used to construct URLs for external resources and APIs.
//   - AuthAPIKey: This property is used to store the API key for authentication purposes. It is used to authenticate requests made to the Ruller API.
=======
// Config ...
>>>>>>> cache-ruller
type Config struct {
	ResourceLoaderType       string `mapstructure:"FEATWS_RULLER_RESOURCE_LOADER_TYPE"`
	ResourceLoaderURL        string `mapstructure:"FEATWS_RULLER_RESOURCE_LOADER_URL"`
	ResourceLoaderHeaders    http.Header
	ResourceLoaderHeadersStr string `mapstructure:"FEATWS_RULLER_RESOURCE_LOADER_HEADERS"`

	Port             string `mapstructure:"PORT"`
	DefaultRules     string `mapstructure:"FEATWS_RULLER_DEFAULT_RULES"`
	DisableSSLVerify bool   `mapstructure:"FEATWS_DISABLE_SSL_VERIFY"`

	ResolverBridgeURL        string `mapstructure:"FEATWS_RULLER_RESOLVER_BRIDGE_URL"`
	ResolverBridgeHeaders    http.Header
	ResolverBridgeHeadersStr string `mapstructure:"FEATWS_RULLER_RESOLVER_BRIDGE_HEADERS"`

	ExternalHost string `mapstructure:"EXTERNAL_HOST"`

	AuthAPIKey string `mapstructure:"FEATWS_RULLER_API_KEY"`
}

var config = &Config{}

var loaded = false

// LoadConfig loads configuration settings from a file and sets default values for missing settings.
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
	viper.SetDefault("FEATWS_DISABLE_SSL_VERIFY", false)
	viper.SetDefault("EXTERNAL_HOST", "localhost:8000")
	viper.SetDefault("FEATWS_RULLER_API_KEY", "")

	err = viper.ReadInConfig()
	if err != nil {
		if err2, ok := err.(*os.PathError); !ok {
			err = err2
			log.Errorf("Error on Load Config: %v", err)
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

// GetConfig returns the loaded configuration or panics if there was an error loading it.
func GetConfig() *Config {
	if !loaded {
		err := LoadConfig()
		loaded = true
		if err != nil {
			panic(fmt.Sprintf("load config error: %s", err))
		}
	}
	return config
}

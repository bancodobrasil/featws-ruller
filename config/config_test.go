package config

import (
	"bytes"
	"io"
	"testing"

	"github.com/spf13/viper"
)

var configFileExample = []byte(`
	FEATWS_RULLER_RESOURCE_LOADER_TYPE="http"
	FEATWS_RULLER_RESOURCE_LOADER_URL="test/resourceLoaderUrl"
	FEATWS_RULLER_RESOURCE_LOADER_HEADERS="testHeader"
	FEATWS_RULLER_RESOLVER_BRIDGE_URL="test/bridgeUrl"
	FEATWS_RULLER_RESOLVER_BRIDGE_HEADERS="bridgeHeaders"
	FEATWS_RULLER_DEFAULT_RULES="testDefaultRules"
	PORT="4000"
	`)

var mockConfig = &configFileExample

func initConfig() {
	viper.Reset()
	var r io.Reader
	viper.SetConfigType("env")
	r = bytes.NewReader(configFileExample)
	viper.Unmarshal(r, mockConfig)
}

func TestLoadConfig(t *testing.T) {

}

func TestConfig(t *testing.T) {

}

package config

import (
	"github.com/spf13/viper"
)

type resourceLoader struct {
	Type    string            `mapstructure:"FEATWS_RULLER_RESOURCE_LOADER_TYPE"`
	URL     string            `mapstructure:"FEATWS_RULLER_RESOURCE_LOADER_URL"`
	Headers map[string]string `mapstructure:"FEATWS_RULLER_RESOURCE_LOADER_HEADERS"`
}

// O config ira armazenar todas as configurações da aplicação
// Os valores serão lidos pelo pelo viper de um arquivo de configuração ou de variaveis de ambiente
type Config struct {
	ResourceLoader resourceLoader
	Port           string `mapstructure:"FEATWS_RULLER_PORT"`
	DefaultRules   string `mapstructure:"FEATWS_RULLER_DEFAULT_RULES"`
}

// LoadConfig irá ler as configurações do arquivo ou varivaveis de ambiente
//path string
func LoadConfig() (config Config, err error) {
	//viper.AddConfigPath(path)

	viper.AutomaticEnv()
	//viper.SetDefault("Type", "http")
	viper.SetDefault("FEATWS_RULLER_RESOURCE_LOADER_HEADERS", map[string]string{})
	viper.SetDefault("FEATWS_RULLER_DEFAULT_RULES", "")
	viper.SetDefault("FEATWS_RULLER_PORT", "8000")

	/*err = viper.ReadInConfig()
	if err != nil {
		return
	}*/

	cfg := Config{
		ResourceLoader: resourceLoader{},
	}
	err = viper.Unmarshal(&cfg)
	return cfg, err
}

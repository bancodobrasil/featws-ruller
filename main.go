package main

import (
	"log"

	"github.com/bancodobrasil/featws-ruller/config"
	"github.com/gin-gonic/gin"
)

//TODO SUBSTITUIR PELO VIPER OU PORTA PADRAO
/*func getEnv(key, fallback string) string {
	value := viper.GetEnv(key)
	if value == nil {
		return fallback
	}

	return value.Error()

}*/

// DefaultKnowledgeBaseName its default name of Knowledge Base
const DefaultKnowledgeBaseName = "default"

// DefaultKnowledgeBaseVersion its default version of Knowledge Base
const DefaultKnowledgeBaseVersion = "latest"

//TODO RETIRATR ISSO AQUI
var Config config.Config

// Hello returns a greeting for the named person.
func main() {

	Config, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Não foi possível carregar as configurações")
	}

	if Config.DefaultRules != "" {
		defaultGRL := Config.DefaultRules
		log.Printf("Carregando '%s' como folha de regras default!", defaultGRL)
		err := loadLocalGRL(defaultGRL, DefaultKnowledgeBaseName, DefaultKnowledgeBaseVersion)
		if err != nil {
			panic(err)
		}
	} else {
		log.Println("Não foram carregadas regras default!")
	}

	router := gin.New()

	setupServer(router)

	port := Config.Port

	router.Run(":" + port)

}

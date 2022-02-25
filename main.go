package main

import (
	"log"

	"github.com/bancodobrasil/featws-ruller/config"
	"github.com/gin-gonic/gin"
)

// DefaultKnowledgeBaseName its default name of Knowledge Base
const DefaultKnowledgeBaseName = "default"

// DefaultKnowledgeBaseVersion its default version of Knowledge Base
const DefaultKnowledgeBaseVersion = "latest"

// my var Config is ...
var Config = config.Config{}

func main() {

	err := config.LoadConfig(&Config)
	if err != nil {
		log.Fatalf("Não foi possível carregar as configurações: %s\n", err)
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

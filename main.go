package main

import (
	"log"

	"github.com/bancodobrasil/featws-ruller/config"
	"github.com/bancodobrasil/featws-ruller/controller"
	"github.com/bancodobrasil/featws-ruller/route"
	"github.com/bancodobrasil/featws-ruller/service"
	"github.com/gin-gonic/gin"
)

func main() {

	err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Não foi possível carregar as configurações: %s\n", err)
	}

	cfg := config.GetConfig()

	if cfg.DefaultRules != "" {
		defaultGRL := cfg.DefaultRules
		log.Printf("Carregando '%s' como folha de regras default!", defaultGRL)
		err := service.LoadLocalGRL(defaultGRL, controller.DefaultKnowledgeBaseName, controller.DefaultKnowledgeBaseVersion)
		if err != nil {
			panic(err)
		}
	} else {
		log.Println("Não foram carregadas regras default!")
	}

	router := gin.New()

	route.SetupServer(router)

	port := cfg.Port

	router.Run(":" + port)

}

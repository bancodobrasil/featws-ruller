package main

import (
	"log"

	"github.com/bancodobrasil/featws-ruller/config"
	v1 "github.com/bancodobrasil/featws-ruller/controllers/v1"
	"github.com/bancodobrasil/featws-ruller/routes"
	"github.com/bancodobrasil/featws-ruller/services"
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
		err := services.LoadLocalGRL(defaultGRL, v1.DefaultKnowledgeBaseName, v1.DefaultKnowledgeBaseVersion)
		if err != nil {
			panic(err)
		}
	} else {
		log.Println("Não foram carregadas regras default!")
	}

	router := gin.New()

	routes.SetupRoutes(router)

	port := cfg.Port

	router.Run(":" + port)

}

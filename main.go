package main

import (
	"os"

	"github.com/bancodobrasil/featws-ruller/config"
	_ "github.com/bancodobrasil/featws-ruller/docs"
	"github.com/bancodobrasil/featws-ruller/routes"
	"github.com/bancodobrasil/featws-ruller/services"
	ginMonitor "github.com/bancodobrasil/gin-monitor"
	"github.com/bancodobrasil/goauth"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	ginlogrus "github.com/toorop/gin-logrus"
)

func setupLog() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	// log.SetFormatter(&log.JSONFormatter{})

	log.SetOutput(os.Stdout)

	log.SetLevel(log.DebugLevel)
}

// @title FeatWS Ruler

// @version 1.0

// @Description O projeto Ruler é uma implementação do motor de regras [grule-rule-engine](https://github.com/hyperjumptech/grule-rule-engine), que é utilizado para avaliar regras no formato .grl . O Ruler permite que as regras definidas em arquivos .grl sejam avaliadas de maneira automática e eficiente, ajudando a automatizar as decisões tomadas pelo FeatWS. Isso possibilita que o sistema possa analisar e classificar grandes quantidades de informações de maneira rápida e precisa.
// @Description
// @Description Ao utilizar as regras fornecidas pelo projeto Ruler, o FeatWS é capaz de realizar análises de regras em larga escala e fornecer resultados precisos e relevantes para seus usuários. Isso é especialmente importante em áreas como análise de sentimentos em mídias sociais, detecção de fraudes financeiras e análise de dados em geral.
// @Description
// @Description Antes de realizar os testes no Swagger, é necessário autorizar o acesso clicando no botão **Authorize**, ao lado, e inserindo a senha correspondente. Após inserir o campo **value** e clicar no botão **Authorize**, o Swagger estará disponível para ser utilizado.
// @Description
// @Description A seguir é explicado com mais detalhes sobre os endpoints:
// @Description  	- **/Eval**: Esse endpoint é utilizado apenas para aplicações que possuem uma única folha de regra padrão.
// @Description  	- **/Eval/{knowledgeBase}**: Nesse endpoint, é necessário informar o parâmetro com o nome da folha de regra desejada e, como resultado, será retornado a última versão da folha de regra correspondente.
// @Description  	- **/Eval/{knowledgeBase}/{version}**: Nesse endpoint é necessário colocar o parâmetro do nome da folha de regra como também o número da versão da folha de regra que você deseja testar a regra.
// @Description
// @Description		**Parameters / Parâmetros**
// @Description		No **knowledgeBase**, você pode especificar o nome da folha de regras que deseja utilizar. Já o **version** você coloca a versão que você deseja avaliar. Em **Paramenter**, é possível enviar os parametros que você deseja testar na folha de regra.

// @termsOfService http://swagger.io/terms/

// @contact.name API Support

// @contact.url http://www.swagger.io/support

// @contact.email support@swagger.io

// @license.name Apache 2.0

// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8000

// @BasePath /api/v1

// @securityDefinitions.apikey Authentication Api Key
// @in header
// @name X-API-Key

// @x-extension-openapi {"example": "value on a json format"}

// This function sets up a server using the Gin framework and loads default rules if specified in the
// configuration.
func main() {

	setupLog()

	err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Unable to load configurations: %s\n", err)
	}

	cfg := config.GetConfig()

	if cfg.DefaultRules != "" {
		defaultGRL := cfg.DefaultRules
		log.Debugf("Loading '%s' as defult rule sheet!", defaultGRL)
		err := services.EvalService.LoadLocalGRL(defaultGRL, services.DefaultKnowledgeBaseName, services.DefaultKnowledgeBaseVersion)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Warnln("Default rule sheet not loaded!")
	}

	monitor, err := ginMonitor.New("v0.3.2-rc1", ginMonitor.DefaultErrorMessageKey, ginMonitor.DefaultBuckets)
	if err != nil {
		log.Panic(err)
	}

	gin.DefaultWriter = log.StandardLogger().WriterLevel(log.DebugLevel)
	gin.DefaultErrorWriter = log.StandardLogger().WriterLevel(log.ErrorLevel)

	router := gin.New()

	goauth.BootstrapMiddleware()

	router.Use(ginlogrus.Logger(log.StandardLogger()), gin.Recovery())
	router.Use(monitor.Prometheus())
	router.GET("metrics", gin.WrapH(promhttp.Handler()))
	routes.SetupRoutes(router)
	routes.APIRoutes(router)

	port := cfg.Port

	router.Run(":" + port)

}

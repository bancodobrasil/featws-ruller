package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
)

type resourceLoader struct {
	Type    string            `yaml:"type"`
	URL     string            `yaml:"url"`
	Headers map[string]string `yaml:"headers"`
}

// Config its used to load config for resource loader
type Config struct {
	ResourceLoader resourceLoader `yaml:"resource-loader"`
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// DefaultKnowledgeBaseName its default name of Knowledge Base
const DefaultKnowledgeBaseName = "default"

// DefaultKnowledgeBaseVersion its default version of Knowledge Base
const DefaultKnowledgeBaseVersion = "latest"

var config = Config{}

// Hello returns a greeting for the named person.
func main() {

	arg := os.Args[1:]

	bytes, err := ioutil.ReadFile(arg[0])
	if err != nil {
		log.Fatalf("error: %v", err)
		panic(err)
	}

	err = yaml.Unmarshal(bytes, &config)
	if err != nil {
		log.Fatalf("error: %v", err)
		panic(err)
	}

	if len(arg) > 1 {
		defaultGRL := arg[1]
		log.Printf("Carregando '%s' como folha de regras default!", defaultGRL)

		err := loadLocalGRL(defaultGRL, DefaultKnowledgeBaseName, DefaultKnowledgeBaseVersion)
		if err != nil {
			panic(err)
		}
	} else {
		log.Println("Não foram carregadas regras default!")
	}
	router := gin.Default()

	srv := setupServer(router)

	log.Fatal(srv.ListenAndServe())
}

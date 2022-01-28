package main

import (
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type ResourceLoader struct {
	Type    string            `yaml:"type"`
	Url     string            `yaml:"url"`
	Headers map[string]string `yaml:"headers"`
}

type Config struct {
	ResourceLoader ResourceLoader `yaml:"resource-loader"`
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

const DEFAULT_KNOWLEDGEBASE_NAME = "default"
const DEFAULT_KNOWLEDGEBASE_VERSION = "latest"

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

		err := loadLocalGRL(defaultGRL, DEFAULT_KNOWLEDGEBASE_NAME, DEFAULT_KNOWLEDGEBASE_VERSION)
		if err != nil {
			panic(err)
		}
	} else {
		log.Println("Não foram carregadas regras default!")
	}

	srv := setupServer()

	log.Fatal(srv.ListenAndServe())
}

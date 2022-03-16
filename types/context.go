package types

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/bancodobrasil/featws-ruller/config"
	log "github.com/sirupsen/logrus"
)

// RemoteLoaded ...
type RemoteLoaded struct {
	Resolver string
}

// RemoteLoadeds ...
type RemoteLoadeds map[string]RemoteLoaded

// Context its used to store parameters and temporary variables during rule assertions
type Context struct {
	TypedMap
	RemoteLoadeds RemoteLoadeds
}

// NewContextFromMap method create a new Context
func NewContext() *Context {
	return NewContextFromMap(make(map[string]interface{}))
}

// NewContextFromMap method create a new Context from map
func NewContextFromMap(values map[string]interface{}) *Context {
	instance := &Context{
		TypedMap:      *NewTypedMapFromMap(values),
		RemoteLoadeds: make(map[string]RemoteLoaded),
	}
	instance.Getter = interface{}(instance).(Getter)
	return instance
}

// RegistryRemoteLoaded ...
func (c *Context) RegistryRemoteLoaded(param string, resolver string) {
	c.RemoteLoadeds[param] = RemoteLoaded{
		Resolver: resolver,
	}
}

func (c *Context) load(param string) interface{} {
	value := c.resolve(c.RemoteLoadeds[param].Resolver, param)
	c.Put(param, value)
	return value
}

func (c *Context) isRemoteLoaded(param string) bool {
	_, ok := c.RemoteLoadeds[param]
	return ok
}

// GetEntry ...
func (c *Context) GetEntry(param string) interface{} {
	value := c.TypedMap.GetEntry(param)

	if value == nil && c.isRemoteLoaded(param) {
		return c.load(param)
	}

	return value
}

type resolveInputV1 struct {
	Resolver string                 `json:"resolver"`
	Context  map[string]interface{} `json:"context"`
	Load     []string               `json:"load"`
}

type resolveOutputV1 struct {
	Context map[string]interface{} `json:"context"`
	Errors  map[string]interface{} `json:"errors"`
}

func (c *Context) resolve(resolver string, param string) interface{} {
	config := config.GetConfig()

	url := fmt.Sprintf("%s/api/v1/resolve", config.ResolverBridgeURL)

	url = strings.ReplaceAll(url, "//api/v1", "/api/v1")

	input := resolveInputV1{
		Resolver: resolver,
		Context:  c.GetEntries(),
		Load:     []string{param},
	}

	log.Debugf("Resolving with '%s': %v", url, input)

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(input)
	if err != nil {
		panic(err.Error())
	}

	req, err := http.NewRequest("POST", url, &buf)
	if err != nil {
		panic(err.Error())
	}

	req.Header = config.ResolverBridgeHeaders

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err.Error())
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}

	log.Debugf("Resolving with '%s': %v > %s", url, input, string(data))

	output := resolveOutputV1{}
	err = json.Unmarshal(data, &output)
	if err != nil {
		panic(err.Error())
	}

	if len(output.Errors) > 0 {
		panic(fmt.Sprintf("%s", output.Errors))
	}

	return output.Context[param]
}

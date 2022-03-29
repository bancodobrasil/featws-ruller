package types

import (
	"bytes"
	"crypto/tls"
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
	Resolver
	Loader
}

// Resolver ...
type Resolver interface {
	resolve(resolver string, param string) interface{}
}

// Loader ...
type Loader interface {
	load(param string) interface{}
}

// HTTPClient ...
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client ...
var Client HTTPClient = newHTTPClient()

func newHTTPClient() *http.Client {
	config := config.GetConfig()
	client := &http.Client{}
	if config.DisableSSLVerify {
		transCfg := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // ignore expired SSL certificates
		}
		client.Transport = transCfg
	}
	return client
}

// NewContext method create a new Context
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
	if c.Loader != nil {
		return c.Loader.load(param)
	}
	return c.loadImpl(param)
}

func (c *Context) loadImpl(param string) interface{} {
	remote, ok := c.RemoteLoadeds[param]
	if !ok {
		panic("The param it's not registry as remote loaded")
	}

	value := c.resolve(remote.Resolver, param)
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
	if c.Resolver != nil {
		return c.Resolver.resolve(resolver, param)
	}
	return c.resolveImpl(resolver, param)
}

func (c *Context) resolveImpl(resolver string, param string) interface{} {
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
		panic("error on encode input")
	}

	req, err := http.NewRequest("POST", url, &buf)
	if err != nil {
		panic("error on create Request")
	}

	req.Header = config.ResolverBridgeHeaders

	resp, err := Client.Do(req)
	if err != nil {
		panic("error on execute request")
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		panic("error on read the body")
	}

	log.Infof("Resolving with '%s': %v > %s", url, input, string(data))

	output := resolveOutputV1{}
	err = json.Unmarshal(data, &output)
	if err != nil {
		panic("error on response decoding")
	}

	if len(output.Errors) > 0 {
		panic(fmt.Sprintf("%s", output.Errors))
	}

	return output.Context[param]
}

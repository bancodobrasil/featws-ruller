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
	From     string
}

// RemoteLoadeds ...
type RemoteLoadeds map[string]RemoteLoaded

// Context its used to store parameters and temporary variables during rule assertions
type Context struct {
	TypedMap
	RemoteLoadeds  RemoteLoadeds
	RequiredParams []string
	Resolver
	Loader
	RequiredConfigured bool
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
		TypedMap:       *NewTypedMapFromMap(values),
		RemoteLoadeds:  make(map[string]RemoteLoaded),
		RequiredParams: []string{},
	}
	instance.Getter = interface{}(instance).(Getter)
	return instance
}

// RegistryRequiredParams ...
func (c *Context) RegistryRequiredParams(params ...string) {

	for _, param := range params {
		if c.isRemoteLoaded(param) {
			continue
		}
		c.RequiredParams = append(c.RequiredParams, param)
		if !c.Has(param) {
			c.addError("requiredParamErrors", param, fmt.Errorf("parameter %s is required", param))
		}

	}
}

func (c *Context) addError(key, param string, err interface{}) {
	if !c.Has(key) {
		c.Put(key, NewTypedMap())
	}

	switch r := err.(type) {
	case *log.Entry:
		c.GetMap(key).AddItem(param, r.Message)
	case log.Entry:
		c.GetMap(key).AddItem(param, r.Message)
	case string:
		c.GetMap(key).AddItem(param, r)
	case error:
		c.GetMap(key).AddItem(param, r.Error())
	default:
		c.GetMap(key).AddItem(param, fmt.Sprintf("%v", r))
	}
}

// RegistryRemoteLoadedWithFrom ...
func (c *Context) RegistryRemoteLoadedWithFrom(param string, resolver string, from string) {
	if from == "" {
		from = param
	}
	c.RemoteLoadeds[param] = RemoteLoaded{
		Resolver: resolver,
		From:     from,
	}
}

// RegistryRemoteLoaded ...
func (c *Context) RegistryRemoteLoaded(param string, resolver string) {
	c.RegistryRemoteLoadedWithFrom(param, resolver, "")
}

func (c *Context) load(param string) interface{} {
	if c.Loader != nil {
		return c.Loader.load(param)
	}
	return c.loadImpl(param)
}

func (c *Context) loadImpl(param string) interface{} {
	defer func() {
		r := recover()
		if r == nil {
			return
		}
		c.addError("errors", param, r)
	}()
	remote, ok := c.RemoteLoadeds[param]
	if !ok {
		log.Panic("The param it's not registry as remote loaded")
	}

	value := c.resolve(remote.Resolver, remote.From)
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
	Resolver string                 `json:"resolver,omitempty"`
	Context  map[string]interface{} `json:"context,omitempty"`
	Load     []string               `json:"load,omitempty"`
}

type resolveOutputV1 struct {
	Context map[string]interface{} `json:"context,omitempty"`
	Errors  map[string]interface{} `json:"errors,omitempty"`
	Error   string                 `json:"error,omitempty"`
}

func (c *Context) resolve(resolver string, param string) interface{} {
	if c.Resolver != nil {
		return c.Resolver.resolve(resolver, param)
	}
	return c.resolveImpl(resolver, param)
}

func (c *Context) resolveImpl(resolver string, param string) interface{} {
	config := config.GetConfig()

	url := fmt.Sprintf("%s/api/v1/resolve/%s", config.ResolverBridgeURL, resolver)

	url = strings.ReplaceAll(url, "//api/v1", "/api/v1")

	input := resolveInputV1{
		// Resolver: resolver,
		Context: c.GetEntries(),
		Load:    []string{param},
	}

	log.Debugf("Resolving with '%s': %v", url, input)

	var buf bytes.Buffer

	err := json.NewEncoder(&buf).Encode(input)
	if err != nil {
		log.Panic("error on encode input")
	}

	log.Debugf("Resolving with '%s' decoded: %v", url, buf.String())

	req, err := http.NewRequest("POST", url, &buf)
	if err != nil {
		log.Panic("error on create Request")
	}

	req.Header = config.ResolverBridgeHeaders

	resp, err := Client.Do(req)
	if err != nil {
		log.Panic("error on execute request")
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Panic("error on read the body")
	}

	log.Debugf("Resolving with '%s': %v > %s", url, input, string(data))

	output := resolveOutputV1{}
	err = json.Unmarshal(data, &output)
	if err != nil {
		log.Panic("error on response decoding")
	}

	if len(output.Errors) > 0 {
		log.Panic(fmt.Sprintf("%s", output.Errors))
	}

	if output.Error != "" {
		panic(output.Error)
	}

	return output.Context[param]
}

// SetRequiredConfigured ...
func (c *Context) SetRequiredConfigured() {
	c.RequiredConfigured = true
}

// IsReady ...
func (c *Context) IsReady() bool {
	return c.RequiredConfigured && !c.Has("requiredParamErrors")
}

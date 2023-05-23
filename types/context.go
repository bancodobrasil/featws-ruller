package types

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/bancodobrasil/featws-ruller/config"
	telemetry "github.com/bancodobrasil/gin-telemetry"
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
	RawContext context.Context
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

// The `isRemoteLoaded` method is a function defined on the `Context` struct in the Go programming
// language. It takes a single parameter `param` of type `string` and returns a boolean value. The
// purpose of this method is to check whether a given parameter is registered as a remote loaded
// parameter in the `RemoteLoadeds` field of the `Context` struct.
func (c *Context) isRemoteLoaded(param string) bool {
	_, ok := c.RemoteLoadeds[param]
	return ok
}

// GetEntry is used to retrieve a value from the `Context` struct. It first tries to get the value
// from the `TypedMap` field of the `Context` struct using the `GetEntry` method. If the value is not
// found in the `TypedMap`, it checks if the parameter is registered as a remote loaded parameter using
// the `isRemoteLoaded` method. If it is, it calls the `load` method to load the value from a remote
// resolver and stores it in the `TypedMap` before returning it. If the parameter is not found in the
// `TypedMap` and is not registered as a remote loaded parameter, it returns `nil`.
func (c *Context) GetEntry(param string) interface{} {
	value := c.TypedMap.GetEntry(param)

	if value == nil && c.isRemoteLoaded(param) {
		return c.load(param)
	}

	return value
}

// The `resolveInputV1` type in Go represents a set of optional parameters for resolving data,
// including a resolver string, context map, and load array.
//
// Property:
//   - Resolver {string}: The "Resolver" property is a string that specifies the name of the resolver function to be executed. A resolver function is a function that retrieves data from a data source and returns it to the client.
//   - Context: The `Context` property is a map of key-value pairs that can be used to provide additional information or context to the resolver function. The keys in the map are strings, and the values can be of any type. This property is optional and can be omitted if no additional context is needed.
//   - Load {[]string}: The `Load` property is an optional array of strings that specifies the fields to be loaded by the resolver. This is useful when you want to optimize the performance of your application by only loading the necessary data. The values in the `Load` array correspond to the fields in the GraphQL query that need
type resolveInputV1 struct {
	Resolver string                 `json:"resolver,omitempty"`
	Context  map[string]interface{} `json:"context,omitempty"`
	Load     []string               `json:"load,omitempty"`
}

// The type `resolveOutputV1` has three optional fields: `Context`, `Errors`, and `Error`, which are
// used for JSON serialization.
//
// Property:
//   - Context: Context is a map of key-value pairs that represent the context of the resolved output. This can include any relevant information that was used to generate the output, such as input parameters or metadata. The context is optional and may be omitted if not needed.
//   - Errors: The `Errors` property is a map that can contain any additional error information related to the operation being performed. It is optional and can be omitted if there are no errors to report. The keys in the map represent the error codes or error types, and the values can be any additional information related to
//     -Error {string}: The "Error" property is a string that represents a general error message. It is optional and may be used to provide additional information about any errors that occurred during the execution of a function or method.
type resolveOutputV1 struct {
	Context map[string]interface{} `json:"context,omitempty"`
	Errors  map[string]interface{} `json:"errors,omitempty"`
	Error   string                 `json:"error,omitempty"`
}

// This function is responsible for resolving a parameter using a remote resolver. It first checks if
// the `Resolver` field of the `Context` struct is not nil, and if it is not, it calls the `resolve`
// method of the `Resolver` interface to resolve the parameter. If the `Resolver` field is nil, it
// calls the `resolveImpl` method of the `Context` struct to resolve the parameter. The `resolveImpl`
// method constructs a request to the resolver bridge API, sends the request, and returns the resolved
// value.
func (c *Context) resolve(resolver string, param string) interface{} {
	if c.Resolver != nil {
		return c.Resolver.resolve(resolver, param)
	}
	return c.resolveImpl(resolver, param)
}

// The `resolveImpl` function is responsible for resolving a parameter using a remote resolver. It
// constructs a request to the resolver bridge API, sends the request, and returns the resolved value.
// The function takes two arguments: `resolver` and `param`, where `resolver` is the name of the
// resolver to use and `param` is the name of the parameter to resolve.
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

	ctx := c.RawContext
	var req *http.Request

	if ctx != nil {
		req, err = http.NewRequestWithContext(ctx, "POST", url, &buf)
	} else {
		req, err = http.NewRequest("POST", url, &buf)
	}

	if err != nil {
		log.Panic("error on create Request")
	}

	req.Header = config.ResolverBridgeHeaders

	if !telemetry.MiddlewareDisabled && ctx != nil {
		telemetry.Inject(ctx, req.Header)
	}

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

// SetRequiredConfigured sets the `RequiredConfigured` field of the `Context` struct to `true`. This field is
// used to keep track of whether all the required parameters have been configured or not. When this
// field is set to `true`, it means that all the required parameters have been configured and the
// context is ready to be used.
func (c *Context) SetRequiredConfigured() {
	c.RequiredConfigured = true
}

// IsReady returns a boolean value indicating whether the context is ready or not. A context is
// considered ready if all the required parameters have been configured and there are no errors
// related to missing required parameters.
func (c *Context) IsReady() bool {
	return c.RequiredConfigured && !c.Has("requiredParamErrors")
}

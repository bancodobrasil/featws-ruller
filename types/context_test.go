package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/bancodobrasil/featws-ruller/config"
	"github.com/sirupsen/logrus"
)

func TestMain(m *testing.M) {
	// setup()
	code := m.Run()
	// shutdown()
	os.Exit(code)
}

// TestNewContext start
func TestNewContext(t *testing.T) {
	got := NewContext()

	expected := &Context{
		TypedMap: *NewTypedMap(),
	}

	if reflect.TypeOf(got) != reflect.TypeOf(expected) {
		t.Error("You got an error while try to generate new context")
	}
}

// TestNewContext stop

// TestRegistryRemoteLoaded start
func TestRegistryRemoteLoaded(t *testing.T) {
	ct := NewContext()
	ct.RegistryRemoteLoaded("myparam", "myresolver")
	got := ct.isRemoteLoaded("myparam")
	expected := true
	if got != expected {
		t.Errorf("Test Fail, we want %v, we got %v", got, expected)
	}

}

// TestRegistryRemoteLoaded stop

// TestRegistryNotRemoteLoaded start
func TestRegistryNotRemoteLoaded(t *testing.T) {
	ct := NewContext()
	// ct.RegistryRemoteLoaded("myparam", "myresolver")
	got := ct.isRemoteLoaded("myparam")
	expected := false
	if got != expected {
		t.Errorf("Test Fail, we want %v, we got %v", got, expected)
	}

}

// TestRegistryNotRemoteLoaded stop

// TestLoad start
type MockContextLoad struct {
	Context
	t *testing.T
}

func (m *MockContextLoad) resolve(resolver string, param string) interface{} {
	if resolver != "myresolver" {
		m.t.Error("the resolvers are not the same")
	}

	if param != "myRemoteParam" {
		m.t.Error("the params are not the same")
	}
	return "myresult"

}
func TestLoad(t *testing.T) {
	ctx := &MockContextLoad{
		Context: *NewContext(),
		t:       t,
	}

	ctx.Resolver = ctx

	ctx.RegistryRemoteLoaded("myRemoteParam", "myresolver")

	want := ctx.load("myRemoteParam")
	expected := "myresult"

	if want != expected {
		t.Errorf("Couldn't load the resolve")
	}

}

// TestLoad stop

// TestLoadPanicNotRemoted start
func TestLoadPanicNotRemoted(t *testing.T) {
	ctx := NewContext()
	ctx.load("myRemoteParam")

	got := ctx.GetMap("errors").GetSlice("myRemoteParam")[0]
	expected := "The param it's not registry as remote loaded"

	if !reflect.DeepEqual(got, expected) {
		t.Error("The error message it's not throwed")
	}
}

// TestLoadPanicNotRemoted stop

// TestGetEntryRemoteLoaded start
type MockContextGetEntry struct {
	Context
	t *testing.T
}

func (m *MockContextGetEntry) load(param string) interface{} {
	if param != "myRemoteParam" {
		m.t.Error("the params are not the same")
	}
	return "myresult"

}
func TestGetEntryRemoteLoaded(t *testing.T) {
	ctx := &MockContextGetEntry{
		Context: *NewContext(),
		t:       t,
	}

	ctx.Loader = ctx

	ctx.RegistryRemoteLoaded("myRemoteParam", "myresolver")
	got := ctx.GetEntry("myRemoteParam")

	expected := "myresult"

	if got != expected {
		t.Errorf("Couldn't Get the entries")
	}

}

// TestGetEntryRemoteLoaded stop

// TestGetEntryContext start
func TestGetEntryContext(t *testing.T) {
	ctx := NewContext()
	ctx.Put("mystring", "teste")
	got := ctx.GetEntry("mystring")

	expected := "teste"

	if got != expected {
		t.Errorf("Couldn't Get the entries")
	}

}

// TestGetEntryContext stop

// TestEncodePanic start
type MockHTTPClientEncodePanic struct {
	http.Client
	t *testing.T
}

func (m *MockHTTPClientEncodePanic) Do(req *http.Request) (*http.Response, error) {
	defer req.Body.Close()

	data, err := io.ReadAll(req.Body)
	if err != nil {
		panic(err.Error())
	}

	input := resolveInputV1{}
	err = json.Unmarshal(data, &input)
	if err != nil {
		panic(err.Error())
	}

	if input.Resolver != "resolver_name" {
		m.t.Error("the resolver is wrong")
	}

	expectedContext := make(map[string]interface{})
	expectedContext["mystring"] = "teste"

	if !reflect.DeepEqual(input.Context, expectedContext) {
		m.t.Error("The contexts are not the same")
	}

	expectedLoad := []string{"param_name"}

	if !reflect.DeepEqual(input.Load, expectedLoad) {
		m.t.Error("The loads dont match each other")
	}

	stringReader := strings.NewReader("{}")

	return &http.Response{
		Body: io.NopCloser(stringReader),
	}, nil
}
func TestEncodePanic(t *testing.T) {
	defer func() {
		r := recover()
		if r.(*logrus.Entry).Message != "error on encode input" {
			t.Error("The panic message it's not throwed")
		}
	}()
	Client = &MockHTTPClientEncodePanic{
		t: t,
	}

	ctx := NewContext()
	ctx.Put("mystring", math.Inf(0))

	ctx.resolve("resolver_name", "param_name")

}

// TestEncodePanic stop

// TestRequestPanic start
func TestRequestPanic(t *testing.T) {
	defer func() {
		r := recover()
		config.LoadConfig()
		if r.(*logrus.Entry).Message != "error on create Request" {
			t.Error("The panic message it's not throwed")
		}
	}()

	config := config.GetConfig()
	config.ResolverBridgeURL = "\n"

	ctx := NewContext()

	ctx.resolve("resolver_name", "param_name")
}

//TestRequestPanic stop

// TestRequestExecutePanic start
type MockHTTPClientExecutePanic struct {
	http.Client
	t *testing.T
}

func (m *MockHTTPClientExecutePanic) Do(req *http.Request) (*http.Response, error) {
	return nil, fmt.Errorf("mock do error")
}

func TestRequestExecutePanic(t *testing.T) {
	defer func() {
		r := recover()
		if r.(*logrus.Entry).Message != "error on execute request" {
			t.Error("The panic message it's not throwed")
		}
	}()

	Client = &MockHTTPClientExecutePanic{
		t: t,
	}

	ctx := NewContext()

	ctx.resolve("resolver_name", "param_name")
}

// TestRequestExecutePanic stop

// TestResponseDecodePanic start
type MockHTTPClientResponseDecodePanic struct {
	http.Client
	t *testing.T
}

func (m *MockHTTPClientResponseDecodePanic) Do(req *http.Request) (*http.Response, error) {
	data, err := io.ReadAll(req.Body)
	if err != nil {
		panic(err.Error())
	}

	input := resolveInputV1{}
	err = json.Unmarshal(data, &input)
	if err != nil {
		panic(err.Error())
	}

	if input.Resolver != "resolver_name" {
		m.t.Error("the resolver is wrong")
	}

	expectedContext := make(map[string]interface{})
	expectedContext["mystring"] = "teste"

	if !reflect.DeepEqual(input.Context, expectedContext) {
		m.t.Error("The contexts are not the same")
	}

	expectedLoad := []string{"param_name"}

	if !reflect.DeepEqual(input.Load, expectedLoad) {
		m.t.Error("The loads dont match each other")
	}

	stringReader := strings.NewReader("")

	return &http.Response{
		Body: io.NopCloser(stringReader),
	}, nil
}
func TestResponseDecodePanic(t *testing.T) {
	defer func() {
		r := recover()
		if r.(*logrus.Entry).Message != "error on response decoding" {
			t.Error("The panic message it's not throwed")
		}
	}()
	Client = &MockHTTPClientResponseDecodePanic{
		t: t,
	}

	ctx := NewContext()
	ctx.Put("mystring", "teste")

	ctx.resolve("resolver_name", "param_name")
}

// TestResponseDecodePanic stop

// TestPanicOnReadBody start
type MockHTTPClientReadBodyPanic struct {
	http.Client
	t *testing.T
}

func (m *MockHTTPClientReadBodyPanic) Read(p []byte) (n int, err error) {
	return 0, errors.New(`everything is broken /o\`)
}

func (m *MockHTTPClientReadBodyPanic) Do(req *http.Request) (*http.Response, error) {

	data, err := io.ReadAll(req.Body)
	if err != nil {
		panic(err.Error())
	}

	input := resolveInputV1{}
	err = json.Unmarshal(data, &input)
	if err != nil {
		panic(err.Error())
	}

	if input.Resolver != "resolver_name" {
		m.t.Error("the resolver is wrong")
	}

	expectedContext := make(map[string]interface{})
	expectedContext["mystring"] = "teste"

	if !reflect.DeepEqual(input.Context, expectedContext) {
		m.t.Error("The contexts are not the same")
	}

	expectedLoad := []string{"param_name"}

	if !reflect.DeepEqual(input.Load, expectedLoad) {
		m.t.Error("The loads dont match each other")
	}

	return &http.Response{
		Body: io.NopCloser(&MockHTTPClientReadBodyPanic{}),
	}, nil
}
func TestPanicOnReadBody(t *testing.T) {
	defer func() {
		r := recover()
		if r.(*logrus.Entry).Message != "error on read the body" {
			t.Error("The panic message it's not throwed")
		}
	}()
	Client = &MockHTTPClientReadBodyPanic{
		t: t,
	}

	ctx := NewContext()
	ctx.Put("mystring", "teste")

	ctx.resolve("resolver_name", "param_name")
}

// TestPanicOnReadBody stop

// TestResponseDecodeMoreThenZero start
type MockHTTPClientResponseDecodeMoreThenOneError struct {
	http.Client
	t *testing.T
}

func (m *MockHTTPClientResponseDecodeMoreThenOneError) Do(req *http.Request) (*http.Response, error) {
	data, err := io.ReadAll(req.Body)
	if err != nil {
		panic(err.Error())
	}

	input := resolveInputV1{}
	err = json.Unmarshal(data, &input)
	if err != nil {
		panic(err.Error())
	}

	if input.Resolver != "resolver_name" {
		m.t.Error("the resolver is wrong")
	}

	expectedContext := make(map[string]interface{})
	expectedContext["mystring"] = "teste"

	if !reflect.DeepEqual(input.Context, expectedContext) {
		m.t.Error("The contexts are not the same")
	}

	expectedLoad := []string{"param_name"}

	if !reflect.DeepEqual(input.Load, expectedLoad) {
		m.t.Error("The loads dont match each other")
	}

	stringReader := strings.NewReader(`{
		"errors": {
			"myparam": "myerror"
		}
	}`)

	return &http.Response{
		Body: io.NopCloser(stringReader),
	}, nil
}

func TestResponseDecodeMoreThenZero(t *testing.T) {
	defer func() {
		r := recover()
		if r.(*logrus.Entry).Message != "map[myparam:myerror]" {
			t.Error("The panic message it's not throwed")
		}
	}()
	Client = &MockHTTPClientResponseDecodeMoreThenOneError{
		t: t,
	}

	ctx := NewContext()
	ctx.Put("mystring", "teste")

	ctx.resolve("resolver_name", "param_name")

}

// TestResponseDecodeMoreThenZero stop

// TestResolve start
type MockHTTPClient struct {
	http.Client
	t *testing.T
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	data, err := io.ReadAll(req.Body)
	if err != nil {
		panic(err.Error())
	}

	input := resolveInputV1{}
	err = json.Unmarshal(data, &input)
	if err != nil {
		panic(err.Error())
	}

	if input.Resolver != "resolver_name" {
		m.t.Error("the resolver is wrong")
	}

	expectedContext := make(map[string]interface{})
	expectedContext["mystring"] = "teste"

	if !reflect.DeepEqual(input.Context, expectedContext) {
		m.t.Error("The contexts are not the same")
	}

	expectedLoad := []string{"param_name"}

	if !reflect.DeepEqual(input.Load, expectedLoad) {
		m.t.Error("The loads dont match each other")
	}

	stringReader := strings.NewReader(`{
		"context": {
			"param_name": "myresult"
		}
	}`)

	return &http.Response{
		Body: io.NopCloser(stringReader),
	}, nil
}

func TestResolve(t *testing.T) {
	Client = &MockHTTPClient{
		t: t,
	}
	ctx := NewContext()
	ctx.Put("mystring", "teste")

	got := ctx.resolve("resolver_name", "param_name")
	expected := "myresult"

	if got != expected {
		t.Error("failed to resolve")
	}

}

// TestResolve stop

// TestResolveUnexpectedError start
type MockHTTPClientUnexpectedError struct {
	http.Client
	t *testing.T
}

func (m *MockHTTPClientUnexpectedError) Do(req *http.Request) (*http.Response, error) {
	data, err := io.ReadAll(req.Body)
	if err != nil {
		panic(err.Error())
	}

	input := resolveInputV1{}
	err = json.Unmarshal(data, &input)
	if err != nil {
		panic(err.Error())
	}

	if input.Resolver != "resolver_name" {
		m.t.Error("the resolver is wrong")
	}

	expectedLoad := []string{"param_name"}

	if !reflect.DeepEqual(input.Load, expectedLoad) {
		m.t.Error("The loads dont match each other")
	}

	stringReader := strings.NewReader(`{"error":"error message"}`)

	return &http.Response{
		Body: io.NopCloser(stringReader),
	}, nil
}

func TestResolveUnexpectedError(t *testing.T) {
	defer func() {
		r := recover()
		if r != "error message" {
			t.Error("The panic message it's not throwed")
		}
	}()
	Client = &MockHTTPClientUnexpectedError{
		t: t,
	}
	ctx := NewContext()

	ctx.resolve("resolver_name", "param_name")
}

// TestResolve stop

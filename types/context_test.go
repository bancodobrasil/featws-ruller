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

// This test sets up and shuts down the test environment.
func TestMain(m *testing.M) {
	// setup()
	code := m.Run()
	// shutdown()
	os.Exit(code)
}

// TestNewContext checks if a new context is generated correctly.
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

// TestRegistryRemoteLoadedWithFrom checks if a registry is loaded remotely with a specific parameter,
// resolver, and from value.
func TestRegistryRemoteLoadedWithFrom(t *testing.T) {
	ctx := NewContext()
	ctx.RegistryRemoteLoadedWithFrom("myparam", "myresolver", "myfrom")
	got := ctx.isRemoteLoaded("myparam")
	expected := true
	if got != expected {
		t.Errorf("Test Fail, we want %v, we got %v", got, expected)
	}

}

// TestRegistryRemoteLoaded checks if a registry remote is loaded.
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

// TestRegistryNotRemoteLoaded checks if a registry is not remotely loaded.
func TestRegistryNotRemoteLoaded(t *testing.T) {
	ct := NewContext()
	// ct.RegistryRemoteLoaded("myparam", "myresolver")
	got := ct.isRemoteLoaded("myparam")
	expected := false
	if got != expected {
		t.Errorf("Test Fail, we want %v, we got %v", got, expected)
	}

}

// MockContextLoad is a struct that embeds the `Context` type and includes a pointer to a `testing.T` object.
// The code snippet defines a struct named `MockContextLoad` with the properties:
//
// Property:
//   - t: is a pointer to the testing.T struct that provide methods for logging and reporting test results. It is commonly used in unit tests to assert expected behavior and report any failures. In this case, it is being used as a property of
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

// This is a unit test function that tests the load function of a context object.
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

// MockContextLoadWithFrom embeds the `Context` type and includes a testing `T` object.
//
// Property:
//   - Context: The `MockContextLoadWithFrom` struct is a custom struct that embeds the `Context` interface and has an additional property `t` of type `*testing.T`. This struct is likely used in unit tests to mock a context object and provide a testing context for functions that require a context.
//   - t: is a pointer to the testing.T struct, which is used for logging and reporting test results in Go unit tests. It is typically passed as an argument to test functions and methods.
type MockContextLoadWithFrom struct {
	Context
	t *testing.T
}

// resolve method takes two parameters `resolver` and `param` and returns an interface. It checks if the `resolver`
// parameter is equal to "myresolver" and the `param` parameter is equal to "myfrom". If they are not
// equal, it raises an error. If they are equal, it returns the string "myresult".
func (m *MockContextLoadWithFrom) resolve(resolver string, param string) interface{} {
	if resolver != "myresolver" {
		m.t.Error("the resolvers are not the same")
	}

	if param != "myfrom" {
		m.t.Error("the params are not the same")
	}
	return "myresult"

}

// TestLoadWithFrom tests the "load" function with a specific input and expected output.
func TestLoadWithFrom(t *testing.T) {
	ctx := &MockContextLoadWithFrom{
		Context: *NewContext(),
		t:       t,
	}

	ctx.Resolver = ctx

	ctx.RegistryRemoteLoadedWithFrom("myRemoteParam", "myresolver", "myfrom")

	want := ctx.load("myRemoteParam")
	expected := "myresult"

	if want != expected {
		t.Errorf("Couldn't load the resolve")
	}
}

// TestLoadPanicNotRemoted checks if an error message is thrown when a remote parameter isn't loaded.
func TestLoadPanicNotRemoted(t *testing.T) {
	ctx := NewContext()
	ctx.load("myRemoteParam")

	got := ctx.GetMap("errors").GetSlice("myRemoteParam")[0]
	expected := "The param it's not registry as remote loaded"

	if !reflect.DeepEqual(got, expected) {
		t.Error("The error message it's not throwed")
	}
}

// MockContextGetEntry is a struct that embeds the `Context` type and includes a testing `t` object.
//
// Property:
//   - Context: `MockContextGetEntry` is a struct type.
//   - t: is a pointer to the testing.T struct, which is used in Go's testing package for writing tests and reporting test results. It provides methods for logging test output, marking tests as failed or skipped, and measuring test execution time.
type MockContextGetEntry struct {
	Context
	t *testing.T
}

// load is a struct called `MockContextGetEntry`, this method takes a string parameter called `param` and returns an interface.
func (m *MockContextGetEntry) load(param string) interface{} {
	if param != "myRemoteParam" {
		m.t.Error("the params are not the same")
	}
	return "myresult"

}

// TestGetEntryRemoteLoaded checks if a remote entry is loaded correctly.
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

// TestGetEntryContext tests the GetEntryContext method of a context object in Go programming language.
func TestGetEntryContext(t *testing.T) {
	ctx := NewContext()
	ctx.Put("mystring", "teste")
	got := ctx.GetEntry("mystring")

	expected := "teste"

	if got != expected {
		t.Errorf("Couldn't Get the entries")
	}

}

// MockHTTPClientEncodePanic is a struct that embeds the http.Client type and includes a testing.T object for testing purposes.
// The `MockHTTPClientEncodePanic` struct is a custom implementation of the `http.Client` struct that includes an additional property `t` of type `*testing.T`. This struct is likely used in unit tests to simulate a panic during encoding.
//
// Property:
//   - t: is a pointer to a testing.T object, which is used for testing in Go. It provides methods for reporting test failures and logging test output.
type MockHTTPClientEncodePanic struct {
	http.Client
	t *testing.T
}

// Do is a mock HTTP client that reads the request body, unmarshals it into a struct, and checks if the values of the struct fields match the expected
// values. If the values do not match, it returns an error. If the values match, it returns an HTTP
// response with an empty JSON object. The method also includes panic statements to handle errors
// during the reading and unmarshaling of the request body.
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

// This function tests if a panic message is thrown when encoding input fails.
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

// TestRequestPanic tests if a panic message is thrown when a request cannot be created.
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

// MockHTTPClientExecutePanic is a struct that embeds the `http.Client` type and includes a testing `t` object for panic handling.
// The `MockHTTPClientExecutePanic` struct is a custom implementation of the `http.Client`
// struct that includes an additional property `t` of type `*testing.T`. This struct is likely used in
// unit tests to simulate a panic during an HTTP request and capture the panic for testing purposes.
//
// Property:
//   - t: is a pointer to a testing.T object, which is used for logging and reporting test results in Go unit tests. It allows the MockHTTPClientExecutePanic struct to log any errors or failures that occur during testing.
type MockHTTPClientExecutePanic struct {
	http.Client
	t *testing.T
}

// A method `Do` for a struct `MockHTTPClientExecutePanic` in Go. This method takes a pointer to an
// `http.Request` as input and returns a `nil` response and an error with a message "mock do error".
// This is likely used for testing purposes to simulate an error occurring during an HTTP request.
func (m *MockHTTPClientExecutePanic) Do(req *http.Request) (*http.Response, error) {
	return nil, fmt.Errorf("mock do error")
}

// TestRequestExecutePanic tests if a panic message is thrown when executing a request.
func TestRequestExecutePanic(t *testing.T) {
	defer func() {
		r := recover()
		if r.(*logrus.Entry).Message != "error on execute request: mock do error" {
			t.Error("The panic message it's not throwed")
		}
	}()

	Client = &MockHTTPClientExecutePanic{
		t: t,
	}

	ctx := NewContext()

	ctx.resolve("resolver_name", "param_name")
}

// TestResponseDecodePanic  is a struct that embeds the `http.Client` type and includes a testing `t` object for panic
// handling.
// @property  - - `MockHTTPClientResponseDecodePanic`: a struct type
// that represents a mock HTTP client response that will cause a panic
// when decoded.
//
// Property:
//   - t: is a pointer to a testing.T object, which is used for testing in Go. It provides methods for reporting test failures and logging test output.
type MockHTTPClientResponseDecodePanic struct {
	http.Client
	t *testing.T
}

// Do method is a mock HTTP client that is used for testing purposes. The method reads the request body, unmarshals it into a `resolveInputV1` struct, and
// checks if the `Context` and `Load` fields match the expected values. If they don't match, it raises
// an error. Finally, it returns an empty HTTP response.
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

	// if input.Resolver != "resolver_name" {
	// 	m.t.Error("the resolver is wrong")
	// }

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

// TestResponseDecodePanic tests if a panic message is thrown when decoding a response.
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

// The type MockHTTPClientReadBodyPanic is a struct that extends the http.Client type and includes a
// testing.T object for error reporting.
// The code defines a struct named `MockHTTPClientReadBodyPanic` which embeds the
// `http.Client` struct. This means that `MockHTTPClientReadBodyPanic` has access to all the fields and
// methods of `http.Client`.
//
// Property:
//   - t: is a pointer to a testing.T object, which is used for testing in Go. It provides methods for reporting test failures and logging test output.
type MockHTTPClientReadBodyPanic struct {
	http.Client
	t *testing.T
}

// Read is a mock HTTP client struct `MockHTTPClientReadBodyPanic`. This method returns an error with a message "everything is broken
// /o\". This is likely used for testing error handling in code that uses the HTTP client.
func (m *MockHTTPClientReadBodyPanic) Read(p []byte) (n int, err error) {
	return 0, errors.New(`everything is broken /o\`)
}

// Do method for a mock HTTP client that reads the body of an HTTP request and unmarshals it into a `resolveInputV1` struct. It then checks if the `Resolver` field of
// the struct is equal to a specific value and if the `Context` and `Load` fields of the struct match
// expected values. If any of these checks fail, it raises an error. Finally, it returns a mock HTTP
// response with an empty body. If there is an error reading the request body or unmarshaling it into
// the struct, the method pan
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

	// if input.Resolver != "resolver_name" {
	// 	m.t.Error("the resolver is wrong")
	// }

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
		if r.(*logrus.Entry).Message != "error on read the body: everything is broken /o\\" {
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

// The type MockHTTPClientResponseDecodeMoreThenOneError is a struct that extends the http.Client type and includes a
// testing.T object for error reporting.
// The code defines a struct named `MockHTTPClientResponseDecodeMoreThenOneError` which embeds the
// `http.Client` struct. This means that `MockHTTPClientResponseDecodeMoreThenOneError` has access to all the fields and
// methods of `http.Client`.
//
// Property:
//   - t: is a pointer to a testing.T object, which is used for testing in Go. It provides methods for reporting test failures and logging test output.
type MockHTTPClientResponseDecodeMoreThenOneError struct {
	http.Client
	t *testing.T
}

// Do method is a mock HTTP client that reads the request body, unmarshals it into a `resolveInputV1` struct, checks if the context and load fields of the struct
// match the expected values, and then returns a response with a JSON body containing an error message
// for a specific parameter. This is likely being used for testing purposes to simulate an error
// response from an external API.
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

	// if input.Resolver != "resolver_name" {
	// 	m.t.Error("the resolver is wrong")
	// }

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

// TestResponseDecodeMoreThenZero tests if a panic message is thrown with a specific error message.
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

// The typeMockHTTPClient is a struct that extends the http.Client type and includes a
// testing.T object for error reporting.
// The code defines a struct named `MockHTTPClient` which embeds the
// `http.Client` struct. This means that `MockHTTPClient` has access to all the fields and
// methods of `http.Client`.
//
// Property:
//   - t: is a pointer to a testing.T object, which is used for testing in Go. It provides methods for reporting test failures and logging test output.
type MockHTTPClient struct {
	http.Client
	t *testing.T
}

// Do is a mock HTTP client in Go. This method reads the request body, unmarshals it into a `resolveInputV1` struct, and checks if the `Context` and `Load`
// fields of the struct match the expected values. If they don't match, it raises an error. Finally, it
// returns a mock HTTP response with a predefined JSON string as the response body.
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

	// if input.Resolver != "resolver_name" {
	// 	m.t.Error("the resolver is wrong")
	// }

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

// TestResolve is a unit test function in Go that tests the resolve function of a context object.
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

// The type MockHTTPClientUnexpectedError is a struct that extends the http.Client type and includes a
// testing.T object for error reporting.
// The code defines a struct named `MockHTTPClientUnexpectedError` which embeds the
// `http.Client` struct. This means that `MockHTTPClientUnexpectedError` has access to all the fields and
// methods of `http.Client`.
//
// Property:
//   - t: is a pointer to a testing.T object, which is used for testing in Go. It provides methods for reporting test failures and logging test output.
type MockHTTPClientUnexpectedError struct {
	http.Client
	t *testing.T
}

// Do is a mock HTTP client that returns an unexpected error response. It reads the request body, unmarshals it into a `resolveInputV1` struct, and checks if the
// `Load` field matches an expected value. If it doesn't match, it raises an error. It then returns an
// HTTP response with an error message in the body.
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

	// if input.Resolver != "resolver_name" {
	// 	m.t.Error("the resolver is wrong")
	// }

	expectedLoad := []string{"param_name"}

	if !reflect.DeepEqual(input.Load, expectedLoad) {
		m.t.Error("The loads dont match each other")
	}

	stringReader := strings.NewReader(`{"error":"error message"}`)

	return &http.Response{
		Body: io.NopCloser(stringReader),
	}, nil
}

// TestResolveUnexpectedError is a test function in Go that checks if a panic message is thrown correctly.
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

package v1

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/bancodobrasil/featws-ruller/common/errors"

	"github.com/bancodobrasil/featws-ruller/types"
	"github.com/gin-gonic/gin"
	"github.com/hyperjumptech/grule-rule-engine/ast"
	"github.com/hyperjumptech/grule-rule-engine/builder"
	"github.com/hyperjumptech/grule-rule-engine/pkg"
)

// The function creates a mock Gin context and HTTP response recorder for testing purposes.
func mockGin() (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// test request, must instantiate a request first
	req := &http.Request{
		URL:    &url.URL{},
		Header: make(http.Header), // if you need to test headers
	}
	// finally set the request to the gin context
	c.Request = req

	return c, w
}

// This is a struct that implements the IEval interface and includes a testing.T object for testing
// purposes, but does not require a knowledge base or version.
//
// Property:
//   - `EvalServiceTestEvalHandlerWithoutKnowledgeBaseAndVersion` is a struct type that embeds the `services.IEval` interface and has a `t` property of type `*testing.T`.
//   - t: is a pointer to the testing.T struct, which is used for logging and reporting test results in Go unit tests. It is typically passed as an argument to test functions.
type EvalServiceTestEvalHandlerWithoutKnowledgeBaseAndVersion struct {
	services.IEval
	t *testing.T
}

// This function is implementing the `LoadRemoteGRL` method of the `services.IEval` interface for
// testing purposes. It checks if the `knowledgeBaseName` and `version` parameters passed to the method
// are equal to the default values defined in the `services` package. If they are not equal, it logs an
// error using the `testing.T` object and returns `nil`. If they are equal, it simply returns `nil`.
func (s EvalServiceTestEvalHandlerWithoutKnowledgeBaseAndVersion) LoadRemoteGRL(knowledgeBaseName string, version string) error {
	if knowledgeBaseName != services.DefaultKnowledgeBaseName || version != services.DefaultKnowledgeBaseVersion {
		s.t.Error("Did not load the default")
	}
	return nil
}

// This function is implementing the `GetKnowledgeLibrary` method of the `services.IEval` interface for
// testing purposes. It returns a new instance of the `ast.KnowledgeLibrary` struct, which is an
// in-memory representation of a knowledge base that can be used to store and manage rules. This
// function is used in the test case `TestEvalHandlerWithoutKnowledgeBaseAndVersion` to create a mock
// implementation of the `services.IEval` interface that does not require a knowledge base or version
// to be loaded.
func (s EvalServiceTestEvalHandlerWithoutKnowledgeBaseAndVersion) GetKnowledgeLibrary() *ast.KnowledgeLibrary {
	return ast.NewKnowledgeLibrary()
}

func (s EvalServiceTestEvalHandlerWithoutKnowledgeBaseAndVersion) GetKnowledgeBase(knowledgeBaseName string, version string) (*ast.KnowledgeBase, *errors.RequestError) {
	return nil, &errors.RequestError{Message: "KnowledgeBase or version not found", StatusCode: 404}
}

// This is a test function for the EvalHandler, which checks if the function returns a
// 404 error and a specific error message when called without a knowledge base or version.
func TestEvalHandlerWithoutKnowledgeBaseAndVersion(t *testing.T) {

	services.EvalService = EvalServiceTestEvalHandlerWithoutKnowledgeBaseAndVersion{
		t: t,
	}

	c, r := mockGin()
	EvalHandler()(c)
	gotStatus := r.Result().Status
	expectedStatus := "404 Not Found"

	if gotStatus != expectedStatus {
		t.Error("got error on request evalHandler func")
	}

	gotBody := r.Body.String()
	expectedBody := "KnowledgeBase or version not found"

	if gotBody != expectedBody {
		t.Error("got error on request evalHandler func")
	}
}

// The type EvalServiceTestEvalHandlerLoadError is a struct that implements the IEval interface and
// includes a testing.T object.
//
// Property:
//   - `EvalServiceTestEvalHandlerLoadError` is a struct type that embeds the `services.IEval` interface and has an additional property `t` of type `*testing.T`.
//   - t: is a pointer to the testing.T struct, which is used in Go's testing package to provide methods for writing tests and reporting test results. It is commonly used to log test failures and errors.
type EvalServiceTestEvalHandlerLoadError struct {
	services.IEval
	t *testing.T
}

// This function is implementing the `LoadRemoteGRL` method of the `services.IEval` interface for
// testing purposes. It returns a mock error message when called, which is used to simulate an error
// that might occur when loading a knowledge base or version. This function is used in the test case
// `TestEvalHandlerLoadError` to create a mock implementation of the `services.IEval` interface that
// returns an error when the `EvalHandler` function is called.
func (s EvalServiceTestEvalHandlerLoadError) LoadRemoteGRL(knowledgeBaseName string, version string) error {
	return fmt.Errorf("mock load error")
}

// This function is implementing the `GetKnowledgeLibrary` method of the `services.IEval` interface for
// testing purposes. It returns a new instance of the `ast.KnowledgeLibrary` struct, which is an
// in-memory representation of a knowledge base that can be used to store and manage rules. This
// function is used in the test case `TestEvalHandlerLoadError` to create a mock implementation of the
// `services.IEval` interface that does not require a knowledge base or version to be loaded.
func (s EvalServiceTestEvalHandlerLoadError) GetKnowledgeLibrary() *ast.KnowledgeLibrary {
	return ast.NewKnowledgeLibrary()
}

func (s EvalServiceTestEvalHandlerLoadError) GetKnowledgeBase(knowledgeBaseName string, version string) (*ast.KnowledgeBase, *errors.RequestError) {
	return nil, &errors.RequestError{Message: "Error on load knowledgeBase and/or version", StatusCode: 500}
}

// This is a test function that checks if the EvalHandler function returns an internal server
// error and a specific error message when there is an error loading knowledge base and/or version.
func TestEvalHandlerLoadError(t *testing.T) {

	services.EvalService = EvalServiceTestEvalHandlerLoadError{
		t: t,
	}

	c, r := mockGin()
	EvalHandler()(c)
	gotStatus := r.Code
	expectedStatus := http.StatusInternalServerError

	if gotStatus != expectedStatus {
		t.Error("got error on request evalHandler func")
	}

	gotBody := r.Body.String()
	expectedBody := "Error on load knowledgeBase and/or version"

	if gotBody != expectedBody {
		t.Error("got error on request evalHandler func")
	}
}

// This is a struct that implements the IEval interface and has a knowledge library and testing object
// as fields.
//
// Property:
//   - `EvalServiceTestEvalHandlerWithDefaultKnowledgeBase` is a struct type that implements the `services.IEval` interface.
//   - kl: is a pointer to an instance of `ast.KnowledgeLibrary`. It is likely used to store and manage knowledge bases for the evaluation service.
//   - t: is a pointer to a testing.T object, which is used for logging and reporting test results in Go unit tests. It is commonly used to call methods such as t.Errorf() to report test failures.
type EvalServiceTestEvalHandlerWithDefaultKnowledgeBase struct {
	services.IEval
	kl *ast.KnowledgeLibrary
	t  *testing.T
}

// This method takes two parameters, `knowledgeBaseName` and `version`, but it does not perform any action and
// always returns `nil`.
func (s EvalServiceTestEvalHandlerWithDefaultKnowledgeBase) LoadRemoteGRL(knowledgeBaseName string, version string) error {
	return nil
}

// This method returns a pointer to an `ast.KnowledgeLibrary` object which is a knowledge base used for
// evaluating rules. The `s.kl` variable is assumed to be a pointer to an instance of `ast.KnowledgeLibrary`
// which is set else where in the code.
func (s EvalServiceTestEvalHandlerWithDefaultKnowledgeBase) GetKnowledgeLibrary() *ast.KnowledgeLibrary {
	return s.kl
}

// This method takes in a `Context` and an `ast.KnowledgeBase` as parameters and returns a `Result`
// and an error. However, the implementation of the method simply returns a new empty `Result` and
// a `nil` error, without actually doing any evaluation.
func (s EvalServiceTestEvalHandlerWithDefaultKnowledgeBase) Eval(ctx *types.Context, knowledgeBase *ast.KnowledgeBase) (*types.Result, error) {
	return types.NewResult(), nil
}

// Thiss a test function for the EvalHandler function with a default knowledge base.
func (s EvalServiceTestEvalHandlerWithDefaultKnowledgeBase) GetKnowledgeBase(knowledgeBaseName string, version string) (*ast.KnowledgeBase, *errors.RequestError) {
	return s.kl.GetKnowledgeBase(knowledgeBaseName, version), nil
}

func TestEvalHandlerWithDefaultKnowledgeBase(t *testing.T) {

	services.EvalService = EvalServiceTestEvalHandlerWithDefaultKnowledgeBase{
		t:  t,
		kl: ast.NewKnowledgeLibrary(),
	}

	drls := `
		rule DefaultValues salience 10 {
			when 
				true
			then
				Retract("DefaultValues");
		}
	`

	ruleBuilder := builder.NewRuleBuilder(services.EvalService.GetKnowledgeLibrary())
	bs := pkg.NewBytesResource([]byte(drls))
	ruleBuilder.BuildRuleFromResource(services.DefaultKnowledgeBaseName, services.DefaultKnowledgeBaseVersion, bs)

	c, r := mockGin()

	stringReader := strings.NewReader("{}")
	c.Request.Body = io.NopCloser(stringReader)

	EvalHandler()(c)
	gotStatus := r.Code
	expectedStatus := http.StatusOK

	if gotStatus != expectedStatus {
		t.Error("got error on request evalHandler func")
	}

	gotBody := r.Body.String()
	expectedBody := "{}"

	if gotBody != expectedBody {
		t.Error("got error on request evalHandler func")
	}
}

// This is a struct that implements the IEval interface and includes a knowledge library and a testing
// object.
//
// Property:
// `EvalServiceTestEvalHandlerWithDefaultKnowledgeBaseAndWrongJSON` is a struct type
// that has three properties.
//   - kl:  is a pointer to an instance of `ast.KnowledgeLibrary`. It is likely used to store and manage knowledge bases for the evaluation service.
//   - t: is a pointer to the testing.T struct, which is used for logging and reporting test. Its typically passed as an argument to test functions.
type EvalServiceTestEvalHandlerWithDefaultKnowledgeBaseAndWrongJSON struct {
	services.IEval
	kl *ast.KnowledgeLibrary
	t  *testing.T
}

// This method takes in two parameters `knowledgeBaseName` and `version` of type string and returns an
// error. In this implementation, the method does not perform any action and simply returns a nil
// error.
func (s EvalServiceTestEvalHandlerWithDefaultKnowledgeBaseAndWrongJSON) LoadRemoteGRL(knowledgeBaseName string, version string) error {
	return nil
}

// This method returns a pointer to an `ast.KnowledgeLibrary` object which is a part of the same struct.
func (s EvalServiceTestEvalHandlerWithDefaultKnowledgeBaseAndWrongJSON) GetKnowledgeLibrary() *ast.KnowledgeLibrary {
	return s.kl
}

// This method takes in a `Context` and a `KnowledgeBase` as arguments and returns a `Result` and an error. However, the implementation
// of the method is not doing anything useful as it simply returns an empty `Result` and a `nil` error.
func (s EvalServiceTestEvalHandlerWithDefaultKnowledgeBaseAndWrongJSON) Eval(ctx *types.Context, knowledgeBase *ast.KnowledgeBase) (*types.Result, error) {
	return types.NewResult(), nil
}

// This is a test that tests the EvalHandler function with a default knowledge base and wrong JSON input.
func (s EvalServiceTestEvalHandlerWithDefaultKnowledgeBaseAndWrongJSON) GetKnowledgeBase(knowledgeBaseName string, version string) (*ast.KnowledgeBase, *errors.RequestError) {
	return s.kl.GetKnowledgeBase(knowledgeBaseName, version), nil
}

func TestEvalHandlerWithDefaultKnowledgeBaseAndWrongJSON(t *testing.T) {

	services.EvalService = EvalServiceTestEvalHandlerWithDefaultKnowledgeBaseAndWrongJSON{
		t:  t,
		kl: ast.NewKnowledgeLibrary(),
	}

	drls := `
		rule DefaultValues salience 10 {
			when 
				true
			then
				Retract("DefaultValues");
		}
	`

	ruleBuilder := builder.NewRuleBuilder(services.EvalService.GetKnowledgeLibrary())
	bs := pkg.NewBytesResource([]byte(drls))
	ruleBuilder.BuildRuleFromResource(services.DefaultKnowledgeBaseName, services.DefaultKnowledgeBaseVersion, bs)

	c, r := mockGin()

	stringReader := strings.NewReader("")
	c.Request.Body = io.NopCloser(stringReader)

	EvalHandler()(c)
	gotStatus := r.Code
	expectedStatus := http.StatusInternalServerError

	if gotStatus != expectedStatus {
		t.Error("got error on request evalHandler func")
	}

	gotBody := r.Body.String()
	expectedBody := "Error on json decode"

	if gotBody != expectedBody {
		t.Error("we expect error and the didn't came out")
	}
}

// This is a struct that implements the IEval interface and includes a knowledge library and a testing
// object.
//
// Property:
//   - 1. `EvalServiceTestEvalHandlerWithDefaultKnowledgeBaseEvalError` is a struct type that implements the `services.IEval` interface.
//   - kl: is a pointer to an instance of `ast.KnowledgeLibrary`. It is likely used to store and manage knowledge bases for the evaluation service.
//   - t: this property is a pointer to a testing.T object, which is used for logging and reporting test results in Go unit tests. It is typically passed as an argument to test functions.
type EvalServiceTestEvalHandlerWithDefaultKnowledgeBaseEvalError struct {
	services.IEval
	kl *ast.KnowledgeLibrary
	t  *testing.T
}

// LoadRemoteGRL is a struct `EvalServiceTestEvalHandlerWithDefaultKnowledgeBaseEvalError`, this
// method takes in two parameters `knowledgeBaseName` and `version` of type string and returns an
// error. In this implementation, the method always returns `nil`, indicating that there was no error
// in loading the remote knowledge base.
func (s EvalServiceTestEvalHandlerWithDefaultKnowledgeBaseEvalError) LoadRemoteGRL(knowledgeBaseName string, version string) error {
	return nil
}

// GetKnowledgeLibrary is a struct called `EvalServiceTestEvalHandlerWithDefaultKnowledgeBaseEvalError`, t
// his method returns a pointer to an `ast.KnowledgeLibrary` object which is a part of the same struct.
func (s EvalServiceTestEvalHandlerWithDefaultKnowledgeBaseEvalError) GetKnowledgeLibrary() *ast.KnowledgeLibrary {
	return s.kl
}

// This function is implementing the `Eval` method of the `services.IEval` interface for testing
// purposes. It returns a `nil` result and a mock error message when called, which is used to simulate
// an error that might occur during the evaluation of a rule. This function is used in the test case
// `TestEvalHandlerWithDefaultKnowledgeBaseEvalError` to create a mock implementation of the
// `services.IEval` interface that returns an error when the `EvalHandler` function is called with a
// default knowledge base.
func (s EvalServiceTestEvalHandlerWithDefaultKnowledgeBaseEvalError) Eval(ctx *types.Context, knowledgeBase *ast.KnowledgeBase) (*types.Result, error) {
	return nil, fmt.Errorf("mock error")
}

// This is a test that tests the EvalHandler function with a default knowledge base and
// an evaluation error.
func (s EvalServiceTestEvalHandlerWithDefaultKnowledgeBaseEvalError) GetKnowledgeBase(knowledgeBaseName string, version string) (*ast.KnowledgeBase, *errors.RequestError) {
	return s.kl.GetKnowledgeBase(knowledgeBaseName, version), nil
}
func TestEvalHandlerWithDefaultKnowledgeBaseEvalError(t *testing.T) {

	services.EvalService = EvalServiceTestEvalHandlerWithDefaultKnowledgeBaseEvalError{
		t:  t,
		kl: ast.NewKnowledgeLibrary(),
	}

	drls := `
		rule DefaultValues salience 10 {
			when 
				true
			then
				Retract("DefaultValues");
		}
	`

	ruleBuilder := builder.NewRuleBuilder(services.EvalService.GetKnowledgeLibrary())
	bs := pkg.NewBytesResource([]byte(drls))
	ruleBuilder.BuildRuleFromResource(services.DefaultKnowledgeBaseName, services.DefaultKnowledgeBaseVersion, bs)

	c, r := mockGin()

	stringReader := strings.NewReader("{}")
	c.Request.Body = io.NopCloser(stringReader)

	EvalHandler()(c)
	gotStatus := r.Code
	expectedStatus := http.StatusInternalServerError

	if gotStatus != expectedStatus {
		t.Error("got error on request evalHandler func")
	}

	gotBody := r.Body.String()
	expectedBody := "Error on eval"

	if gotBody != expectedBody {
		t.Error("we expect error and the didn't came out")
	}

}

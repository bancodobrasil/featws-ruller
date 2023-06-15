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
	"github.com/bancodobrasil/featws-ruller/services"
	"github.com/bancodobrasil/featws-ruller/types"
	"github.com/gin-gonic/gin"
	"github.com/hyperjumptech/grule-rule-engine/ast"
	"github.com/hyperjumptech/grule-rule-engine/builder"
	"github.com/hyperjumptech/grule-rule-engine/pkg"
)

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

type EvalServiceTestEvalHandlerWithoutKnowledgeBaseAndVersion struct {
	services.IEval
	t *testing.T
}

func (s EvalServiceTestEvalHandlerWithoutKnowledgeBaseAndVersion) LoadRemoteGRL(knowledgeBaseName string, version string) error {
	if knowledgeBaseName != services.DefaultKnowledgeBaseName || version != services.DefaultKnowledgeBaseVersion {
		s.t.Error("Did not load the default")
	}
	return nil
}

func (s EvalServiceTestEvalHandlerWithoutKnowledgeBaseAndVersion) GetKnowledgeLibrary() *ast.KnowledgeLibrary {
	return ast.NewKnowledgeLibrary()
}
func (s EvalServiceTestEvalHandlerWithoutKnowledgeBaseAndVersion) GetKnowledgeBase(knowledgeBaseName string, version string) (*ast.KnowledgeBase, *errors.RequestError) {
	return nil, &errors.RequestError{Message: "KnowledgeBase or version not found", StatusCode: 404}
}
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

type EvalServiceTestEvalHandlerLoadError struct {
	services.IEval
	t *testing.T
}

func (s EvalServiceTestEvalHandlerLoadError) LoadRemoteGRL(knowledgeBaseName string, version string) error {
	return fmt.Errorf("mock load error")
}

func (s EvalServiceTestEvalHandlerLoadError) GetKnowledgeLibrary() *ast.KnowledgeLibrary {
	return ast.NewKnowledgeLibrary()
}
func (s EvalServiceTestEvalHandlerLoadError) GetKnowledgeBase(knowledgeBaseName string, version string) (*ast.KnowledgeBase, *errors.RequestError) {
	return nil, &errors.RequestError{Message: "Error on load knowledgeBase and/or version", StatusCode: 500}
}

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

type EvalServiceTestEvalHandlerWithDefaultKnowledgeBase struct {
	services.IEval
	kl *ast.KnowledgeLibrary
	t  *testing.T
}

func (s EvalServiceTestEvalHandlerWithDefaultKnowledgeBase) LoadRemoteGRL(knowledgeBaseName string, version string) error {
	return nil
}

func (s EvalServiceTestEvalHandlerWithDefaultKnowledgeBase) GetKnowledgeLibrary() *ast.KnowledgeLibrary {
	return s.kl
}

func (s EvalServiceTestEvalHandlerWithDefaultKnowledgeBase) Eval(ctx *types.Context, knowledgeBase *ast.KnowledgeBase) (*types.Result, error) {
	return types.NewResult(), nil
}
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

type EvalServiceTestEvalHandlerWithDefaultKnowledgeBaseAndWrongJSON struct {
	services.IEval
	kl *ast.KnowledgeLibrary
	t  *testing.T
}

func (s EvalServiceTestEvalHandlerWithDefaultKnowledgeBaseAndWrongJSON) LoadRemoteGRL(knowledgeBaseName string, version string) error {
	return nil
}

func (s EvalServiceTestEvalHandlerWithDefaultKnowledgeBaseAndWrongJSON) GetKnowledgeLibrary() *ast.KnowledgeLibrary {
	return s.kl
}

func (s EvalServiceTestEvalHandlerWithDefaultKnowledgeBaseAndWrongJSON) Eval(ctx *types.Context, knowledgeBase *ast.KnowledgeBase) (*types.Result, error) {
	return types.NewResult(), nil
}
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

type EvalServiceTestEvalHandlerWithDefaultKnowledgeBaseEvalError struct {
	services.IEval
	kl *ast.KnowledgeLibrary
	t  *testing.T
}

func (s EvalServiceTestEvalHandlerWithDefaultKnowledgeBaseEvalError) LoadRemoteGRL(knowledgeBaseName string, version string) error {
	return nil
}

func (s EvalServiceTestEvalHandlerWithDefaultKnowledgeBaseEvalError) GetKnowledgeLibrary() *ast.KnowledgeLibrary {
	return s.kl
}

func (s EvalServiceTestEvalHandlerWithDefaultKnowledgeBaseEvalError) Eval(ctx *types.Context, knowledgeBase *ast.KnowledgeBase) (*types.Result, error) {
	return nil, fmt.Errorf("mock error")
}
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

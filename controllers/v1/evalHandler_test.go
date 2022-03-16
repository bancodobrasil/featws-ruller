package v1

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

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
	if knowledgeBaseName != DefaultKnowledgeBaseName || version != DefaultKnowledgeBaseVersion {
		s.t.Error("Did not load the default")
	}
	return nil
}

func (s EvalServiceTestEvalHandlerWithoutKnowledgeBaseAndVersion) GetKnowledgeLibrary() *ast.KnowledgeLibrary {
	return ast.NewKnowledgeLibrary()
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
	expectedBody := "KnowledgeBase or version not founded!"

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
	ruleBuilder.BuildRuleFromResource(DefaultKnowledgeBaseName, DefaultKnowledgeBaseVersion, bs)

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

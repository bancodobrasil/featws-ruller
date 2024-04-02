package services

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"text/template"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/hyperjumptech/grule-rule-engine/ast"
	"github.com/hyperjumptech/grule-rule-engine/builder"
	"github.com/hyperjumptech/grule-rule-engine/engine"
	"github.com/hyperjumptech/grule-rule-engine/pkg"

	"github.com/bancodobrasil/featws-ruller/common/errors"
	"github.com/bancodobrasil/featws-ruller/config"
	"github.com/bancodobrasil/featws-ruller/processor"
	"github.com/bancodobrasil/featws-ruller/types"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// DefaultKnowledgeBaseName its default name of Knowledge Base
const DefaultKnowledgeBaseName = "default"

// DefaultKnowledgeBaseVersion its default version of Knowledge Base
const DefaultKnowledgeBaseVersion = "latest"

// LoadLocalGRL loads a GRL (Grule Rule Language) file from a local path and builds a rule from it
// using the `builder.NewRuleBuilder` function. It takes in the path of the GRL file, the name of the
// knowledge base, and the version of the knowledge base as parameters. It returns an error if there is
// any issue building the rule from the resource.
func (s Eval) LoadLocalGRL(grlPath string, knowledgeBaseName string, version string) error {
	ruleBuilder := builder.NewRuleBuilder(s.knowledgeLibrary)
	fileRes := pkg.NewFileResource(grlPath)
	return ruleBuilder.BuildRuleFromResource(knowledgeBaseName, version, fileRes)
}

// The type `knowledgeBaseInfo` contains information about a knowledge base (rule sheet), including its name and version.
//
// Property
//   - KnowledgeBaseName: The property is a string that represents the name of a knowledge base. A knowledge base is a repository of information used to support decision-making, problem-solving, and other activities. In this case, each knowledge base is a rule sheet and contains the rules within it.
//   - Version: It's a string that represents the version number of the rule sheet, that is, the version of the knowledge base you want to use.
type knowledgeBaseInfo struct {
	KnowledgeBaseName string
	Version           string
}

// LoadRemoteGRL function is responsible for loading GRL (Grule Rule Language) rules from a remote location, such as a GitLab repository,
// and constructing a rule from them using the builder.NewRuleBuilder function. It takes the knowledge base name (rulesheet) and the
// knowledge base version as parameters.
func (s Eval) LoadRemoteGRL(ctx context.Context, knowledgeBaseName string, version string) error {
	cfg := config.GetConfig()
	ruleBuilder := builder.NewRuleBuilder(s.knowledgeLibrary)

	if cfg.ResourceLoader.Type == "http" {
		urlGRL := cfg.ResourceLoader.HTTP.URL
		urlGRL = strings.Replace(urlGRL, "{knowledgeBase}", "{{.KnowledgeBaseName}}", -1)
		urlGRL = strings.Replace(urlGRL, "{version}", "{{.Version}}", -1)

		info := knowledgeBaseInfo{
			KnowledgeBaseName: knowledgeBaseName,
			Version:           version,
		}

		urlTemplate := template.New("UrlTemplate")

		// "Parse" parses a string into a template
		urlTemplate, _ = urlTemplate.Parse(urlGRL)

		var doc bytes.Buffer
		// standard output to print merged data
		err := urlTemplate.Execute(&doc, info)
		if err != nil {
			log.Error("error on load Remote GRL: %w", err)
			return err
		}

		urlGRL = doc.String()
		hearders := cfg.ResourceLoader.HTTP.Headers

		fileRes := pkg.NewURLResourceWithHeaders(urlGRL, hearders)
		return ruleBuilder.BuildRuleFromResource(knowledgeBaseName, version, fileRes)
	}

	if cfg.ResourceLoader.Type == "minio" {

		minioClient, err := minio.New(cfg.ResourceLoader.Minio.Endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(cfg.ResourceLoader.Minio.AccessKey, cfg.ResourceLoader.Minio.SecretKey, ""),
			Secure: cfg.ResourceLoader.Minio.UseSSL,
		})
		if err != nil {
			log.Error("error on create minio client: %w", err)
			return err
		}

		path := cfg.ResourceLoader.Minio.PathTemplate

		path = strings.Replace(path, "{knowledgeBase}", "{{.KnowledgeBaseName}}", -1)
		path = strings.Replace(path, "{version}", "{{.Version}}", -1)

		info := knowledgeBaseInfo{
			KnowledgeBaseName: knowledgeBaseName,
			Version:           version,
		}

		pathTemplate := template.New("PathTemplate")

		// "Parse" parses a string into a template
		pathTemplate, _ = pathTemplate.Parse(path)

		var doc bytes.Buffer
		// standard output to print merged data
		err = pathTemplate.Execute(&doc, info)
		if err != nil {
			log.Error("error on load Remote GRL: %w", err)
			return err
		}

		path = doc.String()

		opts := minio.GetObjectOptions{}

		obj, err := minioClient.GetObject(ctx, cfg.ResourceLoader.Minio.Bucket, path, opts)
		if err != nil {
			log.Error("error on prepare a presigned url: %w", err)
			return err
		}

		res := pkg.NewReaderResource(obj)
		return ruleBuilder.BuildRuleFromResource(knowledgeBaseName, version, res)
	}

	err := fmt.Errorf("error on resolve couldn't load")
	log.Error(err)
	return err
}

// evalMutex is a variable of type sync.Mutex. This variable is used to synchronize access to the Eval method of the Eval struct, which is concurrently
// called by multiple goroutines. By using a mutex, it ensures that only one goroutine can execute the Eval method at a time, preventing race conditions and ensuring proper execution of the method.
var evalMutex sync.Mutex

var getMutex sync.Mutex

var loadWg sync.WaitGroup

var evalWg sync.WaitGroup

// IEval interface defines methods for loading and evaluating knowledge bases in Go.
//
// Property
//   - GetKnowledgeLibrary: is a method that returns a pointer to an ast.KnowledgeLibrary object. This object represents a collection of knowledge bases and their associated rules and facts.
//   - GetDefaultKnowledgeBase: is a method of the IEval interface that returns the default knowledge base of the implementation. A knowledge base is a collection of rules and facts that are used to make inferences and deductions. The default knowledge base is the one that is used if no specific knowledge base is provided during
//   - LoadLocalGRL: is a method that loads a GRL (Guideline Representation Language) file from the local file system and adds its contents to a specified knowledge base with a given version. The method takes in the path of the GRL file, the name of the knowledge base, and the version
//   - {error} LoadRemoteGRL - LoadRemoteGRL is a method that loads a GRL (Guideline Representation Language) file from a remote location into the knowledge base specified by the knowledgeBaseName and version parameters. This method is used to retrieve the rules and facts from a remote source and add them to the knowledge base for evaluation
//   - Eval - Eval is a method that takes in a context and a knowledge base and evaluates the rules in the knowledge base based on the context. It returns a result and an error if there was an issue during evaluation.
type IEval interface {
	GetKnowledgeLibrary() *ast.KnowledgeLibrary
	GetDefaultKnowledgeBase() *ast.KnowledgeBase
	GetKnowledgeBase(ctx context.Context, knowledgeBaseName string, version string) (*ast.KnowledgeBase, *errors.RequestError)
	LoadLocalGRL(grlPath string, knowledgeBaseName string, version string) error
	LoadRemoteGRL(ctx context.Context, knowledgeBaseName string, version string) error
	Eval(ctx *types.Context, knowledgeBase *ast.KnowledgeBase) (*types.Result, error)
}

// EvalService is a variable type of `IEval` and initializing it with a new instance of the `Eval` struct created by calling the `NewEval()`
// function. This variable can be used to access the methods defined in the `IEval` interface.
var EvalService IEval = NewEval(config.GetConfig())

// Eval type contains a reference to a knowledge library in Go's abstract syntax tree.
//
// Property:
//   - knowledgeLibrary - `knowledgeLibrary` is a pointer to an `ast.KnowledgeLibrary` object. Itis a property of the `Eval` struct.
type Eval struct {
	knowledgeLibrary *ast.KnowledgeLibrary
	expirationMap    map[knowledgeBaseInfo]time.Time
	versionTTL       int64
}

// NewEval  creates a new instance of the Eval struct with an empty knowledge library.
func NewEval(config *config.Config) Eval {
	return Eval{
		knowledgeLibrary: ast.NewKnowledgeLibrary(),
		expirationMap:    map[knowledgeBaseInfo]time.Time{},
		versionTTL:       config.KnowledgeBaseVersionTTL,
	}
}

// GetKnowledgeLibrary function is a method of the `Eval` struct that returns a pointer to an
// `ast.KnowledgeLibrary` object. This object represents a collection of knowledge bases and their
// associated rules and facts. The `knowledgeLibrary` property of the `Eval` struct is a pointer to an
// `ast.KnowledgeLibrary` object, and this function simply returns that pointer. This function allows
// other parts of the code to access the `ast.KnowledgeLibrary` object stored in the `Eval` struct.
func (s Eval) GetKnowledgeLibrary() *ast.KnowledgeLibrary {
	return s.knowledgeLibrary
}

// GetDefaultKnowledgeBase function is a method of the `Eval` struct that returns a pointer to
// the default knowledge base of the implementation. A knowledge base is a collection of rules and
// facts that are used to make inferences and deductions. The default knowledge base is the one that is
// used if no specific knowledge base is provided during evaluation.
func (s Eval) GetDefaultKnowledgeBase() *ast.KnowledgeBase {

	return s.GetKnowledgeLibrary().GetKnowledgeBase(DefaultKnowledgeBaseName, DefaultKnowledgeBaseVersion)
}

// GetKnowledgeBase is a method in the `Eval` struct that retrieves a knowledge base handling a possible
// expiration in the rulesheet if it reach the expiration date or loads it from a remote source if it is
// not found in the cache. It takes in the name and version of
// the knowledge base as parameters and returns a pointer to the `ast.KnowledgeBase` struct and a
// `*errors.RequestError` if there is an error. The method first checks if the knowledge base is expired.
// If it does not exist or has expired, it loads the knowledge
// base from a remote source using the `LoadRemoteGRL`
func (s Eval) GetKnowledgeBase(ctx context.Context, knowledgeBaseName string, version string) (*ast.KnowledgeBase, *errors.RequestError) {

	getMutex.Lock()
	defer getMutex.Unlock()

	info := knowledgeBaseInfo{KnowledgeBaseName: knowledgeBaseName, Version: version}

	base := s.GetKnowledgeLibrary().GetKnowledgeBase(knowledgeBaseName, version)

	expirable := true

	// If the version is a number, it isn't expirable
	if _, err := strconv.Atoi(version); err == nil {
		expirable = false
	}

	expired := expirable && s.isKnowledgeBaseVersionExpired(info)

	// If the version isn't expired and there are rules, we must retrieve the version
	if !expired && len(base.RuleEntries) > 0 {
		return base, nil
	}

	loadWg.Add(1)
	defer loadWg.Done()

	log.Trace("Waiting: [evalWg]")
	evalWg.Wait()
	log.Trace("Pass: [evalWg]")

	// If the version is expired, we must invalidate its rules
	if expired {
		for key := range base.RuleEntries {
			s.knowledgeLibrary.RemoveRuleEntry(key, knowledgeBaseName, version)
		}
	}

	err := s.LoadRemoteGRL(ctx, knowledgeBaseName, version)

	if err != nil {
		log.Errorf("Erro on load: %v", err)
		return nil, &errors.RequestError{Message: "Error on load KnowledgeBase and/or version", StatusCode: 500}
	}

	base = s.GetKnowledgeLibrary().GetKnowledgeBase(knowledgeBaseName, version)

	if len(base.RuleEntries) == 0 {
		return nil, &errors.RequestError{Message: "KnowledgeBase or version not found", StatusCode: 404}
	}

	if expirable {
		s.expirationMap[info] = time.Now().Add(time.Duration(s.versionTTL) * time.Second)
	}

	return base, nil

}

// Eval ...
func (s Eval) Eval(ctx *types.Context, knowledgeBase *ast.KnowledgeBase) (result *types.Result, err error) {

	// FIXME Remove synchronization on eval
	evalMutex.Lock()
	defer evalMutex.Unlock()

	log.Trace("Waiting: [loadWg]")
	loadWg.Wait()
	log.Trace("Pass: [loadWg]")

	evalWg.Add(1)
	defer evalWg.Done()

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recovered from panic: %v", r)
			log.Error(err)
		}
	}()
	dataCtx := ast.NewDataContext()

	processor := processor.NewProcessor()

	result = types.NewResult()

	err = dataCtx.Add("processor", processor)
	if err != nil {
		log.Error("error on add processor to data context: \n the result was: %w \n the error was: %w", result, err)
		return
	}

	err = dataCtx.Add("ctx", ctx)
	if err != nil {
		log.Error("error on add context to data context: \n the result was: %w \n the error was: %w", result, err)
		return
	}

	err = dataCtx.Add("result", result)
	if err != nil {
		log.Error("error on add result to data context: \n the result was: %w \n the error was: %w", result, err)
		return
	}

	eng := engine.NewGruleEngine()
	err = eng.Execute(dataCtx, knowledgeBase)
	if err != nil {
		log.Error("error on execute the grule engine: %w", err)
		return
	}

	if ctx.Has("errors") && len(ctx.GetMap("errors").GetEntries()) > 0 {
		result.Put("errors", ctx.GetMap("errors").GetEntries())
	}

	if ctx.Has("requiredParamErrors") && len(ctx.GetMap("requiredParamErrors").GetEntries()) > 0 {
		result.Put("requiredParamErrors", ctx.GetMap("requiredParamErrors").GetEntries())
	}

	log.Trace("Context:\n\t", ctx.GetEntries(), "\n\n")
	log.Trace("Features:\n\t", result.GetFeatures(), "\n\n")

	return
}

func (s Eval) isKnowledgeBaseVersionExpired(info knowledgeBaseInfo) bool {

	expireDate, ok := s.expirationMap[info]

	if !ok {
		return false
	}

	if expireDate.After(time.Now()) {
		return false
	}

	return true
}

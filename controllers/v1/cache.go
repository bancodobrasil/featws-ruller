package v1

import (
	"time"

	"github.com/bancodobrasil/featws-ruller/services"
	"github.com/hyperjumptech/grule-rule-engine/ast"
)

type cache struct {
	KnowledgeBase     *ast.KnowledgeBase
	KnowledgeBaseName string
	Version           string
	ExpirationDate    time.Time
}

func getCache(knowledgeBaseName string, version string) cache {
	c := cache{}

	c.KnowledgeBase = services.EvalService.GetKnowledgeLibrary().GetKnowledgeBase(knowledgeBaseName, version)
	c.ExpirationDate = time.Now().Add(time.Minute * 5)

	return c
}

func isValid(c *cache) bool {
	if c.KnowledgeBase == nil {
		return false
	}
	if c.ExpirationDate.Before(time.Now()) {
		return true
	}
	return false
}

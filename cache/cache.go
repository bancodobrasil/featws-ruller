package cache

import (
	"time"

	"github.com/bancodobrasil/featws-ruller/services"
	"github.com/hyperjumptech/grule-rule-engine/ast"
)

type Cache struct {
	KnowledgeBase     *ast.KnowledgeBase
	KnowledgeBaseName string
	Version           string
	ExpirationDate    time.Time
}

func GetCache(knowledgeBaseName string, version string) Cache {
	c := Cache{}

	c.KnowledgeBase = services.EvalService.GetKnowledgeLibrary().GetKnowledgeBase(knowledgeBaseName, version)
	c.ExpirationDate = time.Now().Add(time.Minute * 5)

	return c
}

func IsValid(c *Cache) bool {
	if c.KnowledgeBase == nil {
		return false
	}
	if c.ExpirationDate.Before(time.Now()) {
		return true
	}
	return false
}

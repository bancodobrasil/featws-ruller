package types

import (
	"fmt"
	"strconv"
)

type Context struct {
	entries map[string]interface{}
}

func NewContext() *Context {
	return &Context{
		entries: make(map[string]interface{}),
	}
}

func (c *Context) Put(param string, value interface{}) {
	c.entries[param] = value
}

func (c *Context) Get(param string) interface{} {
	return c.entries[param]
}

func (c *Context) GetString(param string) string {
	return fmt.Sprintf("%s", c.entries[param])
}

func (c *Context) GetInt(param string) int {
	value, _ := strconv.Atoi(c.GetString(param))
	return value
}

func (c *Context) GetBool(param string) bool {
	value, _ := strconv.ParseBool(c.GetString(param))
	return value
}

func (c *Context) GetEntries() map[string]interface{} {
	return c.entries
}

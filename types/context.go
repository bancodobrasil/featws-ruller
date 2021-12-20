package types

import "strconv"

type Context struct {
	entries map[string]string
}

func NewContext() *Context {
	return &Context{
		entries: make(map[string]string),
	}
}

func (c *Context) Put(param string, value string) {
	c.entries[param] = value
}

func (c *Context) Get(param string) string {
	return c.entries[param]
}

func (c *Context) GetInt(param string) int {
	value, _ := strconv.Atoi(c.Get(param))
	return value
}

func (c *Context) GetBool(param string) bool {
	value, _ := strconv.ParseBool(c.Get(param))
	return value
}

func (c *Context) GetEntries() map[string]string {
	return c.entries
}

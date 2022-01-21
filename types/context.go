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

func (c *Context) Has(param string) bool {
	_, exists := c.entries[param]
	return exists
}

func (c *Context) Get(param string) interface{} {
	return c.entries[param]
}

func (c *Context) GetSlice(param string) []interface{} {
	return c.Get(param).([]interface{})
}

func (c *Context) GetString(param string) string {
	return fmt.Sprintf("%s", c.entries[param])
}

func (c *Context) GetInt(param string) int {
	value, _ := strconv.Atoi(c.GetString(param))
	return value
}

func (c *Context) GetFloat(param string) float64 {
	value, _ := strconv.ParseFloat(c.GetString(param))
	return value
}

func (c *Context) GetBool(param string) bool {
	value, _ := strconv.ParseBool(c.GetString(param))
	return value
}

func (c *Context) GetEntries() map[string]interface{} {
	return c.entries
}

func (c *Context) AddItem(param string, item interface{}) []interface{} {
	if !c.Has(param) {
		c.Put(param, []interface{}{})
	}
	list := c.GetSlice(param)

	list = append(list, item)

	c.Put(param, list)

	return list
}

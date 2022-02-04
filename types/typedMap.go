package types

import (
	"fmt"
	"strconv"
)

type TypedMap struct {
	entries map[string]interface{}
}

func NewTypedMap() *TypedMap {
	return &TypedMap{
		entries: make(map[string]interface{}),
	}
}

func (c *TypedMap) Put(param string, value interface{}) {
	c.entries[param] = value
}

func (c *TypedMap) Has(param string) bool {
	_, exists := c.entries[param]
	return exists
}

func (c *TypedMap) Get(param string) interface{} {
	return c.entries[param]
}

func (c *TypedMap) GetSlice(param string) []interface{} {
	return c.Get(param).([]interface{})
}

func (c *TypedMap) GetString(param string) string {
	return fmt.Sprintf("%s", c.entries[param])
}

func (c *TypedMap) GetInt(param string) int64 {
	value := c.Get(param)
	strValue, ok := value.(string)
	if ok {
		intValue, _ := strconv.Atoi(strValue)
		return int64(intValue)
	}
	return int64(value.(int64))
}

func (c *TypedMap) GetFloat(param string) float64 {
	value, _ := strconv.ParseFloat(c.GetString(param), 64)
	return value
}

func (c *TypedMap) GetBool(param string) bool {
	value, _ := strconv.ParseBool(c.GetString(param))
	return value
}

func (c *TypedMap) GetEntries() map[string]interface{} {
	return c.entries
}

func (c *TypedMap) AddItem(param string, item interface{}) []interface{} {
	if !c.Has(param) {
		c.Put(param, []interface{}{})
	}
	list := c.GetSlice(param)

	list = append(list, item)

	c.Put(param, list)

	return list
}

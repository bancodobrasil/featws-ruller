package types

import (
	"fmt"
	"strconv"
)

// TypedMap its a map with method to gets entries with specific types
type TypedMap struct {
	entries map[string]interface{}
}

// NewTypedMap method create a new TypedMap
func NewTypedMap() *TypedMap {
	return &TypedMap{
		entries: make(map[string]interface{}),
	}
}

// Put method inserts a generic entry on map
func (c *TypedMap) Put(param string, value interface{}) {
	c.entries[param] = value
}

// Has method verify if a param exists in map
func (c *TypedMap) Has(param string) bool {
	_, exists := c.entries[param]
	return exists
}

// Get method get a generic entry of map
func (c *TypedMap) Get(param string) interface{} {
	return c.entries[param]
}

// GetSlice method get a slice entry of map
func (c *TypedMap) GetSlice(param string) []interface{} {
	return c.Get(param).([]interface{})
}

// GetString method get a string entry of map
func (c *TypedMap) GetString(param string) string {
	return fmt.Sprintf("%s", c.entries[param])
}

// GetInt method get a int entry of map
func (c *TypedMap) GetInt(param string) int64 {
	value := c.Get(param)

	if value == nil {
		return 0
	}

	switch v := value.(type) {
	case string:
		intValue, _ := strconv.Atoi(v)
		return int64(intValue)
	case int:
		return int64(v)
	case int64:
		return v
	default:
		panic("It's not possible to recover this parameter as int64")
	}
}

// GetBool method get a bool entry of map
func (c *TypedMap) GetBool(param string) bool {
	value, _ := strconv.ParseBool(c.GetString(param))
	return value
}

// GetEntries method get all entries of map
func (c *TypedMap) GetEntries() map[string]interface{} {
	return c.entries
}

// AddItem method inserts a item into a slice of map
func (c *TypedMap) AddItem(param string, item interface{}) []interface{} {
	if !c.Has(param) {
		c.Put(param, []interface{}{})
	}
	list := c.GetSlice(param)

	list = append(list, item)

	c.Put(param, list)

	return list
}

// AddItems methos insert some items into a slice of map
func (c *TypedMap) AddItems(param string, items ...interface{}) []interface{} {
	for _, item := range items {
		c.AddItem(param, item)
	}
	return c.GetSlice(param)
}

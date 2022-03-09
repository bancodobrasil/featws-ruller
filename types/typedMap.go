package types

import (
	"fmt"
	"reflect"
	"strconv"
)

type Getter interface {
	GetEntry(param string) interface{}
}

type interfaceMap map[string]interface{}

// TypedMap its a map with method to gets entries with specific types
type TypedMap struct {
	interfaceMap
	Getter
}

// NewTypedMap method create a new TypedMap
func NewTypedMap() *TypedMap {
	instance := &TypedMap{
		interfaceMap: make(interfaceMap),
	}
	return instance
}

// NewTypedMap method create a new TypedMap
func NewTypedMapFromMap(values map[string]interface{}) *TypedMap {
	instance := NewTypedMap()
	for key, value := range values {
		instance.Put(key, value)
	}
	return instance
}

// Put method inserts a generic entry on map
func (c *TypedMap) Put(param string, value interface{}) {
	c.interfaceMap[param] = value
}

// Has method verify if a param exists in map
func (c *TypedMap) Has(param string) bool {
	_, exists := c.interfaceMap[param]
	return exists
}

func (c *TypedMap) GetEntry(param string) interface{} {
	value, ok := c.interfaceMap[param]
	if !ok {
		return nil
	}
	return value
}

// Get method get a generic entry of map
func (c *TypedMap) Get(param string) interface{} {
	if c.Getter != nil {
		return c.Getter.GetEntry(param)
	}
	return c.GetEntry(param)
}

// GetSlice method get a slice entry of map
func (c *TypedMap) GetSlice(param string) []interface{} {
	return c.Get(param).([]interface{})
}

// GetString method get a string entry of map
func (c *TypedMap) GetString(param string) string {
	return fmt.Sprintf("%s", c.Get(param))
}

// GetInt method get a int entry of map
func (c *TypedMap) GetInt(param string) int64 {
	value := c.Get(param)
	strValue, ok := value.(string)
	if ok {
		intValue, _ := strconv.Atoi(strValue)
		return int64(intValue)
	}
	return int64(value.(int64))
}

// GetBool method get a bool entry of map
func (c *TypedMap) GetBool(param string) bool {
	value, _ := strconv.ParseBool(c.GetString(param))
	return value
}

// GetMap method get a TypedMap entry of map
func (c *TypedMap) GetMap(param string) *TypedMap {
	value := c.Get(param)
	if value != nil {
		result := NewTypedMap()
		v := reflect.ValueOf(value)
		if v.Kind() == reflect.Map {
			for _, key := range v.MapKeys() {
				strct := v.MapIndex(key)
				result.Put(key.String(), strct.Interface())
			}
		}
		return result
	}
	return nil
}

// GetEntries method get all entries of map
func (c *TypedMap) GetEntries() map[string]interface{} {
	return c.interfaceMap
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

package types

import (
	"fmt"
	"reflect"
	"strconv"

	log "github.com/sirupsen/logrus"
)

// Getter ...
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

// NewTypedMapFromMap method create a new TypedMap
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

// GetEntry ...
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
	switch v := c.Get(param).(type) {
	case bool:
		return strconv.FormatBool(v)
	default:
		return fmt.Sprintf("%s", v)
	}
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
		log.Panic("It's not possible to recover this parameter as int64")
		panic("It's not possible to recover this parameter as int64")
	}
}

// GetFloat method get a int entry of map
func (c *TypedMap) GetFloat(param string) float64 {
	value := c.Get(param)
	if value == nil {
		return 0
	}
	switch v := value.(type) {
	case string:
		floatValue, _ := strconv.ParseFloat(v, 64)
		return floatValue
	case int:
		return float64(v)
	case int64:
		return float64(v)
	case float64:
		return v
	default:
		log.Panic("fail to retrieve this param as float64")
		panic("fail to retrieve this param as float64")
	}
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
		v := reflect.ValueOf(value)
		kind := v.Kind()
		if kind == reflect.Map {
			result := NewTypedMap()
			for _, key := range v.MapKeys() {
				strct := v.MapIndex(key)
				result.Put(key.String(), strct.Interface())
			}
			return result
		}

		tp, ok := value.(*TypedMap)
		if ok {
			return tp
		}

		panic("This param it's not a map")
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

// AddItems methos insert some items into a slice of map
func (c *TypedMap) AddItems(param string, items ...interface{}) []interface{} {
	for _, item := range items {
		c.AddItem(param, item)
	}
	return c.GetSlice(param)
}

package processor

import (
	"strconv"
)

// Processor its utilitary class for assertions
type Processor struct{}

// NewProcessor method to create Processors
func NewProcessor() *Processor {
	return &Processor{}
}

// Boolean method to convert boolean string to boolean
func (p *Processor) Boolean(value bool) string {
	return strconv.FormatBool(value)
}

// Contains method to check if entry is into a array
func (p *Processor) Contains(slice []interface{}, val interface{}) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

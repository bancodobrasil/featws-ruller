package types

import (
	"fmt"
	"strconv"
)

type Result struct {
	features map[string]interface{}
}

func NewResult() *Result {
	return &Result{
		features: make(map[string]interface{}),
	}
}

func (r *Result) Put(feat string, value interface{}) {
	r.features[feat] = value
}

func (r *Result) Get(param string) interface{} {
	return r.features[param]
}

func (r *Result) GetString(param string) string {
	return fmt.Sprintf("%s", r.features[param])
}

func (r *Result) GetInt(param string) int {
	value, _ := strconv.Atoi(r.GetString(param))
	return value
}

func (r *Result) GetFloat(param string) float64 {
	value, _ := strconv.ParseFloat(r.GetString(param))
	return value
}

func (r *Result) GetBool(param string) bool {
	value, _ := strconv.ParseBool(r.GetString(param))
	return value
}

func (r *Result) GetFeatures() map[string]interface{} {
	return r.features
}

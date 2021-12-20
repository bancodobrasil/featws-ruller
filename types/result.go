package types

import "strconv"

type Result struct {
	features map[string]string
}

func NewResult() *Result {
	return &Result{
		features: make(map[string]string),
	}
}

func (r *Result) Put(feat string, value string) {
	r.features[feat] = value
}

func (r *Result) Get(param string) string {
	return r.features[param]
}

func (r *Result) GetInt(param string) int {
	value, _ := strconv.Atoi(r.Get(param))
	return value
}

func (r *Result) GetBool(param string) bool {
	value, _ := strconv.ParseBool(r.Get(param))
	return value
}

func (r *Result) GetFeatures() map[string]string {
	return r.features
}

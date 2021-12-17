package types

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

func (r *Result) GetFeatures() map[string]string {
	return r.features
}

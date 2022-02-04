package types

type Result struct {
	TypedMap
}

func NewResult() *Result {
	return &Result{
		TypedMap: *NewTypedMap(),
	}
}

func (r *Result) GetFeatures() map[string]interface{} {
	return r.GetEntries()
}

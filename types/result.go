package types

// Result its used to store features in end rule assertions
type Result struct {
	TypedMap
}

// NewResult method create a new Result
func NewResult() *Result {
	return &Result{
		TypedMap: *NewTypedMap(),
	}
}

// GetFeatures get list of features
func (r *Result) GetFeatures() map[string]interface{} {
	return r.GetEntries()
}

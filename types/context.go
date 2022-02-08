package types

// Context its used to store parameters and temporary variables during rule assertions
type Context struct {
	TypedMap
}

// NewContext method create a new Context
func NewContext() *Context {
	return &Context{
		TypedMap: *NewTypedMap(),
	}
}

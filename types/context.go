package types

type Context struct {
	TypedMap
}

func NewContext() *Context {
	return &Context{
		TypedMap: *NewTypedMap(),
	}
}

package types

import (
	"reflect"
	"testing"
)

func TestNewContext(t *testing.T) {
	got := NewContext()

	expected := &Context{
		TypedMap: *NewTypedMap(),
	}

	if reflect.TypeOf(got) != reflect.TypeOf(expected) {
		t.Error("You got an error while try to generate new context")
	}
}

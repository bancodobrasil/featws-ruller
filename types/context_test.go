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

func TestRegistryRemoteLoaded(t *testing.T) {
	ct := NewContext()
	ct.RegistryRemoteLoaded("myparam", "myresolver")
	got := ct.isRemoteLoaded("myparam")
	expected := true
	if got != expected {
		t.Errorf("Test Fail, we want %v, we got %v", got, expected)
	}

}
func TestRegistryNotRemoteLoaded(t *testing.T) {
	ct := NewContext()
	// ct.RegistryRemoteLoaded("myparam", "myresolver")
	got := ct.isRemoteLoaded("myparam")
	expected := false
	if got != expected {
		t.Errorf("Test Fail, we want %v, we got %v", got, expected)
	}

}
func TestLoad(t *testing.T) {

}

// func TestIsRemoteLoaded(t *testing.T) {
// 	ct := NewContext()
// 	param := "myparam"
// 	rlds := ct.RemoteLoadeds
// 	resolver := &RemoteLoaded{
// 		Resolver: "myresolver",
// 	}
// 	rlds[param] = *resolver

// 	got := ct.isRemoteLoaded(param)
// 	expected := true

// 	if got != expected {
// 		t.Errorf("Test Fail, we want %v, we got %v", got, expected)
// 	}
// }

func TestGetEntryContext(t *testing.T) {
	// ct := NewContext()
	// ct.Put("mystring", "test")
	// param := "mystring"
}

func TestResolve(t *testing.T) {

}

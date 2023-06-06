package types

import (
	"reflect"
	"testing"
)

// This function tests the creation of a new Result object and checks if it has the expected type.
func TestNewResult(t *testing.T) {
	got := NewResult()

	expected := &Result{
		TypedMap: *NewTypedMap(),
	}

	if reflect.TypeOf(got) != reflect.TypeOf(expected) {
		t.Error("You got an error while try to generate a new result")
	}
}

// This is a test checks if a list of features can be retrieved correctly from a result object.
func TestGetFeatures(t *testing.T) {
	result := NewResult()
	result.Put("myint", 10)
	result.Put("mystring", "test")
	result.Put("mybool", true)

	got := result.GetFeatures()

	expected := map[string]interface{}{
		"myint":    10,
		"mystring": "test",
		"mybool":   true,
	}

	if reflect.DeepEqual(got, expected) != true {
		t.Error("You got an error while try to get list of features")
	}
}

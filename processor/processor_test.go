package processor

import (
	"reflect"
	"testing"
)

// TestNewProcessor checks if a new instance of a Processor is created correctly.
func TestNewProcessor(t *testing.T) {
	got := NewProcessor()

	expected := &Processor{}

	if reflect.TypeOf(got) != reflect.TypeOf(expected) {
		t.Error("you got an error while try to generate a new Processor")

	}
}

// TestBooleanFalse checks if a boolean value is correctly converted to a string.
func TestBooleanFalse(t *testing.T) {
	got := (NewProcessor().Boolean(false))

	expect := "false"

	if got != expect {
		t.Error("you got an error while try to convert boolean to string")
	}

}

// TestBooleanTrue checks if a boolean value is correctly converted to a string.
func TestBooleanTrue(t *testing.T) {
	got := (NewProcessor().Boolean(true))

	expect := "true"

	if got != expect {
		t.Error("you got an error while try to convert boolean to string")
	}

}

// TestContainsTrue checks if a value is present in a slice using a custom processor.
func TestContainsTrue(t *testing.T) {
	testArray := []interface{}{1, 2, 3, 4}
	value := 4
	got := (NewProcessor().Contains(testArray, value))
	expect := true

	if got != expect {
		t.Error("you got an error while try to check the value inside the slice")
	}

}

// TestContainsFalse checks if a value is not present in a given array.
func TestContainsFalse(t *testing.T) {
	testArray := []interface{}{1, 2, 3, 4}
	value := 5
	got := (NewProcessor().Contains(testArray, value))
	expect := false

	if got != expect {
		t.Error("you got an error while try to check the value inside the slice")
	}

}

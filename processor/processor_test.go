package processor

import (
	"reflect"
	"testing"
)

func TestNewProcessor(t *testing.T) {
	exp_return := &Processor{}

	got := NewProcessor()

	if reflect.TypeOf(got) != reflect.TypeOf(exp_return) {
		t.Error("you got an error while try to generate a new Processor")

	}
}

func TestBooleanFalse(t *testing.T) {
	got := (NewProcessor().Boolean(false))

	expect := "false"

	if got != expect {
		t.Error("you got an error while try to convert boolean to string")
	}

}

func TestBooleanTrue(t *testing.T) {
	got := (NewProcessor().Boolean(true))

	expect := "true"

	if got != expect {
		t.Error("you got an error while try to convert boolean to string")
	}

}

func TestContainsTrue(t *testing.T) {
	testArray := []interface{}{1, 2, 3, 4}
	value := 4
	got := (NewProcessor().Contains(testArray, value))
	expect := true

	if got != expect {
		t.Error("you got an error while try to check the value inside the slice")
	}

}
func TestContainsFalse(t *testing.T) {
	testArray := []interface{}{1, 2, 3, 4}
	value := 5
	got := (NewProcessor().Contains(testArray, value))
	expect := false

	if got != expect {
		t.Error("you got an error while try to check the value inside the slice")
	}

}

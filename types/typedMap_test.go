package types

import (
	"reflect"
	"testing"
)

func TestNewTypedMap(t *testing.T) {
	got := NewTypedMap()

	expected := &TypedMap{
		entries: make(map[string]interface{}),
	}

	if reflect.DeepEqual(got, expected) != true {
		t.Error("You got an error while try to generate a new typed map")
	}

}

func TestPut(t *testing.T) {
	tm := NewTypedMap()
	tm.Put("mystring", "test")
	got := tm.GetEntries()["mystring"]

	expected := "test"

	if got != expected {
		t.Error("You got an error while try to put value into map")
	}

}

func TestHasTrue(t *testing.T) {
	tm := NewTypedMap()
	tm.Put("mystring", "test")
	got := tm.Has("mystring")

	expected := true

	if got != expected {
		t.Error("Test fail, the param is not into map")
	}
}

func TestHasFalse(t *testing.T) {
	tm := NewTypedMap()
	tm.entries["mystring"] = 1
	param := "myint"

	got := tm.Has(param)

	expected := false

	if got != expected {
		t.Error("Test fail, the param is not into map")
	}
}

func TestGet(t *testing.T) {
	tm := NewTypedMap()
	tm.entries["mystring"] = 1
	param := "mystring"

	got := tm.Get(param)

	expected := 1

	if got != expected {
		t.Error("Test fail, the param doesn't exists into map")
	}
}

func TestGetSlice(t *testing.T) {
	tm := NewTypedMap()
	myArray := []interface{}{1, 2, 3, 4, 5}
	tm.Put("myArray", myArray)
	param := "myArray"
	got := tm.GetSlice(param)
	expected := myArray

	if reflect.DeepEqual(got, expected) != true {
		t.Error("Test fail, the arrays doesn't match each other")
	}
}

func TestGetString(t *testing.T) {
	tm := NewTypedMap()
	tm.Put("mystring", "test")
	got := tm.GetString("mystring")
	expected := "test"

	if got != expected {
		t.Error("Test fail, string doesn't exist into map")
	}
}

func TestGetIntWithInteger(t *testing.T) {
	tm := NewTypedMap()
	tm.Put("myint", 10)
	got := tm.GetInt("myint")
	expected := int64(10)

	if got != expected {
		t.Error("Couldn't get the integer into the map")
	}
}

func TestGetIntWithString(t *testing.T) {
	tm := NewTypedMap()
	tm.Put("myint", "10")
	got := tm.GetInt("myint")
	expected := int64(10)

	if got != expected {
		t.Error("Couldn't get the integer into the map")
	}
}

func TestGetIntWithInt64(t *testing.T) {
	tm := NewTypedMap()
	tm.Put("myint", int64(10))
	got := tm.GetInt("myint")
	expected := int64(10)

	if got != expected {
		t.Error("Couldn't get the integer into the map")
	}
}

func TestGetIntWithNoParam(t *testing.T) {
	tm := NewTypedMap()
	tm.Put("myint", "")
	got := tm.GetInt("")
	expected := int64(0)

	if got != expected {
		t.Error("Couldn't get the integer into the map")
	}
}

func TestGetIntWithPanic(t *testing.T) {
	defer func() {
		r := recover()
		if r != "It's not possible to recover this parameter as int64" {
			t.Error("The panic message it's not throwed")
		}
	}()
	tm := NewTypedMap()
	tm.Put("myint", false)
	tm.GetInt("myint")
}

func TestGetBool(t *testing.T) {
	tm := NewTypedMap()
	tm.Put("mybool", "true")
	got := tm.GetBool("mybool")
	expect := true

	if got != expect {
		t.Error("Couldn't get the bool value into the map")
	}

}

func TestGetEntries(t *testing.T) {
	tm := NewTypedMap()
	tm.Put("mystring", "test")
	tm.Put("myint", 10)
	tm.Put("mybool", false)

	got := tm.GetEntries()

	expected := map[string]interface{}{
		"mystring": "test",
		"myint":    10,
		"mybool":   false,
	}

	if reflect.DeepEqual(got, expected) != true {
		t.Error("You got an error while try to get the entries of the map")
	}

}

func TestAddItemWithOneItem(t *testing.T) {
	tm := NewTypedMap()
	myslice := []interface{}{"test1", "test2"}
	tm.Put("myslice", myslice)

	got := tm.AddItem("myslice", "test3")
	expected := []interface{}{"test1", "test2", "test3"}

	if reflect.DeepEqual(got, expected) != true {
		t.Error("Couldn't add an item into a slice of map")
	}

}

func TestAddItemWithThreeItemsNonInitlizated(t *testing.T) {
	tm := NewTypedMap()
	tm.AddItem("myslice", "test1")
	tm.AddItem("myslice", "test2")
	tm.AddItem("myslice", "test3")
	got := tm.GetSlice("myslice")

	expected := []interface{}{"test1", "test2", "test3"}

	if reflect.DeepEqual(got, expected) != true {
		t.Error("The the expected map doesn't match with the obtained one")
	}

}

func TestAddItemWithThreeItems(t *testing.T) {
	tm := NewTypedMap()
	myslice := []interface{}{"test1", "test2"}
	tm.Put("myslice", myslice)
	tm.AddItem("myslice", "test3")
	tm.AddItem("myslice", "test4")
	tm.AddItem("myslice", "test5")
	got := tm.GetSlice("myslice")

	expected := []interface{}{"test1", "test2", "test3", "test4", "test5"}

	if reflect.DeepEqual(got, expected) != true {
		t.Error("The the expected map doesn't match with the obtained one")
	}

}

func TestAddItemsWithThreeItemsNonInitlizated(t *testing.T) {
	tm := NewTypedMap()
	got := tm.AddItems("myslice", "test1", "test2", "test3")

	expected := []interface{}{"test1", "test2", "test3"}

	if reflect.DeepEqual(got, expected) != true {
		t.Error("The the expected map doesn't match with the obtained one")
	}

}

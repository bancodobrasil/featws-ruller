package types

import (
	"reflect"
	"testing"

	"github.com/sirupsen/logrus"
)

// This is a test that checks if a new instance of a TypedMap is generated correctly.
func TestNewTypedMap(t *testing.T) {
	got := NewTypedMap()

	expected := &TypedMap{
		interfaceMap: make(map[string]interface{}),
	}

	if reflect.DeepEqual(got, expected) != true {
		t.Error("You got an error while try to generate a new typed map")
	}

}

// This tests creates a new typed map from a regular map.
func TestNewTypedMapFromMap(t *testing.T) {
	m := interfaceMap{"mystring": "test"}
	got := NewTypedMapFromMap(m).GetEntries()
	expected := m
	if reflect.DeepEqual(got, expected) {
		t.Error("You got an error while try to generate a new typed map from a map")
	}
}

// This test checks if a value can be successfully added to a typed map.
func TestPut(t *testing.T) {
	tm := NewTypedMap()
	tm.Put("mystring", "test")
	got := tm.GetEntries()["mystring"]

	expected := "test"

	if got != expected {
		t.Error("You got an error while try to put value into map")
	}
}

// This test function checks if a specific key exists in a typed map.
func TestHasTrue(t *testing.T) {
	tm := NewTypedMap()
	tm.Put("mystring", "test")
	got := tm.Has("mystring")

	expected := true

	if got != expected {
		t.Error("Test fail, the param is not into map")
	}
}

// This test checks if a specific parameter is present in a typed map and returns false if it isn't.
func TestHasFalse(t *testing.T) {
	tm := NewTypedMap()
	tm.Put("mystring", "test")
	param := "myint"

	got := tm.Has(param)

	expected := false

	if got != expected {
		t.Error("Test fail, the param is not into map")
	}
}

// This is a test that tests the Get method of a TypedMap object.
func TestGet(t *testing.T) {
	tm := NewTypedMap()
	tm.Put("mystring", "test")
	param := "mystring"

	got := tm.Get(param)

	expected := "test"

	if got != expected {
		t.Error("Test fail, the param doesn't exists into map")
	}
}

// This is a test function that tests the GetEntry method of a TypedMap object.
func TestGetEntry(t *testing.T) {
	tm := NewTypedMap()
	tm.Put("mystring", "test")
	param := "mystring"

	got := tm.GetEntry(param)

	expected := "test"

	if got != expected {
		t.Error("Test fail, the param doesn't exists into map")
	}
}

// The mockGetter type embeds the Getter interface.
//
// Property:
//   - Getter: The `mockGetter` struct can implement the methods of the `Getter` interface and also add its own methods.
type mockGetter struct {
	Getter
}

// GetEntry function is a method of the `mockGetter` struct that implements the `GetEntry` method of the `Getter` interface. It returns a
// mock value of type `interface{}`. This function is used in the `TestGetWithGetter` test function to
// test the `Get` method of the `TypedMap` object when a `Getter` interface is embedded in the `TypedMap` object.
func (m *mockGetter) GetEntry(param string) interface{} {
	return "mock"
}

// This test function in Go that tests the Get method of a TypedMap object with a mock
// getter.
func TestGetWithGetter(t *testing.T) {
	tm := NewTypedMap()
	tm.Getter = &mockGetter{}
	tm.Put("mystring", "test")
	param := "mystrifdfdng"

	got := tm.Get(param)

	expected := "mock"

	if got != expected {
		t.Error("Test fail, the param doesn't exists into map")
	}
}

// This test checks if a slice retrieved from a typed map matches an expected slice.
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

// The function tests if a string value exists in a typed map and returns an error if it doesn't.
func TestGetString(t *testing.T) {
	tm := NewTypedMap()
	tm.Put("mystring", "test")
	got := tm.GetString("mystring")
	expected := "test"

	if got != expected {
		t.Error("Test fail, string doesn't exist into map")
	}
}

// This test checks if an integer value can be retrieved from a typed map.
func TestGetIntWithInteger(t *testing.T) {
	tm := NewTypedMap()
	tm.Put("myint", 10)
	got := tm.GetInt("myint")
	expected := int64(10)

	if got != expected {
		t.Error("Couldn't get the integer into the map")
	}
}

// This tests if an integer can be retrieved from a string value in a typed map.
func TestGetIntWithString(t *testing.T) {
	tm := NewTypedMap()
	tm.Put("myint", "10")
	got := tm.GetInt("myint")
	expected := int64(10)

	if got != expected {
		t.Error("Couldn't get the integer into the map")
	}
}

// This test checks if an integer value can be retrieved from a typed map.
func TestGetIntWithInt64(t *testing.T) {
	tm := NewTypedMap()
	tm.Put("myint", int64(10))
	got := tm.GetInt("myint")
	expected := int64(10)

	if got != expected {
		t.Error("Couldn't get the integer into the map")
	}
}

// This test checks if an integer can be retrieved from a typed map with no parameter.
func TestGetIntWithNoParam(t *testing.T) {
	tm := NewTypedMap()
	tm.Put("myint", "")
	got := tm.GetInt("")
	expected := int64(0)

	if got != expected {
		t.Error("Couldn't get the integer into the map")
	}
}

// This function tests if a panic message is thrown when trying to recover a non-integer parameter as
// an int64.
func TestGetIntWithPanic(t *testing.T) {
	defer func() {
		r := recover()
		if r.(*logrus.Entry).Message != "It's not possible to recover this parameter as int64" {
			t.Error("The panic message it's not throwed")
		}
	}()
	tm := NewTypedMap()
	tm.Put("myint", false)
	tm.GetInt("myint")
}

// This test checks if a float value can be retrieved from a typed map with no parameters.
func TestGetFloatWithNoParams(t *testing.T) {
	tm := NewTypedMap()
	tm.Put("myint", "")
	got := tm.GetFloat("")
	expected := float64(0)

	if got != expected {
		t.Error("Couldn't get the float value into the map")
	}
}

// This test checks if a string value can be converted to a float and stored in a typed map.
func TestGetFloatWithString(t *testing.T) {
	tm := NewTypedMap()
	tm.Put("myFloatString", "5.5")
	got := tm.GetFloat("myFloatString")
	expected := float64(5.5)

	if got != expected {
		t.Error("Couldn't get the string into the map")
	}
}

// This test checks if an integer can be converted to a float and stored in a typed map.
func TestGetFloatWithInteger(t *testing.T) {
	tm := NewTypedMap()
	tm.Put("myFloatInteger", 1)
	got := tm.GetFloat("myFloatInteger")
	expected := float64(1.0)

	if got != expected {
		t.Error("Couldn't get the integer into the map")
	}
}

// This function tests if an int64 value can be retrieved as a float64 from a typed map.
func TestGetFloatWithInt64(t *testing.T) {
	tm := NewTypedMap()
	tm.Put("myFloatInt64", int64(12))
	got := tm.GetFloat("myFloatInt64")
	expected := float64(12.0)

	if got != expected {
		t.Error("Couldn't get the int64 into the map")
	}
}

// This test checks if a float64 value can be added to a typed map and retrieved correctly.
func TestGetFloatWithFloat64(t *testing.T) {
	tm := NewTypedMap()
	tm.Put("myFloat64", float64(12.562))
	got := tm.GetFloat("myFloat64")
	expected := float64(12.562)

	if got != expected {
		t.Error("Couldn't get the float64 into the map")
	}
}

// This function tests if a panic message is thrown when trying to retrieve a non-float parameter from
// a typed map.
func TestGetFloatWithPanic(t *testing.T) {
	defer func() {
		r := recover()
		if r.(*logrus.Entry).Message != "fail to retrieve this param as float64" {
			t.Error("The panic message it's not throwed")
		}
	}()
	tm := NewTypedMap()
	tm.Put("myWrongParam", false)
	tm.GetFloat("myWrongParam")

}

// This is a test function that checks if a boolean value can be retrieved from a typed map.
func TestGetBool(t *testing.T) {
	tm := NewTypedMap()
	tm.Put("mybool", "true")
	got := tm.GetBool("mybool")
	expect := true

	if got != expect {
		t.Error("Couldn't get the bool value into the map")
	}

}

// This checks if a TypedMap object can retrieve a map and compare it to an expected map.
func TestGetMap(t *testing.T) {
	tm := NewTypedMap()
	mymap := map[string]interface{}{
		"mystring": "test",
		"myint":    15,
		"mybool":   true,
	}

	tm.Put("mymap", mymap)

	got := tm.GetMap("mymap").GetEntries()
	expected := mymap

	if reflect.DeepEqual(got, expected) != true {
		t.Error("The maps aren't equals")
	}

}

// This is a test checks if a map retrieved from a TypedMap object is nil.
func TestGetMapNil(t *testing.T) {
	tm := NewTypedMap()
	tm.Put("mymap", "")

	got := tm.GetMap("")

	if got != nil {
		t.Error("The maps aren't equals")
	}

}

// This test checks if a panic message is thrown when trying to get a map from a non-map parameter in a typed map.
func TestGetMapWithoutMap(t *testing.T) {
	defer func() {
		r := recover()
		if r != "This param it's not a map" {
			t.Error("The panic message it's not throwed")
		}
	}()
	tm := NewTypedMap()
	tm.Put("mymap", "test")
	tm.GetMap("mymap")

}

// This is a test checks if a TypedMap can correctly retrieve its entries.
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

// The function tests adding an item to a slice in a typed map.
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

// This test checks if adding three items to a non-initialized slice in a TypedMap works correctly.
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

// This test checks if adding multiple items to a slice in a typed map works correctly.
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

// This tests the AddItems method of a TypedMap object by adding three items to a
// non-initialized slice and comparing the obtained result with the expected one.
func TestAddItemsWithThreeItemsNonInitlizated(t *testing.T) {
	tm := NewTypedMap()
	got := tm.AddItems("myslice", "test1", "test2", "test3")

	expected := []interface{}{"test1", "test2", "test3"}

	if reflect.DeepEqual(got, expected) != true {
		t.Error("The the expected map doesn't match with the obtained one")
	}

}

// This test checks if a boolean value stored in a TypedMap can be retrieved and
// converted to the expected value.
func TestBooleanPlain(t *testing.T) {
	tm := NewTypedMap()
	tm.Put("mybool", true)

	got := tm.GetBool("mybool")
	expect := true

	if got != expect {
		t.Error("you got an error while try to convert boolean to string")
	}

}

// This function tests if an integer value can be retrieved from a typed map.
func TestIntegerPlain(t *testing.T) {
	tm := NewTypedMap()
	tm.Put("mynumber", 123)

	got := tm.GetInt("mynumber")
	expect := int64(123)

	if got != expect {
		t.Error("you got an error while try to convert int to string")
	}

}

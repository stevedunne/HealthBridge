package assert

import (
	_ "fmt"
	"reflect"
	"testing"
)

// type Assert struct{
// 	T *testing.T
// }

// func NewTesting(testing *testing.T) Assert {
// 	a:=  Assert{ }
// 	a.T = testing
// 	return a
// }

// Checks if the error object is nil
func Error(err error, t *testing.T) {
	if err == nil {
		t.Log("Passed")
	} else {
		t.Errorf("Error %v ", err)
	}
}

func NotNil(actual interface{}, message string, t *testing.T) {
	obj := reflect.ValueOf(actual)
	if obj.IsNil() {
		t.Errorf("Expected not nil but got %v: %s", actual, message)
	} else {
		t.Log("Passed")
	}
}

func IsNil(actual interface{}, message string, t *testing.T) {
	obj := reflect.ValueOf(actual)
	if obj.IsNil() {
		t.Log("Passed")
	} else {
		t.Errorf("Expected nil but got %v: %s", actual, message)
	}
}

// Checks if the actual string matches the expected string
func StrEqual(expected, actual, message string, t *testing.T) {
	if expected == actual {
		t.Log("Passed")
	} else {
		t.Errorf("Expected %v but got %v: %s", expected, actual, message)
	}
}

// Checks if the actual matches the expected string
func IntEqual(expected, actual int, t *testing.T) {
	if expected == actual {
		t.Log("Passed")
	} else {
		t.Errorf("Expected %v but got %v", expected, actual)
	}
}

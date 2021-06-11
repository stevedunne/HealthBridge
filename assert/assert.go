package assert

import (
	"reflect"
	"testing"
)

// Error Checks if the error object is nil
func Error(err error, t *testing.T) {
	if err == nil {
		t.Log("Passed")
	} else {
		t.Errorf("Error %v ", err)
	}
}

//NotNil tests if an object is initialised
func NotNil(actual interface{}, message string, t *testing.T) {
	obj := reflect.ValueOf(actual)
	if obj.IsNil() {
		t.Errorf("Expected not nil but got %v: %s", actual, message)
	} else {
		t.Log("Passed")
	}
}

//IsNil tests for an uninitialised object
func IsNil(actual interface{}, message string, t *testing.T) {
	obj := reflect.ValueOf(actual)
	if obj.IsNil() {
		t.Log("Passed")
	} else {
		t.Errorf("Expected nil but got %v: %s", actual, message)
	}
}

//StrEqual Checks if the actual string matches the expected string
func StrEqual(expected, actual, message string, t *testing.T) {
	if expected == actual {
		t.Log("Passed")
	} else {
		t.Errorf("Expected %v but got %v: %s", expected, actual, message)
	}
}

//IntEqual Checks if the actual matches the expected int
func IntEqual(expected, actual int, t *testing.T) {
	if expected == actual {
		t.Log("Passed")
	} else {
		t.Errorf("Expected %v but got %v", expected, actual)
	}
}

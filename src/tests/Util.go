package tests

import (
	"errors"
	"reflect"
	"runtime/debug"
	"testing"
)

// Assert use this implementation to avoid github assert imports
func Assert(t *testing.T, exp, got interface{}, equal bool) {
	if reflect.DeepEqual(exp, got) != equal {
		debug.PrintStack()
		t.Fatalf("Expecting '%v' got '%v'\n", exp, got)
	}
}

func AssertThrow(t *testing.T, exp error, got error) {
	if !errors.As(exp, &got) {
		t.Errorf("Wrong error thrown")
	}
}

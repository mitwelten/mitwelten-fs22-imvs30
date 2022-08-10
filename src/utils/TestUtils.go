package utils

import (
	"reflect"
	"runtime/debug"
	"strings"
	"testing"
)

// Assert use this implementation to avoid github assert imports
func Assert(t *testing.T, exp, got interface{}, equal bool) {
	if reflect.DeepEqual(exp, got) != equal {
		debug.PrintStack()
		t.Fatalf("Expecting '%v' got '%v'\n", exp, got)
	}
}

func ExpectErrorContains(t *testing.T, expectedPrefix string, err error) {
	if err == nil {
		t.Fatalf("Expected error {'%v...'} but no error thrown\n", expectedPrefix)
	}
	if !strings.Contains(err.Error(), expectedPrefix) {
		t.Fatalf("Expected error {'%v...'} but got {'%v'}\n", expectedPrefix, err.Error())
	}
}

func ExpectNoError(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("Unexpected error {%v}\n", err.Error())
	}
}

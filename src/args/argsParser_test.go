package args

import (
	"errors"
	"mjpeg_multiplexer/src/connection"
	"mjpeg_multiplexer/src/customErrors"
	"testing"
)

var argsMock = []string{"main.exe", "-input", "192.168.137.216:8080 192.168.137.59:8080", "-output", "file", "-output_filename", "out.jpg", "-method", "grid"}
var expectedInputLocations []string
var expectedOutput connection.Output

func TestStreamCommandShouldNotCrash(t *testing.T) {
	// given
	argsMock = []string{"main.exe", "-input", "192.168.137.216:8080 192.168.137.59:8080", "-output", "stream", "-output_port", "8088", "-method", "grid"}

	// then
	_, _ = ParseArgs(argsMock)
}

func TestFileCommandShouldNotCrash(t *testing.T) {
	// given
	argsMock = []string{"main.exe", "-input", "192.168.137.216:8080 192.168.137.59:8080", "-output", "file", "-output_filename", "out.jpg", "-method", "grid"}

	// then
	_, _ = ParseArgs(argsMock)
}

func TestShouldFailWithNotFulfillingMinArguments(t *testing.T) {
	// given
	argsMock = []string{"main.exe", "-input", "192.168.137.216:8080 192.168.137.59:8080", "stream", "-method", "grid"}
	var expected error = &customErrors.ErrArgParserUnfulfilledMinArguments{}

	// when
	var _, err = ParseArgs(argsMock)

	// then
	if err == nil {
		t.Errorf("Error not thrown")
	}

	println(err.Error())

	if !(errors.As(err, &expected)) {
		t.Errorf("Wrong error thrown")
	}
}

func TestShouldFailWithMissingPort(t *testing.T) {
	// given
	argsMock = []string{"main.exe", "-input", "192.168.137.216:8080 192.168.137.59:8080", "-output", "stream", "-method", "grid"}
	var expected error = &customErrors.ErrArgParserInvalidOutputPort{}

	// when
	var _, err = ParseArgs(argsMock)

	// then
	if err == nil {
		t.Errorf("Error not thrown")
	}

	println(err.Error())

	if !(errors.As(err, &expected)) {
		t.Errorf("Wrong error thrown")
	}
}

func TestShouldFailWithMissingFilename(t *testing.T) {
	// given
	argsMock = []string{"main.exe", "-input", "192.168.137.216:8080 192.168.137.59:8080", "-output", "file", "-method", "grid"}
	var expected error = &customErrors.ErrArgParserInvalidOutputFilename{}

	// when
	var _, err = ParseArgs(argsMock)

	// then
	if err == nil {
		t.Errorf("Error not thrown")
	}

	println(err.Error())

	if !(errors.As(err, &expected)) {
		t.Errorf("Wrong error thrown")
	}
}

func TestShouldFailWithInvalidOutputArgument(t *testing.T) {
	// given
	argsMock = []string{"main.exe", "-input", "192.168.137.216:8080 192.168.137.59:8080", "-output", "XXXX", "-method", "grid"}
	var expected error = &customErrors.ErrArgParserInvalidArgument{}

	// when
	var _, err = ParseArgs(argsMock)

	// then
	if err == nil {
		t.Errorf("Error not thrown")
	}

	println(err.Error())

	if !(errors.As(err, &expected)) {
		t.Errorf("Wrong error thrown")
	}
}

package args

import (
	"mjpeg_multiplexer/src/connection"
	"strings"
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

func TestShouldFailWithMissingPort(t *testing.T) {
	// given
	argsMock = []string{"main.exe", "-input", "192.168.137.216:8080 192.168.137.59:8080", "-output", "stream", "-method", "grid"}

	// then
	var _, err = ParseArgs(argsMock)
	if err == nil {
		t.Errorf("Error not thrown")
	}

	if strings.Compare(err.Error(), "-output 'stream' only valid in combination with -output_port ") != 0 {
		t.Errorf("Wrong error thrown")
	}
}

func TestShouldFailWithMissingFilename(t *testing.T) {

	// given
	argsMock = []string{"main.exe", "-input", "192.168.137.216:8080 192.168.137.59:8080", "-output", "file", "-method", "grid"}

	// then
	var _, err = ParseArgs(argsMock)
	if err == nil {
		t.Errorf("Error not thrown")
	}

	if strings.Compare(err.Error(), "-output 'file' only valid in combination with -output_filename ") != 0 {
		t.Errorf("Wrong error thrown")
	}
}

func TestShouldFailWithInvalidOutputArgument(t *testing.T) {

	// given
	argsMock = []string{"main.exe", "-input", "192.168.137.216:8080 192.168.137.59:8080", "-output", "XXXX", "-method", "grid"}

	// then
	var _, err = ParseArgs(argsMock)
	if err == nil {
		t.Errorf("Error not thrown")
	}

	//fixme :)
	if strings.Compare(err.Error()[:41], "invalid output argument: -output argument") != 0 {
		t.Errorf("Wrong error thrown")
	}
}

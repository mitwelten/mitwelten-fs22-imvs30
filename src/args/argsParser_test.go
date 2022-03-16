package args

import (
	"mjpeg_multiplexer/src/connection"
	"os"
	"testing"
)

var argsMock = []string{"main.exe", "-input", "192.168.137.216:8080 192.168.137.59:8080", "-output", "file", "-output_filename", "out.jpg", "-method", "grid"}
var expectedInputLocations []string
var expectedOutput connection.Output

func TestStreamCommandShouldNotCrash(t *testing.T) {

	// given
	argsMock = []string{"main.exe", "-input", "192.168.137.216:8080 192.168.137.59:8080", "-output", "stream", "-output_port", "8088", "-method", "grid"}

	// when
	os.Args = argsMock

	// then
	_, _ = ParseArgs()

}

func TestFailCommandShouldNotCrash(t *testing.T) {

	// given
	argsMock = []string{"main.exe", "-input", "192.168.137.216:8080 192.168.137.59:8080", "-output", "file", "-output_filename", "out.jpg", "-method", "grid"}

	// when
	os.Args = argsMock

	// then
	_, _ = ParseArgs()
}

func TestShouldFailWithMissingPort(t *testing.T) {

	// given
	argsMock = []string{"main.exe", "-input", "192.168.137.216:8080 192.168.137.59:8080", "-output", "stream", "-method", "grid"}

	// when
	os.Args = argsMock

	// then
	_, _ = ParseArgs()
}

func TestShouldFailWithMissingFilename(t *testing.T) {

	// given
	argsMock = []string{"main.exe", "-input", "192.168.137.216:8080 192.168.137.59:8080", "-output", "file", "-method", "grid"}

	// when
	os.Args = argsMock

	// then
	_, _ = ParseArgs()
}

func TestShouldFailWithInvalidOutputArgument(t *testing.T) {

	// given
	argsMock = []string{"main.exe", "-input", "192.168.137.216:8080 192.168.137.59:8080", "-output", "XXXX", "-method", "grid"}

	// when
	os.Args = argsMock

	// then
	_, _ = ParseArgs()
}

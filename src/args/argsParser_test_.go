package args

import (
	"errors"
	"fmt"
	"mjpeg_multiplexer/src/connection"
	"mjpeg_multiplexer/src/customErrors"
	"testing"
)

var argsMockGrid = []string{"multiplexer", "grid", "--grid_dimension", "2", "1", "input", "localhost:8080,localhost:8080", "output", "8088", "--log_time"}
var argsMockMotion = []string{"main.exe", "motionDetection", "input", "localhost:8080,localhost:8080", "output", "8088"}
var expectedInputLocations []string
var expectedOutput connection.Output

func TestBasicGridCommand(t *testing.T) {
	// given
	gridComd := argsMockGrid

	fmt.Printf("%v\n", gridComd)

	// then
	_, err := ParseArgs(gridComd)
	if err != nil {
		t.Errorf("Error thrown")
	}
}

func TestShouldFailWithNotFulfillingMinArguments(t *testing.T) {
	// given
	argsMockGrid = []string{"main.exe", "-input", "192.168.137.216:8080 192.168.137.59:8080", "stream", "-mode", "grid", "-grid_dimension", "1 2"}
	var expected *customErrors.ErrArgParserUnfulfilledMinArguments

	// when
	var _, err = ParseArgs(argsMockGrid)

	// then
	if err == nil {
		t.Errorf("Error not thrown")
	}

	println(err.Error())

	if !(errors.As(err, &expected)) {
		t.Errorf("Wrong error thrown")
	}

	//alternative
	//if _, ok := err.(*customErrors.ErrArgParserUnfulfilledMinArguments); ok {
	//	fmt.Printf("is of type: customErrors")
	//} else {
	//	fmt.Println("Using Assert: Error NOT of type customErrors error")
	//	t.Errorf("Wrong error thrown")
	//}
}

func TestShouldFailWithMissingPort(t *testing.T) {
	// given
	argsMockGrid = []string{"main.exe", "-input", "192.168.137.216:8080 192.168.137.59:8080", "-output", "stream", "-mode", "grid", "-grid_dimension", "1 2"}
	var expected *customErrors.ErrArgParserInvalidOutputPort

	// when
	var _, err = ParseArgs(argsMockGrid)

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
	argsMockGrid = []string{"main.exe", "-input", "192.168.137.216:8080 192.168.137.59:8080", "-output", "file", "-mode", "grid", "-grid_dimension", "1 2"}
	var expected *customErrors.ErrArgParserInvalidOutputFilename

	// when
	var _, err = ParseArgs(argsMockGrid)

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
	argsMockGrid = []string{"main.exe", "-input", "192.168.137.216:8080 192.168.137.59:8080", "-output", "XXXX", "-mode", "grid", "-grid_dimension", "1 2"}
	var expected *customErrors.ErrArgParserInvalidArgument

	// when
	var _, err = ParseArgs(argsMockGrid)

	// then
	if err == nil {
		t.Errorf("Error not thrown")
	}

	println(err.Error())

	if !(errors.As(err, &expected)) {
		t.Errorf("Wrong error thrown")
	}
}

func TestShouldFailWithInvalidMode(t *testing.T) {
	// given
	argsMockGrid = []string{"main.exe", "-input", "192.168.137.216:8080 192.168.137.59:8080", "-output", "file", "-output_filename", "name", "-mode", "XXXX", "-grid_dimension", "1 2"}
	var expected *customErrors.ErrArgParserInvalidMode

	// when
	var _, err = ParseArgs(argsMockGrid)

	// then
	if err == nil {
		t.Errorf("Error not thrown")
	}

	println(err.Error())

	if !(errors.As(err, &expected)) {
		t.Errorf("Wrong error thrown")
	}
}

func TestShouldFailWithInvalidGridDimension(t *testing.T) {
	// given
	argsMockGrid = []string{"main.exe", "-input", "192.168.137.216:8080 192.168.137.59:8080", "-output", "file", "-output_filename", "name", "-mode", "grid", "-grid_dimension", "1"}
	var expected *customErrors.ErrArgParserInvalidGridDimension

	// when
	var _, err = ParseArgs(argsMockGrid)

	// then
	if err == nil {
		t.Errorf("Error not thrown")
	}

	println(err.Error())

	if !(errors.As(err, &expected)) {
		t.Errorf("Wrong error thrown")
	}
}

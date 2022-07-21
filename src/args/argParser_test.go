package args

import (
	"strings"
	"testing"
)

const prefixModeMissing = "Mode missing"
const prefixInputMissing = "Input missing"
const prefixOutputMissing = "Output missing"
const prefixValueMissing = "Missing value"
const prefixMalformedArgument = "Malformed argument"
const prefixInvalidOption = "Invalid option"

func expectErrorMessage(t *testing.T, expectedPrefix string, err error) {
	if err == nil {
		t.Fatalf("Expected error {'%v...'} but no error thrown\n", expectedPrefix)
	}
	if !strings.Contains(err.Error(), expectedPrefix) {
		t.Fatalf("Expected error {'%v...'} but got {'%v'}\n", expectedPrefix, err.Error())
	}
}
func expectNoError(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("Unexpected error {%v}\n", err.Error())
	}
}

const separator = " "

func parseArgs(input string) error {
	_, err := ParseArgs(strings.Split(input, separator))
	return err
}

// base string: grid input 192.168.137.76:8080,localhost:8081 output 8088 --log_fps --quality 100

func TestAllOk(t *testing.T) {
	var err error

	err = parseArgs("grid input :8080,:8081,:8082,:8083 output 8088 --grid_dimension 2,2 --width 1000 --height 1000 --show_label --log_fps")
	expectNoError(t, err)

	err = parseArgs("grid input :8080,:8081,:8082,:8083,:8084 output 8088 --grid_dimension 2,3 --width 1000 --height 1000 --ignore_aspect_ratio")
	expectNoError(t, err)

	err = parseArgs("panel input :8080,:8081,:8082 output 8088 --width 1000 --cycle --duration 5")
	expectNoError(t, err)

	err = parseArgs("panel input :8080,:8081,:8082,:8083,:8084 output 8088 --motion --log_fps --debug --show_label")
	expectNoError(t, err)

	err = parseArgs("carousel input :8080,:8081,:8082,:8083,:8084 output 8088 --framerate 20 --log_fps")
	expectNoError(t, err)

	err = parseArgs("carousel input :8080,:8081,:8082,:8083,:8084 output 8088 --quality 80 --motion --show_label --labels=1,2,3,4,5 --label_font_size=50 --log_fps")
	expectNoError(t, err)

	err = parseArgs("grid input 192.168.137.76:8080,localhost:8081 output 8088")
	expectNoError(t, err)

	err = parseArgs("grid input 192.168.137.76:8080,localhost:8081 output 8088 --log_fps")
	expectNoError(t, err)

	err = parseArgs("grid input 192.168.137.76:8080,localhost:8081 output 8088 --log_fps --quality 100")
	expectNoError(t, err)

	err = parseArgs("panel input localhost:8081 output 8088 --log_fps --quality 100")
	expectNoError(t, err)

	err = parseArgs("panel input localhost:8081 output 8088 --log_fps --quality=100")
	expectNoError(t, err)

	err = parseArgs("grid input 192.168.137.76:8080,localhost:8081 output 8088 --log_fps --quality 100 --width 800")
	expectNoError(t, err)

	err = parseArgs("grid input 192.168.137.76:8080,localhost:8081 output 8088 --log_fps --quality 100 --grid_dimension 3,2")
	expectNoError(t, err)

	err = parseArgs("carousel input 192.168.137.76:8080,localhost:8081 output 8088 --log_fps --quality 100")
	expectNoError(t, err)
}

func TestModeMissing(t *testing.T) {
	var err error

	err = parseArgs("_ input 192.168.137.76:8080,localhost:8081 output 8088 --log_fps --quality 100")
	expectErrorMessage(t, prefixModeMissing, err)

	err = parseArgs("grid_ input 192.168.137.76:8080,localhost:8081 output 8088 --log_fps --quality 100")
	expectErrorMessage(t, prefixModeMissing, err)
}

func TestInputMissing(t *testing.T) {
	var err error

	err = parseArgs("grid output 8088 --log_fps --quality 100")
	expectErrorMessage(t, prefixInputMissing, err)

	err = parseArgs("grid input_ 192.168.137.76:8080,localhost:8081 output 8088 --log_fps --quality 100")
	expectErrorMessage(t, prefixInputMissing, err)
}

func TestOutputMissing(t *testing.T) {
	var err error

	err = parseArgs("grid input 192.168.137.76:8080,localhost:8081 --log_fps --quality 100")
	expectErrorMessage(t, prefixOutputMissing, err)

	err = parseArgs("grid input 192.168.137.76:8080,localhost:8081 output_ 8088 --log_fps --quality 100")
	expectErrorMessage(t, prefixOutputMissing, err)
}

func TestValueMissing(t *testing.T) {
	var err error

	err = parseArgs("grid input 192.168.137.76:8080,localhost:8081 output 8088 --log_fps --quality")
	expectErrorMessage(t, prefixValueMissing, err)

	err = parseArgs("grid input 192.168.137.76:8080,localhost:8081 output 8088 --log_fps --quality=")
	expectErrorMessage(t, prefixValueMissing, err)

	err = parseArgs("grid input 192.168.137.76:8080,localhost:8081 output 8088 --quality --log_fps")
	expectErrorMessage(t, prefixValueMissing, err)

	err = parseArgs("grid input 192.168.137.76:8080,localhost:8081 output 8088 --log_fps --quality --grid_dimension 3,2")
	expectErrorMessage(t, prefixValueMissing, err)

}
func TestMalformedArgument(t *testing.T) {
	var err error

	err = parseArgs("grid input 192.168.137.76:8080,localhost:8081 output 8088 --quality 100,")
	expectErrorMessage(t, prefixMalformedArgument, err)

	err = parseArgs("grid input 192.168.137.76:8080,localhost:8081 output 8088 --quality=100,")
	expectErrorMessage(t, prefixMalformedArgument, err)

	err = parseArgs("grid input 192.168.137.76:8080,localhost:8081 output 8088 --quality <invalid>")
	expectErrorMessage(t, prefixMalformedArgument, err)

	err = parseArgs("grid input 192.168.137.76:8080,localhost:8081 output 8088 --log_fps --grid_dimension trash")
	expectErrorMessage(t, prefixMalformedArgument, err)

	err = parseArgs("grid input 192.168.137.76:8080,localhost:8081 output 8088 --log_fps --grid_dimension 3,trash")
	expectErrorMessage(t, prefixMalformedArgument, err)

	err = parseArgs("grid input 192.168.137.76:8080,localhost:8081 output 8088 --log_fps --grid_dimension 3,")
	expectErrorMessage(t, prefixMalformedArgument, err)

	err = parseArgs("grid input 192.168.137.76:8080,localhost:8081 output 8088 --log_fps --grid_dimension 3")
	expectErrorMessage(t, prefixMalformedArgument, err)

}

func TestInvalidOption(t *testing.T) {
	var err error

	err = parseArgs("grid input 192.168.137.76:8080,localhost:8081 output 8088 --log_fps --quality_")
	expectErrorMessage(t, prefixInvalidOption, err)

	err = parseArgs("grid input 192.168.137.76:8080,localhost:8081 output 8088 --log_fps -quality")
	expectErrorMessage(t, prefixInvalidOption, err)

	err = parseArgs("grid input 192.168.137.76:8080,localhost:8081 output 8088 --log_fps ---quality")
	expectErrorMessage(t, prefixInvalidOption, err)

	err = parseArgs("grid input 192.168.137.76:8080,localhost:8081 output 8088 --log_fps quality")
	expectErrorMessage(t, prefixInvalidOption, err)

	err = parseArgs("grid input 192.168.137.76:8080,localhost:8081 output 8088 --log_fps --grid")
	expectErrorMessage(t, prefixInvalidOption, err)
}

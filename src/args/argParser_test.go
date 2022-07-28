package args

import (
	"mjpeg_multiplexer/src/utils"
	"strings"
	"testing"
)

const prefixModeMissing = "Mode missing"
const prefixInputMissing = "Input missing"
const prefixOutputMissing = "Output missing"
const prefixValueMissing = "Missing value"
const prefixMalformedArgument = "Malformed argument"
const prefixInvalidOption = "Invalid option"
const prefixMotionWithGrid = "Option '--motion' only available for the modes 'panel' or 'carousel'."
const prefixGridDimensionWithoutGrid = "Option '--grid_dimension=ROWS,COLUMNS' only available for the mode 'grid'."
const prefixPanelCycleWithoutPanel = "Option '--panel_cycle' only available for the mode 'panel'."

const separator = " "

func parseArgs(input string) error {
	_, err := ParseArgs(strings.Split(input, separator))
	return err
}

// base string: grid input 192.168.137.76:8080,localhost:8081 output 8088 --log_fps --quality 100

func TestAllOk(t *testing.T) {
	var err error

	err = parseArgs("grid input :8080,:8081,:8082,:8083 output 8088 --grid_dimension 2,2 --width 1000 --height 1000 --show_label --log_fps")
	utils.ExpectNoError(t, err)

	err = parseArgs("grid input :8080,:8081,:8082,:8083,:8084 output 8088 --grid_dimension 2,3 --width 1000 --height 1000 --ignore_aspect_ratio")
	utils.ExpectNoError(t, err)

	err = parseArgs("panel input :8080,:8081,:8082 output 8088 --width 1000 --panel_cycle --duration 5")
	utils.ExpectNoError(t, err)

	err = parseArgs("panel input :8080,:8081,:8082,:8083,:8084 output 8088 --motion --log_fps --debug --show_label")
	utils.ExpectNoError(t, err)

	err = parseArgs("carousel input :8080,:8081,:8082,:8083,:8084 output 8088 --framerate 20 --log_fps")
	utils.ExpectNoError(t, err)

	err = parseArgs("carousel input :8080,:8081,:8082,:8083,:8084 output 8088 --quality 80 --motion --show_label --labels=1,2,3,4,5 --label_font_size=50 --log_fps")
	utils.ExpectNoError(t, err)

	err = parseArgs("grid input 192.168.137.76:8080,localhost:8081 output 8088")
	utils.ExpectNoError(t, err)

	err = parseArgs("grid input 192.168.137.76:8080,localhost:8081 output 8088 --log_fps")
	utils.ExpectNoError(t, err)

	err = parseArgs("grid input 192.168.137.76:8080,localhost:8081 output 8088 --log_fps --quality 100")
	utils.ExpectNoError(t, err)

	err = parseArgs("panel input localhost:8081 output 8088 --log_fps --quality 100")
	utils.ExpectNoError(t, err)

	err = parseArgs("panel input localhost:8081 output 8088 --log_fps --quality=100")
	utils.ExpectNoError(t, err)

	err = parseArgs("grid input 192.168.137.76:8080,localhost:8081 output 8088 --log_fps --quality 100 --width 800")
	utils.ExpectNoError(t, err)

	err = parseArgs("grid input 192.168.137.76:8080,localhost:8081 output 8088 --log_fps --quality 100 --grid_dimension 3,2")
	utils.ExpectNoError(t, err)

	err = parseArgs("carousel input 192.168.137.76:8080,localhost:8081 output 8088 --log_fps --quality 100")
	utils.ExpectNoError(t, err)
}

func TestModeMissing(t *testing.T) {
	var err error

	err = parseArgs("_ input 192.168.137.76:8080,localhost:8081 output 8088 --log_fps --quality 100")
	utils.ExpectErrorMessage(t, prefixModeMissing, err)

	err = parseArgs("grid_ input 192.168.137.76:8080,localhost:8081 output 8088 --log_fps --quality 100")
	utils.ExpectErrorMessage(t, prefixModeMissing, err)
}

func TestInputMissing(t *testing.T) {
	var err error

	err = parseArgs("grid output 8088 --log_fps --quality 100")
	utils.ExpectErrorMessage(t, prefixInputMissing, err)

	err = parseArgs("grid input_ 192.168.137.76:8080,localhost:8081 output 8088 --log_fps --quality 100")
	utils.ExpectErrorMessage(t, prefixInputMissing, err)

	err = parseArgs("grid input output 8088 --log_fps --quality 100")
	utils.ExpectErrorMessage(t, prefixInputMissing, err)
}

func TestOutputMissing(t *testing.T) {
	var err error

	err = parseArgs("grid input 192.168.137.76:8080,localhost:8081 --log_fps --quality 100")
	utils.ExpectErrorMessage(t, prefixOutputMissing, err)

	err = parseArgs("grid input 192.168.137.76:8080,localhost:8081 output_ 8088 --log_fps --quality 100")
	utils.ExpectErrorMessage(t, prefixOutputMissing, err)

	err = parseArgs("grid input 192.168.137.76:8080,localhost:8081 output --log_fps --quality 100")
	utils.ExpectErrorMessage(t, prefixOutputMissing, err)
}

func TestValueMissing(t *testing.T) {
	var err error

	err = parseArgs("grid input 192.168.137.76:8080,localhost:8081 output 8088 --log_fps --quality")
	utils.ExpectErrorMessage(t, prefixValueMissing, err)

	err = parseArgs("grid input 192.168.137.76:8080,localhost:8081 output 8088 --log_fps --quality=")
	utils.ExpectErrorMessage(t, prefixValueMissing, err)

	err = parseArgs("grid input 192.168.137.76:8080,localhost:8081 output 8088 --quality --log_fps")
	utils.ExpectErrorMessage(t, prefixValueMissing, err)

	err = parseArgs("grid input 192.168.137.76:8080,localhost:8081 output 8088 --log_fps --quality --grid_dimension 3,2")
	utils.ExpectErrorMessage(t, prefixValueMissing, err)

}
func TestMalformedArgument(t *testing.T) {
	var err error

	err = parseArgs("grid input 192.168.137.76:8080,localhost:8081 output 8088 --quality 100,")
	utils.ExpectErrorMessage(t, prefixMalformedArgument, err)

	err = parseArgs("grid input 192.168.137.76:8080,localhost:8081 output 8088 --quality=100,")
	utils.ExpectErrorMessage(t, prefixMalformedArgument, err)

	err = parseArgs("grid input 192.168.137.76:8080,localhost:8081 output 8088 --quality <invalid>")
	utils.ExpectErrorMessage(t, prefixMalformedArgument, err)

	err = parseArgs("grid input 192.168.137.76:8080,localhost:8081 output 8088 --log_fps --grid_dimension trash")
	utils.ExpectErrorMessage(t, prefixMalformedArgument, err)

	err = parseArgs("grid input 192.168.137.76:8080,localhost:8081 output 8088 --log_fps --grid_dimension 3,trash")
	utils.ExpectErrorMessage(t, prefixMalformedArgument, err)

	err = parseArgs("grid input 192.168.137.76:8080,localhost:8081 output 8088 --log_fps --grid_dimension 3,")
	utils.ExpectErrorMessage(t, prefixMalformedArgument, err)

	err = parseArgs("grid input 192.168.137.76:8080,localhost:8081 output 8088 --log_fps --grid_dimension 3")
	utils.ExpectErrorMessage(t, prefixMalformedArgument, err)

}

func TestInvalidOption(t *testing.T) {
	var err error

	err = parseArgs("grid input 192.168.137.76:8080,localhost:8081 output 8088 --log_fps --quality_")
	utils.ExpectErrorMessage(t, prefixInvalidOption, err)

	err = parseArgs("grid input 192.168.137.76:8080,localhost:8081 output 8088 --log_fps -quality")
	utils.ExpectErrorMessage(t, prefixInvalidOption, err)

	err = parseArgs("grid input 192.168.137.76:8080,localhost:8081 output 8088 --log_fps ---quality")
	utils.ExpectErrorMessage(t, prefixInvalidOption, err)

	err = parseArgs("grid input 192.168.137.76:8080,localhost:8081 output 8088 --log_fps quality")
	utils.ExpectErrorMessage(t, prefixInvalidOption, err)

	err = parseArgs("grid input 192.168.137.76:8080,localhost:8081 output 8088 --log_fps --grid")
	utils.ExpectErrorMessage(t, prefixInvalidOption, err)
}

func TestMotionWithGrid(t *testing.T) {
	var err error

	err = parseArgs("panel input :8080,:8081,:8082,:8083 output 8088 --width 1000 --height 1000 --show_label --log_fps --motion")
	utils.ExpectNoError(t, err)

	err = parseArgs("grid input :8080,:8081,:8082,:8083 output 8088 --width 1000 --height 1000 --show_label --log_fps")
	utils.ExpectNoError(t, err)

	err = parseArgs("grid input :8080,:8081,:8082,:8083 output 8088 --width 1000 --height 1000 --show_label --log_fps --motion")
	utils.ExpectErrorMessage(t, prefixMotionWithGrid, err)
}

func TestGridDimensionWithoutGrid(t *testing.T) {
	var err error

	err = parseArgs("grid input :8080,:8081,:8082,:8083 output 8088 --width 1000 --height 1000 --show_label --log_fps --grid_dimension 3,3")
	utils.ExpectNoError(t, err)

	err = parseArgs("panel input :8080,:8081,:8082,:8083 output 8088 --width 1000 --height 1000 --show_label --log_fps --grid_dimension 3,3")
	utils.ExpectErrorMessage(t, prefixGridDimensionWithoutGrid, err)

	err = parseArgs("carousel input :8080,:8081,:8082,:8083 output 8088 --width 1000 --height 1000 --show_label --log_fps --grid_dimension 3,3")
	utils.ExpectErrorMessage(t, prefixGridDimensionWithoutGrid, err)
}

func TestPanelCycleWithoutPanel(t *testing.T) {
	var err error

	err = parseArgs("panel input :8080,:8081,:8082,:8083 output 8088 --width 1000 --height 1000 --show_label --log_fps --panel_cycle")
	utils.ExpectNoError(t, err)

	err = parseArgs("grid input :8080,:8081,:8082,:8083 output 8088 --width 1000 --height 1000 --show_label --log_fps --panel_cycle")
	utils.ExpectErrorMessage(t, prefixPanelCycleWithoutPanel, err)

	err = parseArgs("carousel input :8080,:8081,:8082,:8083 output 8088 --width 1000 --height 1000 --show_label --log_fps --panel_cycle")
	utils.ExpectErrorMessage(t, prefixPanelCycleWithoutPanel, err)
}

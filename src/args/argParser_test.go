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
const prefixGridTooSmall = "Invalid configuration '--grid_dimension"
const prefixInvalidOption = "Invalid option"
const levenshteinString = "Did you mean"

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

	err = parseArgs("panel input :8080,:8081,:8082,:8083,:8084 output 8088 --activity --log_fps --debug --show_label")
	utils.ExpectNoError(t, err)

	err = parseArgs("carousel input :8080,:8081,:8082,:8083,:8084 output 8088 --framerate 20 --log_fps")
	utils.ExpectNoError(t, err)

	err = parseArgs("carousel input :8080,:8081,:8082,:8083,:8084 output 8088 --quality 80 --activity --show_label --labels=1,2,3,4,5 --label_font_size=50 --log_fps")
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
	utils.ExpectErrorContains(t, prefixModeMissing, err)

	err = parseArgs("grid_ input 192.168.137.76:8080,localhost:8081 output 8088 --log_fps --quality 100")
	utils.ExpectErrorContains(t, prefixModeMissing, err)
}

func TestInputMissing(t *testing.T) {
	var err error

	err = parseArgs("grid output 8088 --log_fps --quality 100")
	utils.ExpectErrorContains(t, prefixInputMissing, err)

	err = parseArgs("grid input_ 192.168.137.76:8080,localhost:8081 output 8088 --log_fps --quality 100")
	utils.ExpectErrorContains(t, prefixInputMissing, err)

	err = parseArgs("grid input output 8088 --log_fps --quality 100")
	utils.ExpectErrorContains(t, prefixInputMissing, err)
}

func TestOutputMissing(t *testing.T) {
	var err error

	err = parseArgs("grid input 192.168.137.76:8080,localhost:8081 --log_fps --quality 100")
	utils.ExpectErrorContains(t, prefixOutputMissing, err)

	err = parseArgs("grid input 192.168.137.76:8080,localhost:8081 output_ 8088 --log_fps --quality 100")
	utils.ExpectErrorContains(t, prefixOutputMissing, err)

	err = parseArgs("grid input 192.168.137.76:8080,localhost:8081 output --log_fps --quality 100")
	utils.ExpectErrorContains(t, prefixOutputMissing, err)
}

func TestValueMissing(t *testing.T) {
	var err error

	err = parseArgs("grid input 192.168.137.76:8080,localhost:8081 output 8088 --log_fps --quality")
	utils.ExpectErrorContains(t, prefixValueMissing, err)

	err = parseArgs("grid input 192.168.137.76:8080,localhost:8081 output 8088 --log_fps --quality=")
	utils.ExpectErrorContains(t, prefixValueMissing, err)

	err = parseArgs("grid input 192.168.137.76:8080,localhost:8081 output 8088 --quality --log_fps")
	utils.ExpectErrorContains(t, prefixValueMissing, err)

	err = parseArgs("grid input 192.168.137.76:8080,localhost:8081 output 8088 --log_fps --quality --grid_dimension 3,2")
	utils.ExpectErrorContains(t, prefixValueMissing, err)

}
func TestMalformedArgument(t *testing.T) {
	var err error

	err = parseArgs("grid input 192.168.137.76:8080,localhost:8081 output 8088 --quality 100,")
	utils.ExpectErrorContains(t, prefixMalformedArgument, err)

	err = parseArgs("grid input 192.168.137.76:8080,localhost:8081 output 8088 --quality=100,")
	utils.ExpectErrorContains(t, prefixMalformedArgument, err)

	err = parseArgs("grid input 192.168.137.76:8080,localhost:8081 output 8088 --quality <invalid>")
	utils.ExpectErrorContains(t, prefixMalformedArgument, err)

	err = parseArgs("grid input 192.168.137.76:8080,localhost:8081 output 8088 --log_fps --grid_dimension trash")
	utils.ExpectErrorContains(t, prefixMalformedArgument, err)

	err = parseArgs("grid input 192.168.137.76:8080,localhost:8081 output 8088 --log_fps --grid_dimension 3,trash")
	utils.ExpectErrorContains(t, prefixMalformedArgument, err)

	err = parseArgs("grid input 192.168.137.76:8080,localhost:8081 output 8088 --log_fps --grid_dimension 3,")
	utils.ExpectErrorContains(t, prefixMalformedArgument, err)

	err = parseArgs("grid input 192.168.137.76:8080,localhost:8081 output 8088 --log_fps --grid_dimension 3")
	utils.ExpectErrorContains(t, prefixMalformedArgument, err)

}

func TestInvalidOption(t *testing.T) {
	var err error

	err = parseArgs("grid input 192.168.137.76:8080,localhost:8081 output 8088 --log_fps --quality_")
	utils.ExpectErrorContains(t, prefixInvalidOption, err)

	err = parseArgs("grid input 192.168.137.76:8080,localhost:8081 output 8088 --log_fps -quality")
	utils.ExpectErrorContains(t, prefixInvalidOption, err)

	err = parseArgs("grid input 192.168.137.76:8080,localhost:8081 output 8088 --log_fps ---quality")
	utils.ExpectErrorContains(t, prefixInvalidOption, err)

	err = parseArgs("grid input 192.168.137.76:8080,localhost:8081 output 8088 --log_fps quality")
	utils.ExpectErrorContains(t, prefixInvalidOption, err)

	err = parseArgs("grid input 192.168.137.76:8080,localhost:8081 output 8088 --log_fps --grid")
	utils.ExpectErrorContains(t, prefixInvalidOption, err)
}

func TestGridTooSmall(t *testing.T) {
	var err error

	err = parseArgs("grid input 192.168.137.76:8080,localhost:8081 output 8088 --log_fps --grid_dimension 0,0")
	utils.ExpectErrorContains(t, prefixGridTooSmall, err)

	err = parseArgs("grid input 192.168.137.76:8080,localhost:8081 output 8088 --log_fps --grid_dimension 1,0")
	utils.ExpectErrorContains(t, prefixGridTooSmall, err)

	err = parseArgs("grid input 192.168.137.76:8080,localhost:8081 output 8088 --log_fps --grid_dimension 1,1")
	utils.ExpectErrorContains(t, prefixGridTooSmall, err)

	err = parseArgs("grid input 192.168.137.76:8080,localhost:8081,192.168.137.76:8080,localhost:8081 output 8088 --log_fps --grid_dimension 2,1")
	utils.ExpectErrorContains(t, prefixGridTooSmall, err)

	err = parseArgs("grid input 192.168.137.76:8080,localhost:8081,192.168.137.76:8080,localhost:8081 output 8088 --log_fps --grid_dimension 2,2")
	utils.ExpectNoError(t, err)

	err = parseArgs("grid input 192.168.137.76:8080 output 8088 --log_fps --grid_dimension 1,1")
	utils.ExpectNoError(t, err)

	err = parseArgs("grid input 192.168.137.76:8080 output 8088 --log_fps")
	utils.ExpectNoError(t, err)
}

func TestActivityWithGrid(t *testing.T) {
	var err error

	err = parseArgs("panel input :8080,:8081,:8082,:8083 output 8088 --width 1000 --height 1000 --show_label --log_fps --activity")
	utils.ExpectNoError(t, err)

	err = parseArgs("grid input :8080,:8081,:8082,:8083 output 8088 --width 1000 --height 1000 --show_label --log_fps")
	utils.ExpectNoError(t, err)

	err = parseArgs("grid input :8080,:8081,:8082,:8083 output 8088 --width 1000 --height 1000 --show_label --log_fps --activity")
	utils.ExpectErrorContains(t, errActivityWithGrid, err)
}

func TestGridDimensionWithoutGrid(t *testing.T) {
	var err error

	err = parseArgs("grid input :8080,:8081,:8082,:8083 output 8088 --width 1000 --height 1000 --show_label --log_fps --grid_dimension 3,3")
	utils.ExpectNoError(t, err)

	err = parseArgs("panel input :8080,:8081,:8082,:8083 output 8088 --width 1000 --height 1000 --show_label --log_fps --grid_dimension 3,3")
	utils.ExpectErrorContains(t, errGridDimensionWithoutGrid, err)

	err = parseArgs("carousel input :8080,:8081,:8082,:8083 output 8088 --width 1000 --height 1000 --show_label --log_fps --grid_dimension 3,3")
	utils.ExpectErrorContains(t, errGridDimensionWithoutGrid, err)
}

func TestPanelCycleWithoutPanel(t *testing.T) {
	var err error

	err = parseArgs("panel input :8080,:8081,:8082,:8083 output 8088 --width 1000 --height 1000 --show_label --log_fps --panel_cycle")
	utils.ExpectNoError(t, err)

	err = parseArgs("grid input :8080,:8081,:8082,:8083 output 8088 --width 1000 --height 1000 --show_label --log_fps --panel_cycle")
	utils.ExpectErrorContains(t, errPanelCycleWithoutPanel, err)

	err = parseArgs("carousel input :8080,:8081,:8082,:8083 output 8088 --width 1000 --height 1000 --show_label --log_fps --panel_cycle")
	utils.ExpectErrorContains(t, errPanelCycleWithoutPanel, err)
}
func TestDurationWithGrid(t *testing.T) {
	var err error

	err = parseArgs("grid input :8080,:8081,:8082,:8083 output 8088 --width 1000 --height 1000 --show_label --log_fps --duration 10")
	utils.ExpectErrorContains(t, errDurationWithGrid, err)

	err = parseArgs("panel input :8080,:8081,:8082,:8083 output 8088 --width 1000 --height 1000 --show_label --log_fps --duration 10")
	utils.ExpectNoError(t, err)

	err = parseArgs("carousel input :8080,:8081,:8082,:8083 output 8088 --width 1000 --height 1000 --show_label --log_fps --duration 10")
	utils.ExpectNoError(t, err)
}
func TestBorderWithCarousel(t *testing.T) {
	var err error

	err = parseArgs("grid input :8080,:8081,:8082,:8083 output 8088 --width 1000 --height 1000 --show_label --log_fps --show_border")
	utils.ExpectNoError(t, err)

	err = parseArgs("panel input :8080,:8081,:8082,:8083 output 8088 --width 1000 --height 1000 --show_label --log_fps --show_border")
	utils.ExpectNoError(t, err)

	err = parseArgs("carousel input :8080,:8081,:8082,:8083 output 8088 --width 1000 --height 1000 --show_label --log_fps --show_border")
	utils.ExpectErrorContains(t, errBorderWithCarousel, err)
}

func TestLabelWithoutShowLabel(t *testing.T) {
	var err error

	err = parseArgs("grid input :8080,:8081,:8082,:8083 output 8088 --width 1000 --height 1000 --show_label --labels=l1,l2,l3,l4")
	utils.ExpectNoError(t, err)

	err = parseArgs("grid input :8080,:8081,:8082,:8083 output 8088 --width 1000 --height 1000 --labels=l1,l2,l3,l4")
	utils.ExpectErrorContains(t, errLabelWithoutShowLabel, err)

}

func TestLabelFontSizeWithoutShowLabel(t *testing.T) {
	var err error

	err = parseArgs("grid input :8080,:8081,:8082,:8083 output 8088 --width 1000 --height 1000 --show_label --label_font_size 100")
	utils.ExpectNoError(t, err)

	err = parseArgs("grid input :8080,:8081,:8082,:8083 output 8088 --width 1000 --height 1000 --label_font_size 100")
	utils.ExpectErrorContains(t, errLabelFontSizeWithoutShowLabel, err)
}

func TestLevenshtein(t *testing.T) {
	var err error

	err = parseArgs("grid input :8080,:8081,:8082,:8083 output 8088 --width 1000 --height 1000 --show_label --label_font_size 100")
	utils.ExpectNoError(t, err)

	err = parseArgs("grid input :8080,:8081,:8082,:8083 output 8088 --wdth 1000 --height 1000 --label_font_size 100")
	utils.ExpectErrorContains(t, levenshteinString, err)

	err = parseArgs("grid input :8080,:8081,:8082,:8083 output 8088 --width 1000 --heihgt 1000 --label_font_size 100")
	utils.ExpectErrorContains(t, levenshteinString, err)

	err = parseArgs("grid input :8080,:8081,:8082,:8083 output 8088 --width 1000 --height 1000 --labels_font_size 100")
	utils.ExpectErrorContains(t, levenshteinString, err)

	err = parseArgs("grid input :8080,:8081,:8082,:8083 output 8088 --width 1000 --height 1000 --show_labels")
	utils.ExpectErrorContains(t, levenshteinString, err)
}

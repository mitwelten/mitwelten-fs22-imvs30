package args

import (
	"fmt"
	"github.com/docopt/docopt.go"
	"mjpeg_multiplexer/src/aggregator"
	"mjpeg_multiplexer/src/connection"
	"mjpeg_multiplexer/src/customErrors"
	"mjpeg_multiplexer/src/global"
	"mjpeg_multiplexer/src/multiplexer"
	"strings"
)

const (
	InputLocationSeparator = ","
)

var (
	usage = `Usage:
  multiplexer (motion) (input) URL (output) PORT [options]
  multiplexer (grid) (--grid_dimension GRID_Y GRID_X) (input) URL (output) PORT [options] 
  multiplexer -h | --help
  multiplexer --version

Options:
  -h --help              Shows this screen.
  --input_framerate=n    input framerate in fps [default: -1]
  --output_framerate=n   output framerate in fps[default: -1]
  --output_max_width=n   output width in pixel [default: -1]
  --output_max_height=n  output height in pixel [default: -1]  
  --use_auth   	         Use Authentication
  --log_time	         Log Time verbose
  --verbose              Shows details.
  --version              Shows version.`
)

// parseInput parses input URLS derived from command line arguments
func parseInputUrls(config multiplexer.MultiplexerConfig, inputStr string) multiplexer.MultiplexerConfig {
	inputUrls := strings.Split(inputStr, InputLocationSeparator)
	var conns []connection.Input
	for _, s := range inputUrls {
		conns = append(conns, connection.NewInputHTTP(s))
	}
	config.InputLocations = conns

	return config
}

// ParseArgs parses all arguments derived from command line
func ParseArgs(args []string) (config multiplexer.MultiplexerConfig, err error) {

	// parse all arguments and save them to vars
	arguments, _ := docopt.ParseDoc(usage)
	fmt.Println(arguments)

	// mandatory
	input, _ := arguments.String("URL") // input
	port, _ := arguments.String("PORT") // output

	// mode
	grid, _ := arguments.Bool("grid")
	//_, _ := arguments.Bool("motion") // could be used later in a switch statement
	gridX, _ := arguments.Int("GRID_X")
	gridY, _ := arguments.Int("GRID_Y")

	// options
	inputFramerate, _ := arguments.Float64("--input_framerate")
	outputFramerate, _ := arguments.Float64("--output_framerate")
	outputWidth, _ := arguments.Int("--output_max_height")
	outputHeight, _ := arguments.Int("--output_width")
	useAuth, _ := arguments.Bool("--use_auth")
	logTime, _ := arguments.Bool("--log_time")

	// global config

	// inputURL parsing
	config = parseInputUrls(config, input)

	// output stream
	config.Output, err = connection.NewOutputHTTP(port)
	if err != nil {
		return multiplexer.MultiplexerConfig{}, &customErrors.ErrHttpOpenOutputSocket{}
	}

	// disabled: output file
	// config.Output = connection.NewOutputFile(outputFileNamePtr)

	if outputWidth != -1 && outputHeight != -1 {
		global.Config.MaxWidth = outputWidth
		global.Config.MaxHeight = outputHeight
	}

	// mode
	if grid {
		config.Aggregator = &aggregator.AggregatorGrid{Row: gridY, Col: gridX}
	} else {
		config.Aggregator = &aggregator.AggregatorChange{}
	}

	//	return multiplexer.MultiplexerConfig{}, &customErrors.ErrArgParserInvalidMode{Argument: *modePtr}

	// credentials
	if useAuth {
		config.UseAuth = true
	} else {
		config.UseAuth = false
	}

	// logtime
	if logTime {
		global.Config.LogTime = true
	}

	// InputRates
	global.Config.InputFramerate = inputFramerate
	global.Config.OutputFramerate = outputFramerate

	// non error case, return nil
	return config, nil
}

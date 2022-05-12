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
	"time"
)

const (
	InputLocationSeparator = ","
)

var (
	usage = `Usage:
  multiplexer (motion) (input) URL (output) PORT [options]
  multiplexer (grid) (grid_dimension GRID_Y GRID_X) (input) URL (output) PORT [options] 
  multiplexer -h | --help
  multiplexer --version

Options:
  -h --help             Shows this screen.
  --input_framerate=n   input framerate in fps [default: -1]
  --output_framerate=n  output framerate in fps[default: -1]
  --use_auth   	        Use Authentication
  --log_time	        Log Time verbose
  --verbose             Shows details.
  --version             Shows version.`
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
	grid_x, _ := arguments.Int("GRID_X")
	grid_y, _ := arguments.Int("GRID_Y")

	// options
	input_framerate, _ := arguments.Int("--input_framerate")
	output_framerate, _ := arguments.Int("--output_framerate")
	output_width, _ := arguments.Int("--output_height")
	output_height, _ := arguments.Int("--output_width")
	use_auth, _ := arguments.Bool("--use_auth")
	log_time, _ := arguments.Bool("--log_time")

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

	if output_width != -1 && output_height != -1 {
		global.Config.MaxWidth = output_width
		global.Config.MaxHeight = output_height
	}

	// mode
	if grid {
		config.Aggregator = &aggregator.AggregatorGrid{Row: grid_y, Col: grid_x}
	} else {
		config.Aggregator = &aggregator.AggregatorChange{}
	}

	//	return multiplexer.MultiplexerConfig{}, &customErrors.ErrArgParserInvalidMode{Argument: *modePtr}

	// credentials
	if use_auth {
		config.UseAuth = true
	} else {
		config.UseAuth = false
	}

	// logtime
	if log_time {
		global.Config.LogTime = true
	}

	// minInputDelay
	global.Config.MinimumInputDelay = time.Duration(input_framerate) * time.Millisecond

	// non error case, return nil
	return config, nil
}

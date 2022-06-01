package args

import (
	"fmt"
	"github.com/docopt/docopt.go"
	"mjpeg_multiplexer/src/aggregator"
	"mjpeg_multiplexer/src/connection"
	"mjpeg_multiplexer/src/customErrors"
	"mjpeg_multiplexer/src/global"
	"mjpeg_multiplexer/src/multiplexer"
	"os"
	"strings"
	"time"
)

const (
	InputLocationSeparator = ","
)

var (
	usage = `Usage:
  multiplexer (grid) (--dimension GRID_Y GRID_X) (input) URL (output) PORT [options] 
  multiplexer (carousel) (--duration DURATION) (input) URL (output) PORT [options]
  multiplexer (panel) [--cycle] [--duration DURATION] (input) URL (output) PORT [options]  
  multiplexer -h | --help
  multiplexer --version

Options:
  -h --help              Shows this screen.
  --duration=n           frame duration [default: 10]
  --input_framerate=n    input framerate in fps [default: -1]
  --output_framerate=n   output framerate in fps[default: -1]
  --output_max_width=n   output width in pixel [default: -1]
  --output_max_height=n  output height in pixel [default: -1]
  --output_quality=n     output jpeg quality in percentage [default: 100]
  --border=n             number of black pixels between each image [default: 0]
  --use_auth             Use Authentication
  --label                Show label for input streams
  --log_time             Log Time verbose
  --motion               Enables motion detection to focus the most active frame on selected mode
  --verbose              Shows details. 
  --version              Shows version.`
)

// parseInput parses input URLS derived from command line arguments
func parseInputUrls(config multiplexer.MultiplexerConfig, inputStr string) multiplexer.MultiplexerConfig {
	inputUrls := strings.Split(inputStr, InputLocationSeparator)
	config.InputLocations = []connection.Input{}
	for i, url := range inputUrls {
		global.Config.InputConfigs = append(global.Config.InputConfigs, global.InputConfig{Url: url})
		config.InputLocations = append(config.InputLocations, connection.NewInputHTTP(&global.Config.InputConfigs[i], url))
	}

	return config
}

var printUsage = func(err error, usage_ string) {
	fmt.Println(usage)
	os.Exit(0)
}

// ParseArgs parses all arguments derived from command line
func ParseArgs(args []string) (config multiplexer.MultiplexerConfig, err error) {
	// init custom handler to print full usage on error
	parser := &docopt.Parser{
		HelpHandler:  printUsage,
		OptionsFirst: false,
	}
	arguments, err := parser.ParseArgs(usage, nil, "")

	// parse all arguments and save them to vars
	//arguments, _ := docopt.ParseDoc(usage)
	fmt.Println(arguments)

	// mandatory
	input, _ := arguments.String("URL") // input
	port, _ := arguments.String("PORT") // output

	// mode
	grid, _ := arguments.Bool("grid")
	panel, _ := arguments.Bool("panel")
	carousel, _ := arguments.Bool("carousel")
	// mode options
	gridX, _ := arguments.Int("GRID_X")
	gridY, _ := arguments.Int("GRID_Y")
	duration, _ := arguments.Int("--duration") // carousel or panel-cycle duration in seconds

	panelCycle, _ := arguments.Bool("--cycle") // panel cycle, default false
	// options
	inputFramerate, _ := arguments.Float64("--input_framerate")
	outputFramerate, _ := arguments.Float64("--output_framerate")
	outputMaxWidth, _ := arguments.Int("--output_max_width")
	outputMaxHeight, _ := arguments.Int("--output_max_height")
	outputQuality, _ := arguments.Int("--output_quality")
	outputMargin, _ := arguments.Int("--border")
	useAuth, _ := arguments.Bool("--use_auth")
	logTime, _ := arguments.Bool("--log_time")
	showInputLabel, _ := arguments.Bool("--label")
	useMotion, _ := arguments.Bool("--motion")

	// global config

	// inputURL parsing
	config = parseInputUrls(config, input)

	// mode
	if grid {
		config.Aggregator = &aggregator.AggregatorGrid{Row: gridY, Col: gridX}
	} else if carousel {
		config.Aggregator = &aggregator.AggregatorCarousel{Duration: time.Duration(duration) * time.Second}
	} else if panel {
		config.Aggregator = &aggregator.AggregatorPanel{Duration: time.Duration(duration) * time.Second, CycleFrames: panelCycle}
	} else {
		return multiplexer.MultiplexerConfig{}, &customErrors.ErrArgParserInvalidArgument{}
	}

	// output stream
	config.Output, err = connection.NewOutputHTTP(port, config.Aggregator)
	if err != nil {
		return multiplexer.MultiplexerConfig{}, &customErrors.ErrHttpOpenOutputSocket{}
	}

	// disabled: output file
	// config.Output = connection.NewOutputFile(outputFileNamePtr)

	if outputMaxWidth != -1 || outputMaxHeight != -1 {
		global.Config.MaxWidth = outputMaxWidth
		global.Config.MaxHeight = outputMaxHeight
	}

	//	return multiplexer.MultiplexerConfig{}, &customErrors.ErrArgParserInvalidMode{Argument: *modePtr}

	// credentials
	global.Config.UseAuth = useAuth

	// logtime
	global.Config.LogTime = logTime

	// InputRates
	global.Config.InputFramerate = inputFramerate
	global.Config.OutputFramerate = outputFramerate

	// quality
	global.Config.EncodeQuality = outputQuality

	// border
	global.Config.Border = outputMargin

	// label
	global.Config.ShowInputLabel = showInputLabel

	//motion
	global.Config.UseMotion = useMotion

	// non error case, return nil
	return config, nil
}

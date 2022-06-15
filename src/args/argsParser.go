package args

import (
	"fmt"
	"github.com/docopt/docopt.go"
	"log"
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
	ArgumentSeparator = ","
)

var (
	usage = `Usage:
  multiplexer (grid) (--dimension GRID_Y GRID_X) (input) URL (output) PORT [options] 
  multiplexer (carousel) (--duration DURATION) (input) URL (output) PORT [options]
  multiplexer (panel) [--cycle] [--duration DURATION] (input) URL (output) PORT [options]  
  multiplexer -h | --help
  multiplexer --version

Options:
  -h --help                        Shows this screen.
  --duration=n                     frame duration [default: 10]
  --input_framerate=n              input framerate in fps [default: -1]
  --output_framerate=n             output framerate in fps[default: -1]
  --width=n                        output width in pixel [default: -1]
  --height=n                       output height in pixel [default: -1]
  --ignore_aspect_ratio            todo
  --output_quality=n               output jpeg quality in percentage [default: 100]
  --border                         number of black pixels between each image
  --use_auth                       Use Authentication
  --show_label                     Show label for input streams
  --labels=n                       comma separated list of names to show instead of the camera input url
  --label_font_size=n              input label font size in pixel [default: 32]
  --log_time                       Log Time verbose
  --motion                         Enables motion detection to focus the most active frame on selected mode
  --verbose                        Shows details. 
  --version                        Shows version.`
)

// parseInput parses input URLS derived from command line arguments
func parseInputUrls(config multiplexer.MultiplexerConfig, inputStr string) multiplexer.MultiplexerConfig {
	arr := strings.Split(inputStr, ArgumentSeparator)
	config.InputLocations = []connection.Input{}
	for i, url := range arr {
		global.Config.InputConfigs = append(global.Config.InputConfigs, global.InputConfig{Url: url, Label: url})
		config.InputLocations = append(config.InputLocations, connection.NewInputHTTP(&global.Config.InputConfigs[i], url))
	}

	return config
}

// todo: evtl. trim
func parseSeparatedString(inputStr string) {
	arr := strings.Split(inputStr, ArgumentSeparator)
	if len(global.Config.InputConfigs) != len(arr) {
		log.Fatalf("%v input location present, but %v labels found\n", len(global.Config.InputConfigs), len(arr))
	}
	for i, label := range arr {
		global.Config.InputConfigs[i].Label = label
	}
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
	width, _ := arguments.Int("--width")
	height, _ := arguments.Int("--height")
	ignoreAspectRatio, _ := arguments.Bool("--ignore_aspect_ratio")
	outputQuality, _ := arguments.Int("--output_quality")
	useBorder, _ := arguments.Bool("--border")
	useAuth, _ := arguments.Bool("--use_auth")
	logTime, _ := arguments.Bool("--log_time")
	showInputLabel, _ := arguments.Bool("--show_label")
	useMotion, _ := arguments.Bool("--motion")
	inputLabels, _ := arguments.String("--labels") // input inputLabels
	inputLabelFontSize, _ := arguments.Int("--label_font_size")

	// global config

	// inputURL and label parsing
	config = parseInputUrls(config, input)
	if len(inputLabels) != 0 {
		parseSeparatedString(inputLabels)
	}

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
	config.Output = connection.NewOutputHTTP(port, config.Aggregator)

	// disabled: output file
	// config.Output = connection.NewOutputFile(outputFileNamePtr)

	global.Config.Width = width
	global.Config.Height = height
	global.Config.IgnoreAspectRatio = ignoreAspectRatio

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
	global.Config.ShowBorder = useBorder

	// label
	global.Config.ShowInputLabel = showInputLabel
	global.Config.InputLabelFontSize = inputLabelFontSize

	//motion
	global.Config.UseMotion = useMotion

	// non error case, return nil
	return config, nil
}

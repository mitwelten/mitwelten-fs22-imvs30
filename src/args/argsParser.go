package args

import (
	"fmt"
	"github.com/docopt/docopt.go"
	"log"
	"math"
	"mjpeg_multiplexer/src/aggregator"
	"mjpeg_multiplexer/src/connection"
	"mjpeg_multiplexer/src/customErrors"
	"mjpeg_multiplexer/src/global"
	"mjpeg_multiplexer/src/multiplexer"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	ArgumentSeparator = ","
)

var (
	usage = `Usage:
  multiplexer (grid | carousel | panel) (input) URL (output) PORT [options] 
  multiplexer -h | --help
  multiplexer --version

Options:
  --grid_dimension=ROWS,COLUMNS    todo
  --motion                         Enables motion detection to focus the most active frame on selected mode
  --duration=n                     frame duration [default: 10]
  --cycle                          todo
  --width=n                        output width in pixel [default: -1]
  --height=n                       output height in pixel [default: -1]
  --ignore_aspect_ratio            todo
  --framerate=n                    output framerate in fps[default: -1]
  --quality=n                      output jpeg quality in percentage [default: 100]
  --use_auth                       Use Authentication
  --show_border                    number of black pixels between each image
  --show_label                     Show label for input streams
  --labels=n                       comma separated list of names to show instead of the camera input url
  --label_font_size=n              input label font size in pixel [default: 32]
  --log_fps                       Log Time verbose
  --verbose                        Shows details. 
  --version                        Shows version.
  -h --help                        Shows this screen.
`
)

var (
	helpString = `Usage: multiplexer [grid | panel | carousel] input [URL] output [URL] [options...]
                   <--------- mode --------> <- input -> <- output ->
Mode:
  grid: static grid of images with X rows and Y columns
  panel: dynamic panel of.... Can be used with motion (see --motion)
  carousel: dynamic carousel view.... Can be used with motion (see --motion)
Input:  comma separated list of urls including port
Output: output url including port

Examples: 
  ./multiplexer grid input localhost:8080,localhost:8081 output 8088
  ./multiplexer panel input :8080,:8081,:8082 output 8088 --cycle --width 800 
  ./multiplexer carousel input 192.168.0.1:8080 192.168.0.2:8081 output 8088 --motion

Options:
  --grid_dimension=ROWS,COLUMNS    todo
  --motion                         Enables motion detection to focus the most active frame on selected mode
  --duration=n                     frame duration [default: 10]
  --cycle                          todo
  --width=n                        output width in pixel [default: -1]
  --height=n                       output height in pixel [default: -1]
  --ignore_aspect_ratio            todo
  --framerate=n                    output framerate in fps[default: -1]
  --quality=n                      output jpeg quality in percentage [default: 100]
  --use_auth                       Use Authentication
  --show_border                    number of black pixels between each image
  --show_label                     Show label for input streams
  --labels=n                       comma separated list of names to show instead of the camera input url
  --label_font_size=n              input label font size in pixel [default: 32]
  --log_fps                       Log Time verbose
  --verbose                        Shows details. 
  --version                        Shows version.
  -h --help                        Shows this screen.
`
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
	fmt.Println(helpString)
	os.Exit(1)
}

// ParseArgs parses all arguments derived from command line
func ParseArgs(args []string) (config multiplexer.MultiplexerConfig, err error) {
	// init custom handler to print full usage on error
	parser := &docopt.Parser{
		HelpHandler:  printUsage,
		OptionsFirst: false,
	}
	arguments, err := parser.ParseArgs(usage, nil, "")

	// mode
	grid, _ := arguments.Bool("grid")
	panel, _ := arguments.Bool("panel")
	carousel, _ := arguments.Bool("carousel")
	//input
	input, _ := arguments.String("URL") // input
	//output
	port, _ := arguments.String("PORT") // output

	//options
	gridDimension, _ := arguments.String("--grid_dimension")
	useMotion, _ := arguments.Bool("--motion")
	duration, _ := arguments.Int("--duration") // carousel or panel-cycle duration in seconds
	panelCycle, _ := arguments.Bool("--cycle") // panel cycle, default false
	width, _ := arguments.Int("--width")
	height, _ := arguments.Int("--height")
	ignoreAspectRatio, _ := arguments.Bool("--ignore_aspect_ratio")
	framerate, _ := arguments.Float64("--framerate")
	quality, _ := arguments.Int("--quality")
	showBorder, _ := arguments.Bool("--show_border")
	useAuth, _ := arguments.Bool("--use_auth")
	showInputLabel, _ := arguments.Bool("--show_label")
	inputLabels, _ := arguments.String("--labels") // input inputLabels
	inputLabelFontSize, _ := arguments.Int("--label_font_size")

	logFPS, _ := arguments.Bool("--log_fps")

	// inputURL and label parsing
	config = parseInputUrls(config, input)
	if len(inputLabels) != 0 {
		parseSeparatedString(inputLabels)
	}

	// mode
	if grid {
		var gridX int
		var gridY int
		if len(gridDimension) == 0 {
			gridX = int(math.Ceil(math.Sqrt(float64(len(global.Config.InputConfigs)))))
			gridY = gridX
		} else {
			arr := strings.Split(gridDimension, ArgumentSeparator)
			gridX, err = strconv.Atoi(arr[0])
			if err != nil {
				log.Fatalf("Invalid grid dimension input %v: %v", arr[0], err.Error())
			}
			gridY, err = strconv.Atoi(arr[1])
			if err != nil {
				log.Fatalf("Invalid grid dimension input %v: %v", arr[1], err.Error())
			}
		}
		config.Aggregator = &aggregator.AggregatorGrid{Row: gridX, Col: gridY}
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
	global.Config.LogFPS = logFPS

	// InputRates
	global.Config.OutputFramerate = framerate

	// quality
	global.Config.EncodeQuality = quality

	// border
	global.Config.ShowBorder = showBorder

	// label
	global.Config.ShowInputLabel = showInputLabel
	global.Config.InputLabelFontSize = inputLabelFontSize

	//motion
	global.Config.UseMotion = useMotion

	// non error case, return nil
	return config, nil
}

package args

import (
	"errors"
	"fmt"
	"github.com/docopt/docopt.go"
	"golang.org/x/exp/slices"
	"log"
	"math"
	"mjpeg_multiplexer/src/aggregator"
	"mjpeg_multiplexer/src/global"
	"mjpeg_multiplexer/src/input"
	"mjpeg_multiplexer/src/multiplexer"
	"mjpeg_multiplexer/src/output"
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
  --duration=n                     frame duration [default: 15]
  --cycle                          todo
  --width=n                        output width in pixel [default: -1]
  --height=n                       output height in pixel [default: -1]
  --ignore_aspect_ratio            todo
  --framerate=n                    output framerate in fps[default: -1]
  --quality=n                      output jpeg quality in percentage [default: -1]
  --use_auth                       Use Authentication
  --show_border                    number of black pixels between each image
  --show_label                     Show label for input streams
  --labels=n                       comma separated list of names to show instead of the camera input url
  --label_font_size=n              input label font size in pixel [default: 32]
  --log_fps                        Log Time verbose
  --verbose                        Shows details. 
  --version                        Shows version.
  --always_active                  (hidden) Disables the 'fast mode' when no client is connected
  --disable_passthrough            (hidden) Disables passthrough mode
  --debug                          (hidden)
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
  --grid_dimension=ROWS,COLUMNS    Number of cells in the grid mode
  --motion                         Enables motion detection to focus the most active frame on selected mode
  --cycle                          Enables cycling of the panel layout, see also [--duration] 
  --duration=n                     Duration in seconds before changing the layout (panel and carousel only) [default: 15]
  --width=n                        Output width in pixel [default: -1]
  --height=n                       Output height in pixel [default: -1]
  --ignore_aspect_ratio            Stretches the frames instead of adding a letterbox on resize
  --framerate=n                    Output framerate in fps[default: -1]
  --quality=n                      Output jpeg quality in percentage [default: -1]
  --use_auth                       Use Authentication
  --show_border                    Enables a border in the grid and panel layout between the images
  --show_label                     Show label for input streams
  --labels=n                       Comma separated list of names to show instead of the camera input url
  --label_font_size=n              Input label font size in pixel [default: 32]
  --log_fps                        Logs the current FPS 
  --verbose                        Shows details. 
  --version                        Shows version.
  -h --help                        Shows this screen.
`
)

// parseInput parses input URLS derived from command line arguments
func parseInputUrls(config *multiplexer.MultiplexerConfig, inputStr string) {
	arr := strings.Split(inputStr, ArgumentSeparator)
	config.Inputs = []input.Input{}
	for i, url := range arr {
		config.Inputs = append(config.Inputs, input.NewInputHTTP(i, url))
	}
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

var mode = []string{"grid", "panel", "carousel"}
var optionalFlags = []string{"--cycle", "--ignore_aspect_ratio", "--use_auth", "--show_border", "--show_label", "--log_fps", "--verbose", "--version", "--always_active", "--debug"}
var optionalValues = []string{"--grid_dimension", "--duration", "--width", "--height", "--framerate", "--quality", "--label_font_size"}

func containsHelp(args []string) bool {
	return slices.Contains(args, "--help") || slices.Contains(args, "-h")
}

func containsMode(args []string) bool {
	return slices.Contains(args, "grid") || slices.Contains(args, "panel") || slices.Contains(args, "carousel")
}

func containsInput(args []string) bool {
	return slices.Contains(args, "input")
}

func containsOutput(args []string) bool {
	return slices.Contains(args, "output")
}
func checkOptions(args []string) int {

	for i := 0; i < len(args); i++ {
		el := args[i]

		// mode
		if slices.Contains(mode, el) {
			continue
		}

		// in- and output (skip the specified values)
		if el == "input" || el == "output" {
			i++
			continue
		}

		// flags
		if slices.Contains(optionalFlags, el) {
			continue
		}

		// value fields (skip the specified values)
		if slices.Contains(optionalValues, el) {
			if i == len(args)-1 {
				abort(fmt.Sprintf("Missing value for option '%s'.", el))
			}
			value := args[i+1]

			if el == "--grid-dimension" {
				arr := strings.Split(value, ArgumentSeparator)
				if len(arr) != 2 {
					abort(fmt.Sprintf("Malformed argument '%s'. Expected numerical arguments in '%s x', but got '%s %s'", el, el, el, value))
				}
				_, err1 := strconv.Atoi(arr[0])
				_, err2 := strconv.Atoi(arr[1])

				if err1 != nil || err2 != nil {
					abort(fmt.Sprintf("Malformed argument '%s'. Expected numerical arguments in '%s x', but got '%s %s'", el, el, el, value))
				}

			} else {
				_, err := strconv.Atoi(value)
				if err != nil {
					abort(fmt.Sprintf("Malformed argument '%s'. Expected numerical arguments in '%s x', but got '%s %s'", el, el, el, value))
				}
			}

			i++
			continue
		}

		return i
	}
	return -1
}

func abort(msg string) {
	fmt.Printf("%s See help by adding -h or --help for more information\n", msg)
	os.Exit(1)
}

func validateArgs() {
	args := os.Args[1:]

	if len(args) == 0 || containsHelp(args) {
		printUsage(nil, "")
	}

	if !containsMode(args) {
		abort("Mode missing! Please specify a mode [grid|panel|carousel].")
	}

	if !containsInput(args) {
		abort("Input missing! Please specify an input (eg. 'input localhost:8080').")
	}

	if !containsOutput(args) {
		abort("Output missing! Please specify an output (eg. 'output 8088').")
	}

	i := checkOptions(args)
	if i != -1 {
		abort(fmt.Sprintf("Invalid option '%s'.", args[i]))
	}

}

// ParseArgs parses all arguments derived from command line
func ParseArgs() (config multiplexer.MultiplexerConfig, err error) {
	//todo validate args by parsing os.Args
	validateArgs()

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

	//hidden
	alwaysActive, _ := arguments.Bool("--always_active")
	disablePassthrough, _ := arguments.Bool("--disable_passthrough")
	enableDebug, _ := arguments.Bool("--debug")

	// inputURL and label parsing
	parseInputUrls(&config, input)
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
		return multiplexer.MultiplexerConfig{}, errors.New("invalid mode")
	}

	// output stream
	config.Output = output.NewOutputHTTP(port)

	// disabled: output file
	// config.Output = connection.NewOutputFile(outputFileNamePtr)

	global.Config.Width = width
	global.Config.Height = height
	global.Config.IgnoreAspectRatio = ignoreAspectRatio

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

	//hidden
	global.Config.AlwaysActive = alwaysActive
	global.Config.DisablePassthrough = disablePassthrough
	global.Config.Debug = enableDebug

	// non error case, return nil
	return config, nil
}

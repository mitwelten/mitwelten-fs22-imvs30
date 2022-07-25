package args

import (
	"errors"
	"fmt"
	"github.com/docopt/docopt.go"
	"golang.org/x/exp/slices"
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
  -v --version                     Shows version.
  -h --help                        Shows this screen.`
)

// parseInput parses input URLS derived from command line arguments
// return the number of parsed input urls
func parseInputUrls(config *multiplexer.MultiplexerConfig, inputStr string) int {
	arr := strings.Split(inputStr, ArgumentSeparator)
	config.Inputs = []input.Input{}
	config.InputConfigs = []global.InputConfig{}
	for i, url := range arr {
		config.Inputs = append(config.Inputs, input.NewInputHTTP(i, url))
		config.InputConfigs = append(config.InputConfigs, global.InputConfig{Url: url, Label: url})
	}
	return len(arr)
}

// todo: evtl. trim
func parseSeparatedString(config *multiplexer.MultiplexerConfig, inputStr string) error {
	arr := strings.Split(inputStr, ArgumentSeparator)
	if len(config.Inputs) != len(arr) {
		return createUsageError(fmt.Sprintf("%v input location present, but %v labels found\n", len(config.Inputs), len(arr)))
	}
	for i, label := range arr {
		config.InputConfigs[i].Label = label
	}
	return nil
}

var printUsage = func(err error, usage_ string) {
	fmt.Println(helpString)
	os.Exit(1)
}

var mode = []string{"grid", "panel", "carousel"}
var inOutput = []string{"input", "output"}
var optionalFlags = []string{"--motion", "--cycle", "--ignore_aspect_ratio", "--use_auth", "--show_border", "--show_label", "--log_fps", "--always_active", "--debug"}
var optionalValues = []string{"--grid_dimension", "--duration", "--width", "--height", "--framerate", "--quality", "--label_font_size", "--labels"}

func containsHelp(args []string) bool {
	return slices.Contains(args, "--help") || slices.Contains(args, "-h")
}

func containsVersion(args []string) bool {
	return slices.Contains(args, "--version") || slices.Contains(args, "-v")
}

func createUsageError(msg string) error {
	return errors.New(fmt.Sprintf("%s See help by adding -h or --help for more information.\n", msg))
}

func checkMode(args []string) error {
	if !slices.Contains(args, "grid") && !slices.Contains(args, "panel") && !slices.Contains(args, "carousel") {
		return createUsageError("Mode missing! Please specify a mode [grid|panel|carousel].")
	}
	return nil
}

func checkInput(args []string) error {
	errMsg := "Input missing! Please specify an input with a value (eg. 'input localhost:8080')."

	index := slices.Index(args, "input")
	if index < 0 {
		return createUsageError(errMsg)
	}

	// createUsageError if 'input' is the last word or follow by a keyword
	if index == len(args)-1 || isKeyWord(args[index+1]) {
		return createUsageError(errMsg)
	}

	return nil
}

func checkOutput(args []string) error {
	errMsg := "Output missing! Please specify an output (eg. 'output 8088')."

	index := slices.Index(args, "output")
	if index < 0 {
		return createUsageError(errMsg)
	}

	// createUsageError if 'input' is the last word or follow by a keyword
	if index == len(args)-1 || isKeyWord(args[index+1]) {
		return createUsageError(errMsg)
	}

	return nil
}

func isKeyWord(value string) bool {
	return slices.Contains(mode, value) || slices.Contains(inOutput, value) || slices.Contains(optionalFlags, value) || slices.Contains(optionalValues, value)
}

func checkOptions(args []string) error {

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
		containsEquals := strings.Contains(el, "=")

		// '--quality=100' or 'quality 100' are both legal
		var field string
		if containsEquals {
			arr := strings.Split(el, "=")
			field = arr[0]
		} else {
			field = el
		}

		if slices.Contains(optionalValues, field) {
			var value string

			if containsEquals {
				arr := strings.Split(el, "=")
				value = arr[1]
				if len(value) == 0 {
					return createUsageError(fmt.Sprintf("Missing value for option '%s'.", field))
				}
			} else {
				if i == len(args)-1 {
					return createUsageError(fmt.Sprintf("Missing value for option '%s'.", field))
				}
				value = args[i+1]
			}

			if isKeyWord(value) {
				return createUsageError(fmt.Sprintf("Missing value for option '%s'.", field))
			}

			// if a value option is followed by another one (eg. '--quality --log_fps') we assume that the value is missing

			// special cases:
			// --grid_dimension x,y
			// --label text1,text2,...,textN

			if field == "--grid_dimension" {
				arr := strings.Split(value, ArgumentSeparator)
				if len(arr) != 2 {
					return createUsageError(fmt.Sprintf("Malformed argument '--grid_dimension'. Expected numerical arguments in '--grid_dimension x,y', but got '--grid_dimension %s'", value))
				}
			} else if field != "--labels" {
				_, err := strconv.Atoi(value)
				if err != nil {
					return createUsageError(fmt.Sprintf("Malformed argument '%s'. Expected numerical arguments in '%s x', but got '%s'", el, field, el))
				}
			}

			if !containsEquals {
				i++
			}
			continue
		}

		return createUsageError(fmt.Sprintf("Invalid option '%s'.", args[i]))
	}
	return nil

}

//validateArgs Checks the program args for syntactic errors before passing it into DocOpts
// The following things will be checked:
// - `-h` or `--help`
// - `-v` or `--version`
// - empty argument list
// - missing mode
// - missing input or input parameter
// - missing output or output parameter
// - invalid option or missing option parameter
//
// The method will terminate the program in the first two cases with return code 0 or return an error (or nil if everything is ok) in the other cases
func validateArgs(args []string) error {

	// check for help or version string
	if containsHelp(args) {
		fmt.Println(helpString)
		os.Exit(0)
	}

	if containsVersion(args) {
		fmt.Printf("mjpeg_multiplexer version %v\n", multiplexer.Version)
		os.Exit(0)
	}

	//no args => error
	if len(args) == 0 {
		return errors.New(helpString)
	}

	var err error

	// one mode must be present
	err = checkMode(args)
	if err != nil {
		return err
	}

	// input must be present
	err = checkInput(args)
	if err != nil {
		return err
	}

	// output must be present
	err = checkOutput(args)
	if err != nil {
		return err
	}

	// check for invalid options or missing option parameters
	err = checkOptions(args)
	if err != nil {
		return err
	}

	return nil

}

// ParseArgs parses all arguments derived from command line
func ParseArgs(args []string) (config multiplexer.MultiplexerConfig, err error) {
	err = validateArgs(args)
	if err != nil {
		return multiplexer.MultiplexerConfig{}, err
	}

	// init custom handler to print full usage on error
	parser := &docopt.Parser{
		HelpHandler:  printUsage,
		OptionsFirst: false,
	}
	arguments, err := parser.ParseArgs(usage, args, "")

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
	nInputs := parseInputUrls(&config, input)
	if len(inputLabels) != 0 {
		err = parseSeparatedString(&config, inputLabels)
		if err != nil {
			return multiplexer.MultiplexerConfig{}, err
		}
	}

	// mode
	if grid {

		if useMotion {
			return multiplexer.MultiplexerConfig{}, createUsageError("Option '--motion' only available for the modes 'panel' or 'carousel'.")

		}

		var gridX int
		var gridY int
		if len(gridDimension) == 0 {
			gridX = int(math.Ceil(math.Sqrt(float64(len(config.Inputs)))))
			gridY = gridX
		} else {
			arr := strings.Split(gridDimension, ArgumentSeparator)
			var err1, err2 error
			gridX, err1 = strconv.Atoi(arr[0])
			gridY, err2 = strconv.Atoi(arr[1])

			if err1 != nil || err2 != nil {
				return multiplexer.MultiplexerConfig{}, createUsageError(fmt.Sprintf("Malformed argument '--grid_dimension'. Expected numerical arguments in '--grid_dimension x,y', but got '--grid_dimension %s'", gridDimension))
			}
		}
		config.Aggregator = &aggregator.AggregatorGrid{Row: gridX, Col: gridY}
	} else if carousel {
		config.Aggregator = &aggregator.AggregatorCarousel{Duration: time.Duration(duration) * time.Second}
	} else if panel {
		if nInputs > 8 {
			return multiplexer.MultiplexerConfig{}, errors.New(fmt.Sprintf("Too many inputs for panel mode: %v inputs specified, but only up to 8 are supported. Use `carousel` or `grid` instead.\n", nInputs))

		}
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

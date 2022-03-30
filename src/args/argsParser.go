package args

import (
	"flag"
	"mjpeg_multiplexer/src/aggregator"
	"mjpeg_multiplexer/src/connection"
	"mjpeg_multiplexer/src/customErrors"
	"mjpeg_multiplexer/src/multiplexer"
	"os"
	"strconv"
	"strings"
)

const (
	InputUsage             = "Use -input to determine input streams. Pattern: [IP][:][PORT][SPACE][IP][:][PORT]..."
	OutputUsage            = "Use -output to determine output modus. 'file' or 'stream' possible."
	OutputFileNameUsage    = "filename used to save input to file"
	OutputStreamPortUsage  = "port used for output stream"
	MethodUsage            = "Method, how the output will be mixed. 'grid' is possible. "
	GridUsage              = "Number of rows and columns for the grid. format: '<row> <col>'"
	InputLocationSeparator = " "
)

func parseGrid(config multiplexer.MultiplexerConfig, methodGridPtr *string) (multiplexer.MultiplexerConfig, error) {
	var gridDimension = strings.Split(*methodGridPtr, InputLocationSeparator)
	if len(gridDimension) != 2 {
		return config, &customErrors.ErrArgParserInvalidGridDimension{}
	}

	row, err := strconv.Atoi(gridDimension[0])
	if err != nil {
		return config, &customErrors.ErrArgParserInvalidGridDimension{}
	}
	col, err := strconv.Atoi(gridDimension[1])
	if err != nil {
		return config, &customErrors.ErrArgParserInvalidGridDimension{}
	}
	config.Aggregator = aggregator.AggregatorGrid{Row: row, Col: col}

	return config, nil
}

func parseInput(config multiplexer.MultiplexerConfig, inputStr string) multiplexer.MultiplexerConfig {
	config.InputLocations = strings.Split(inputStr, InputLocationSeparator)
	return config
}

func ParseArgs(args []string) (config multiplexer.MultiplexerConfig, err error) {
	var CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	//---Define all various flags---
	inputPtr := CommandLine.String("input", "", InputUsage)
	outputPtr := CommandLine.String("output", "", OutputUsage) // file oder stream
	outputFileNamePtr := CommandLine.String("output_filename", "", OutputFileNameUsage)
	outputStreamPortPtr := CommandLine.String("output_port", "", OutputStreamPortUsage)
	methodPtr := CommandLine.String("method", "", MethodUsage) // grid, combine etc.
	methodGridPtr := CommandLine.String("grid_dimension", "", GridUsage)

	//---parse the command line into the defined flags---
	CommandLine.Parse(args[1:])

	// first validation
	// check if at least all three mandatory parameters are present
	if len(*inputPtr) == 0 || len(*outputPtr) == 0 || len(*methodPtr) == 0 {
		return multiplexer.MultiplexerConfig{}, &customErrors.ErrArgParserUnfulfilledMinArguments{}
	}
	// stream
	if strings.Compare(*outputPtr, "stream") == 0 {
		if len(*outputStreamPortPtr) == 0 {
			return multiplexer.MultiplexerConfig{}, &customErrors.ErrArgParserInvalidOutputPort{}
		} else {
			config.Output, err = connection.NewOutputHTTP(*outputStreamPortPtr)
			if err != nil {
				return multiplexer.MultiplexerConfig{}, &customErrors.ErrHttpOpenOutputSocket{}
			}
		}
		// file
	} else if strings.Compare(*outputPtr, "file") == 0 {
		if len(*outputFileNamePtr) == 0 {
			return multiplexer.MultiplexerConfig{}, &customErrors.ErrArgParserInvalidOutputFilename{}
		} else {
			config.Output = connection.NewOutputFile(*outputFileNamePtr)
		}
	} else {
		return multiplexer.MultiplexerConfig{}, &customErrors.ErrArgParserInvalidArgument{Argument: *outputPtr}
	}

	// input parsing
	config = parseInput(config, *inputPtr)

	// method
	switch *methodPtr {
	case "grid":
		config, err = parseGrid(config, methodGridPtr)

		if err != nil {
			return multiplexer.MultiplexerConfig{}, err
		}
	default:
		return multiplexer.MultiplexerConfig{}, &customErrors.ErrArgParserInvalidMethod{Argument: *methodPtr}
	}

	// non error case, return nil

	return config, nil
}

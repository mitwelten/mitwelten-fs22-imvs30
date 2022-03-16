package args

import (
	"flag"
	"mjpeg_multiplexer/src/connection"
	"mjpeg_multiplexer/src/customErrors"
	"os"
	"strings"
)

const (
	InputUsage             = "Use -input to determine input streams. Pattern: [IP][:][PORT][SPACE][IP][:][PORT]..."
	OutputUsage            = "Use -output to determine output modus. 'file' or 'stream' possible."
	OutputFileNameUsage    = "filename used to save input to file"
	OutputStreamPortUsage  = "port used for output stream"
	MethodUsage            = "Method, how the output will be mixed. 'combine' or 'grid' possible. "
	InputLocationSeparator = " "
)

type MultiplexerConfig struct {
	InputLocations []string
	Output         connection.Output
}

func parseInput(config MultiplexerConfig, inputStr string) {
	config.InputLocations = strings.Split(inputStr, InputLocationSeparator)
}

func ParseArgs(args []string) (config MultiplexerConfig, err error) {
	var CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	//---Define all various flags---
	inputPtr := CommandLine.String("input", "", InputUsage)
	outputPtr := CommandLine.String("output", "", OutputUsage) // file oder stream
	outputFileNamePtr := CommandLine.String("output_filename", "", OutputFileNameUsage)
	outputStreamPortPtr := CommandLine.String("output_port", "", OutputStreamPortUsage)
	methodPtr := CommandLine.String("method", "", MethodUsage) // grid, combine etc.

	//---parse the command line into the defined flags---
	CommandLine.Parse(args[1:])

	// first validation
	// check if at least all three mandatory parameters are present
	if len(*inputPtr) == 0 || len(*outputPtr) == 0 || len(*methodPtr) == 0 {
		return MultiplexerConfig{}, &customErrors.ErrArgParserUnfulfilledMinArguments{}
	}
	// stream
	if strings.Compare(*outputPtr, "stream") == 0 {
		if len(*outputStreamPortPtr) == 0 {
			return MultiplexerConfig{}, &customErrors.ErrArgParserInvalidOutputPort{}
		} else {
			config.Output, err = connection.NewOutputHTTP(*outputStreamPortPtr)
			if err != nil {
				return MultiplexerConfig{}, &customErrors.ErrHttpOpenOutputSocket{}
			}
		}
		// file
	} else if strings.Compare(*outputPtr, "file") == 0 {
		if len(*outputFileNamePtr) == 0 {
			return MultiplexerConfig{}, &customErrors.ErrArgParserInvalidOutputFilename{}
		} else {
			config.Output = connection.NewOutputFile(*outputFileNamePtr)
		}
	} else {
		return MultiplexerConfig{}, &customErrors.ErrArgParserInvalidArgument{Argument: *outputPtr}
	}

	// input parsing
	parseInput(config, *inputPtr)

	// non error case, return nil
	return config, nil
}

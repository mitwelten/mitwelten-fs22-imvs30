package args

import (
	"errors"
	"flag"
	"fmt"
	"mjpeg_multiplexer/src/connection"
	"mjpeg_multiplexer/src/customErrors"
	"os"
	"strings"
)

type MultiplexerConfig struct {
	InputLocations []string
	Output         connection.Output
}

func parseInput(config MultiplexerConfig, inputStr string) {
	config.InputLocations = strings.Split(inputStr, " ")
}

func ParseArgs(args []string) (config MultiplexerConfig, err error) {
	var CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	//---Define all various flags---
	inputPtr := CommandLine.String("input", "", "Use Input flag to determine input streams")
	outputPtr := CommandLine.String("output", "", "Use Out flag to determine output modus") // file oder stream
	outputFileNamePtr := CommandLine.String("output_filename", "", "")
	outputStreamPortPtr := CommandLine.String("output_port", "", "port for output stream")
	methodPtr := CommandLine.String("method", "", "How will the out be mixed?") // grid, combine etc.

	//---parse the command line into the defined flags---
	//flag.Parse()
	CommandLine.Parse(args[1:])
	// first validation
	// check if at least all three mandatory parameters are present
	if len(*inputPtr) == 0 || len(*outputPtr) == 0 || len(*methodPtr) == 0 {
		println("wrong here...")
		return MultiplexerConfig{}, &customErrors.ArgParserUnfulfilledMinArgumentsError{}
	}
	// stream
	if strings.Compare(*outputPtr, "stream") == 0 {
		if len(*outputStreamPortPtr) == 0 {
			return MultiplexerConfig{}, &customErrors.ArgParserInvalidInputError{}
		} else {
			config.Output, err = connection.NewOutputHTTP(*outputStreamPortPtr)
			if err != nil {
				return MultiplexerConfig{}, errors.New("can't open HTTP output")
			}
		}
		// file
	} else if strings.Compare(*outputPtr, "file") == 0 {
		if len(*outputFileNamePtr) == 0 {
			return MultiplexerConfig{}, errors.New("-output 'file' only valid in combination with -output_filename ")
		} else {
			config.Output = connection.NewOutputFile(*outputFileNamePtr)
		}
	} else {
		return MultiplexerConfig{}, errors.New("invalid output argument: -output argument '" + *outputPtr + "' not valid. Use -output 'stream' or -output 'file'")
	}

	// input parsing
	parseInput(config, *inputPtr)

	fmt.Println(*inputPtr)
	fmt.Println(*outputPtr)
	fmt.Println(*outputFileNamePtr)
	fmt.Println(*outputStreamPortPtr)
	fmt.Println(*methodPtr)

	return config, nil
}

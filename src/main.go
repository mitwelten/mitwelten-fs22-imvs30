package main

import (
	"flag"
	"fmt"
	"mjpeg_multiplexer/src/aggregator"
	"mjpeg_multiplexer/src/connection"
	"mjpeg_multiplexer/src/mjpeg"
	"os"
	"strings"
	"sync"
)

// as per https://docs.fileformat.com/image/jpeg/

var InputLocations []string
var Output connection.Output

func run() {
	var wg sync.WaitGroup

	var channels []chan mjpeg.Frame

	for _, connectionString := range InputLocations {
		wg.Add(1)

		var input = connection.NewInputHTTP(connectionString)
		input.Open()
		var channel = connection.ListenToInput(input)
		channels = append(channels, channel)
	}

	//var aggregatedChannels = aggregator.Merge2Images(channels)
	//var aggregatedChannels = aggregator.MergeImages(channels)
	var aggregatedChannels = aggregator.CombineChannels(channels)

	wg.Add(1)
	//var output = connection.NewOutputFile("out.jpg")
	//var output = connection.NewOutputHTTP("8082")
	connection.RunOutput(Output, aggregatedChannels)

	wg.Wait()
}

func parseInput(inputStr string) {
	InputLocations = strings.Split(inputStr, " ")
}

func parseArgs() {
	//---Define all various flags---
	inputPtr := flag.String("input", "", "Use Input flag to determine input streams")
	outputPtr := flag.String("output", "", "Use Out flag to determine output modus") // file oder stream
	outputFileNamePtr := flag.String("output_filename", "", "")
	outputStreamPortPtr := flag.String("output_port", "", "port for output stream")
	methodPtr := flag.String("method", "", "How will the out be mixed?") // grid, combine etc.

	//---parse the command line into the defined flags---
	flag.Parse()

	// first validation
	// check if at least all three mandatory parameters are present
	if len(*inputPtr) == 0 || len(*outputPtr) == 0 || len(*methodPtr) == 0 {
		fmt.Println("expected at least '-input' '-output' and '-method' arguments")
		os.Exit(2)
	}
	// stream
	if strings.Compare(*outputPtr, "stream") == 0 {
		if len(*outputStreamPortPtr) == 0 {
			fmt.Println("-output 'stream' only valid in combination with -output_port ")
			os.Exit(3)
		} else {
			Output = connection.NewOutputFile(*outputStreamPortPtr)
		}
		// file
	} else if strings.Compare(*outputPtr, "file") == 0 {
		if len(*outputFileNamePtr) == 0 {
			fmt.Println("-output 'file' only valid in combination with -output_filename ")
			os.Exit(4)
		} else {
			Output = connection.NewOutputFile(*outputFileNamePtr)
		}
	} else {
		fmt.Println("invalid output argument: -output argument '" + *outputPtr + "' not valid. Use -output 'stream' or -output 'file'")
		os.Exit(5)
	}

	// input parsing
	parseInput(*inputPtr)

	fmt.Println(*inputPtr)
	fmt.Println(*outputPtr)
	fmt.Println(*outputFileNamePtr)
	fmt.Println(*outputStreamPortPtr)
	fmt.Println(*methodPtr)

}
func main() {
	println("Running the MJPEG-multiFLEXer")
	parseArgs()
	run()
}

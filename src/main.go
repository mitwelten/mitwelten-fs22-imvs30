package main

import (
	"fmt"
	"mjpeg_multiplexer/src/aggregator"
	"mjpeg_multiplexer/src/args"
	"mjpeg_multiplexer/src/connection"
	"mjpeg_multiplexer/src/mjpeg"
	"os"
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

func main() {
	println(os.Args)

	// loop over all arguments by index and value
	for i, arg := range os.Args {
		// print index and value
		fmt.Println("item", i, "is", arg)
	}

	println("Running the MJPEG-multiFLEXer")
	config, err := args.ParseArgs(os.Args)
	if err != nil {
		panic(err.Error())
	}

	InputLocations = config.InputLocations
	Output = config.Output

	run()
}

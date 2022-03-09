package main

import (
	"mjpeg_multiplexer/src/aggregator"
	"mjpeg_multiplexer/src/connection"
	"mjpeg_multiplexer/src/mjpeg"
)

import (
	"os"
	"sync"
)

// as per https://docs.fileformat.com/image/jpeg/

func run(args []string) {
	var wg sync.WaitGroup

	var channels []chan mjpeg.Frame

	for _, connectionString := range args {
		wg.Add(1)

		var input = connection.NewInputHTTP(connectionString)
		input.Open()
		var channel = connection.RunSource(input)
		channels = append(channels, channel)
	}

	//var aggregatedChannels = aggregator.Merge2Images(channels)
	//var aggregatedChannels = aggregator.MergeImages(channels)
	var aggregatedChannels = aggregator.CombineChannels(channels)

	wg.Add(1)
	//var output = connection.NewOutputFile("out.jpg")
	var output = connection.NewOutputHTTP("8081")
	connection.RunOutput(output, aggregatedChannels)

	wg.Wait()
}

func main() {
	println("Running the MJPEG-multiFLEXer")
	var args = os.Args[1:]
	run(args)
}

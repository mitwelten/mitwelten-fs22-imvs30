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
	for _, port := range args {
		wg.Add(1)
		var source = connection.NewHTTPSource(port)
		source.Open()
		var channel = connection.RunSource(source)
		channels = append(channels, channel)
	}
	var aggregatedChannels = aggregator.MergeImages(channels)
	//var aggregatedChannels = aggregator.CombineChannels(channels)

	wg.Add(1)
	var sink = connection.NewFileSink("out.jpg")
	//var sink = connection.NewHTTPSink("8082")
	connection.RunSink(sink, aggregatedChannels)
	//image.Test(sources[0], sources[1])
	wg.Wait()
}

func main() {
	println("Running the MJPEG-multiFLEXer")
	var args = os.Args[1:]
	run(args)
}

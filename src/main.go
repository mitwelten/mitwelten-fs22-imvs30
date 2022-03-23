package main

import (
	"fmt"
	"mjpeg_multiplexer/src/args"
	"mjpeg_multiplexer/src/communication"
	"mjpeg_multiplexer/src/connection"
	"os"
	"sync"
)

// as per https://docs.fileformat.com/image/jpeg/

var Config args.MultiplexerConfig

func run() {
	var wg sync.WaitGroup

	//	var channels []chan mjpeg.Frame
	var frameDatas []*communication.FrameData

	for _, connectionString := range Config.InputLocations {
		wg.Add(1)

		var input = connection.NewInputHTTP(connectionString)
		input.Open()
		var frameData = connection.ListenToInput(input)
		frameDatas = append(frameDatas, frameData)
	}

	var aggregatedChannels = Config.Aggregator.Aggregate(frameDatas...)

	wg.Add(1)
	connection.RunOutput(Config.Output, aggregatedChannels)

	wg.Wait()
}

func main() {
	// loop over all arguments by index and value
	for i, arg := range os.Args {
		// print index and value
		fmt.Println("item", i, "is", arg)
	}

	println("Running the MJPEG-multiFLEXer")
	c, err := args.ParseArgs(os.Args)
	if err != nil {
		panic(err.Error())
	}

	Config = c

	run()
}

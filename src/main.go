package main

import (
	"io"
	"log"
	"mjpeg_multiplexer/src/args"
	"mjpeg_multiplexer/src/communication"
	"mjpeg_multiplexer/src/connection"
	"os"
	"sync"
)

var Config args.MultiplexerConfig

const logFile string = "multiplexer.log"

func setupLog() {
	// log setup
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	mw := io.MultiWriter(os.Stdout, file)
	log.SetOutput(mw)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

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
	setupLog()
	log.Println("Running the MJPEG-multiFLEXer")

	log.Println("parsing args...")
	// loop over all arguments by index and value
	for i, arg := range os.Args {
		// print index and value
		log.Println("item", i, "is", arg)
	}

	c, err := args.ParseArgs(os.Args)
	if err != nil {
		panic(err.Error())
	}

	Config = c

	run()
}

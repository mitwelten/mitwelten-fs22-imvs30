package main

import (
	"io"
	"log"
	"mjpeg_multiplexer/src/args"
	"mjpeg_multiplexer/src/multiplexer"
	"os"
)

var Config multiplexer.MultiplexerConfig

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

	multiplexer.Multiplexer(Config)
}

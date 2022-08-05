package main

import (
	"fmt"
	"io"
	"log"
	"mjpeg_multiplexer/src/args"
	"mjpeg_multiplexer/src/global"
	"mjpeg_multiplexer/src/multiplexer"
	"net/http"
	_ "net/http/pprof"
	"os"
)

const logFile string = "multiplexer.log"

// setupLog opens the log file and configures the log class such that all logs will be written to the terminal + log file
func setupLog() {
	// log setup
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	mw := io.MultiWriter(os.Stdout, file)
	log.SetOutput(mw)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

// profile starts the profiler
// usage: view with `go tool pprof http://localhost:6060/debug/pprof/profile\?seconds\=30`
func profile() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
}

// main entry point: setup, argParsing, start multiplexer
func main() {
	profile()
	setupLog()
	global.SetupInitialConfig()

	multiplexerConfig, err := args.ParseArgs(os.Args[1:])
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	multiplexer.Multiplexer(multiplexerConfig)
}

package main

import (
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

func profile() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	//view with `go tool pprof http://localhost:6060/debug/pprof/profile\?seconds\=30`
}

func main() {
	profile()
	setupLog()
	global.SetupInitialConfig()

	log.Println("Running the MJPEG-multiFLEXer")
	log.Println("parsing args...")
	// loop over all arguments by index and value
	for i, arg := range os.Args {
		// print index and value
		log.Println("item", i, "is", arg)
	}

	multiplexerConfig, err := args.ParseArgs(os.Args)
	if err != nil {
		panic(err.Error())
	}

	log.Printf("Global config: %+v\n", global.Config)
	log.Printf("Multiplexer config: %+v\n", multiplexerConfig)

	multiplexer.Multiplexer(multiplexerConfig)
}

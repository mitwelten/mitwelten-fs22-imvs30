package main

import (
	"mjpeg_multiplexer/src/connection"
)

import (
	"os"
	"sync"
)

// as per https://docs.fileformat.com/image/jpeg/

func run(args []string) {
	var wg sync.WaitGroup

	var sources []connection.Source
	for _, port := range args {
		wg.Add(1)
		var source = connection.NewHTTPSource(port)
		source.Open()
		connection.RunSource(source)
		sources = append(sources, source)
	}

	wg.Add(1)
	var sink = connection.NewFileSink("out_.jpg")
	connection.RunSink(sink, sources)

	wg.Wait()
}

func main() {
	println("Running the MJPEG-multiFLEXer")
	var args = os.Args[1:]
	run(args)
}

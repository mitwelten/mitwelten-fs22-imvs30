//https://go101.org/article/channel-use-cases.html
package main

import "mjpeg_multiplexer/mjpeg"
import "mjpeg_multiplexer/connection"

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
		var source = connection.NewSource(port)
    source.Open()
		var channel = source.Run()
		channels = append(channels, channel)
	}

	wg.Add(1)
	var sink = connection.NewSink("out_.jpg")
	sink.Run(channels)

	wg.Wait()
}

func main() {
	var args = os.Args[1:]
	run(args)
}

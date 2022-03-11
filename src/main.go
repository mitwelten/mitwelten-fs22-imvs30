package main

import (
	"flag"
	"fmt"
	"mjpeg_multiplexer/src/aggregator"
	"mjpeg_multiplexer/src/connection"
	"mjpeg_multiplexer/src/mjpeg"
	"strings"
)

import (
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

func parseArgs() {
	fileCmd := flag.NewFlagSet("file", flag.ExitOnError)
	fileName := fileCmd.String("name", "", "filename")

	streamCmd := flag.NewFlagSet("stream", flag.ExitOnError)
	streamPort := streamCmd.String("port", "", "port address")

	inputCmd := flag.NewFlagSet("input", flag.ExitOnError)
	inputLocations := inputCmd.String("locations", "", "string list of locations")

	if len(os.Args) < 2 {
		fmt.Println("expected 'file' or 'stream' subcommands")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "file":
		fileCmd.Parse(os.Args[2:4])
		Output = connection.NewOutputFile(*fileName)
	case "stream":
		streamCmd.Parse(os.Args[2:4])
		var err error
		Output, err = connection.NewOutputHTTP(*streamPort)
		if err != nil {
			println(err.Error())
			panic("Could not open output socket, aborting...")
		}
	default:
		fmt.Println("expected 'file' or 'stream' subcommands")
		os.Exit(1)
	}

	inputCmd.Parse(os.Args[5:])
	InputLocations = strings.Split(*inputLocations, " ")

}
func main() {
	println("Running the MJPEG-multiFLEXer")
	parseArgs()
	run()
}

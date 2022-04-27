package integration_tests

import (
	"mjpeg_multiplexer/src/aggregator"
	"mjpeg_multiplexer/src/connection"
	"mjpeg_multiplexer/src/multiplexer"
	"os"
	"testing"
	"time"
)

const FILE_NAME string = "RED_BLUE.jpg"
const FILE_NAME_OUTPUT string = "RED_BLUE_OUTPUT.jpg"

func TestRedBlueMerge(t *testing.T) {
	go func() {
		time.Sleep(3 * time.Second)
		t.Error("frame not received within time limit")
		os.Exit(1)
	}()

	SimpleServer("8097", RedFrame())
	SimpleServer("8098", BlueFrame())

	go func() {
		multiplexer.Multiplexer(multiplexer.MultiplexerConfig{
			InputLocations: []connection.Input{connection.NewInputHTTP("localhost:8097"), connection.NewInputHTTP("localhost:8098")},
			Output:         connection.NewOutputFile(FILE_NAME_OUTPUT),
			Aggregator:     &aggregator.AggregatorGrid{Row: 1, Col: 2},
		})
	}()

	time.Sleep(1 * time.Second)
}

package multiplexer

import (
	"mjpeg_multiplexer/src/aggregator"
	"mjpeg_multiplexer/src/communication"
	"mjpeg_multiplexer/src/connection"
	"sync"
)

type MultiplexerConfig struct {
	InputLocations []string
	Output         connection.Output
	Aggregator     aggregator.Aggregator
}

func Multiplexer(config MultiplexerConfig) {
	var wg sync.WaitGroup

	var frameStorage []*communication.FrameStorage

	for _, connectionString := range config.InputLocations {
		wg.Add(1)

		var input = connection.NewInputHTTP(connectionString)
		input.Open()
		var frameData = connection.ListenToInput(input)
		frameStorage = append(frameStorage, frameData)
	}

	var aggregatedFrameStorage = config.Aggregator.Aggregate(frameStorage...)

	wg.Add(1)
	connection.RunOutput(config.Output, aggregatedFrameStorage)

	wg.Wait()
}

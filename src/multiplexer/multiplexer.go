package multiplexer

import (
	"mjpeg_multiplexer/src/aggregator"
	"mjpeg_multiplexer/src/communication"
	"mjpeg_multiplexer/src/connection"
	"sync"
)

type MultiplexerConfig struct {
	InputLocations []connection.Input
	Output         connection.Output
	Aggregator     aggregator.Aggregator
}

func Multiplexer(config MultiplexerConfig) {
	var wg sync.WaitGroup

	var frameStorage []*communication.FrameStorage

	for _, inputConnection := range config.InputLocations {
		wg.Add(1)

		inputConnection.Start()
		var frameData = connection.ListenToInput(inputConnection)
		frameStorage = append(frameStorage, frameData)
	}

	var aggregatedFrameStorage = config.Aggregator.Aggregate(frameStorage...)

	wg.Add(1)
	config.Output.Run(aggregatedFrameStorage)

	wg.Wait()
}

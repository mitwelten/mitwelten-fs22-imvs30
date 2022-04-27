package multiplexer

import (
	"fmt"
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

		err := inputConnection.Start()
		if err != nil {
			panic(fmt.Sprintf("Can't open input stream: %s", err.Error()))
		}
		var frameData = connection.ListenToInput(inputConnection)
		frameStorage = append(frameStorage, frameData)
	}

	config.Aggregator.Aggregate(frameStorage...)

	wg.Add(1)
	config.Output.Run(config.Aggregator)

	wg.Wait()
}

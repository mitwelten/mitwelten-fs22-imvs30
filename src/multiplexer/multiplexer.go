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

	//	var channels []chan mjpeg.Frame
	var frameDatas []*communication.FrameData

	for _, connectionString := range config.InputLocations {
		wg.Add(1)

		var input = connection.NewInputHTTP(connectionString)
		input.Open()
		var frameData = connection.ListenToInput(input)
		frameDatas = append(frameDatas, frameData)
	}

	var aggregatedChannels = config.Aggregator.Aggregate(frameDatas...)

	wg.Add(1)
	connection.RunOutput(config.Output, aggregatedChannels)

	wg.Wait()
}

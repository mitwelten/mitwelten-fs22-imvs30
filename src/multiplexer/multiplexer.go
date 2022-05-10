package multiplexer

import (
	"fmt"
	"mjpeg_multiplexer/src/aggregator"
	"mjpeg_multiplexer/src/communication"
	"mjpeg_multiplexer/src/connection"
	"mjpeg_multiplexer/src/global"
	"mjpeg_multiplexer/src/utils"
	"sync"
)

type MultiplexerConfig struct {
	InputLocations []connection.Input
	Output         connection.Output
	Aggregator     aggregator.Aggregator
	UseAuth        bool
}

func Multiplexer(multiplexerConfig MultiplexerConfig) {
	var wg sync.WaitGroup

	if multiplexerConfig.UseAuth {
		global.Config.Authentications = utils.ParseAuthenticationFile()
	}

	var frameStorage []*communication.FrameStorage

	for _, inputConnection := range multiplexerConfig.InputLocations {
		wg.Add(1)

		err := inputConnection.Start()
		if err != nil {
			panic(fmt.Sprintf("Can't open input stream: %s", err.Error()))
		}
		var frameData = connection.ListenToInput(inputConnection)
		frameStorage = append(frameStorage, frameData)
	}

	multiplexerConfig.Aggregator.Aggregate(frameStorage...)

	wg.Add(1)
	multiplexerConfig.Output.Run(multiplexerConfig.Aggregator)

	wg.Wait()
}

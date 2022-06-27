package multiplexer

import (
	"log"
	"mjpeg_multiplexer/src/aggregator"
	"mjpeg_multiplexer/src/connection"
	"mjpeg_multiplexer/src/global"
	"mjpeg_multiplexer/src/mjpeg"
	"mjpeg_multiplexer/src/utils"
	"sync"
)

type MultiplexerConfig struct {
	InputLocations []connection.Input
	Output         connection.Output
	Aggregator     aggregator.Aggregator
}

func Multiplexer(multiplexerConfig MultiplexerConfig) {
	log.Println("Running the MJPEG-multiFLEXer")
	var wg sync.WaitGroup

	if global.Config.UseAuth {
		utils.ParseAuthenticationFile()
	}

	var frameStorage []*mjpeg.FrameStorage

	for _, inputConnection := range multiplexerConfig.InputLocations {
		wg.Add(1)
		var frameData = connection.ListenToInput(inputConnection)
		frameStorage = append(frameStorage, frameData)
	}

	aggregator.Aggregate(&multiplexerConfig.Aggregator, frameStorage...)
	multiplexerConfig.Output.Run()
	wg.Add(1)
	wg.Wait()
}

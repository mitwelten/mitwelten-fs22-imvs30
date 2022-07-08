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

	var inputStorages []*mjpeg.FrameStorage

	for _, inputConnection := range multiplexerConfig.InputLocations {
		wg.Add(1)
		var frameData = connection.StartInput(inputConnection)
		inputStorages = append(inputStorages, frameData)
	}

	aggregator.StartAggregator(&multiplexerConfig.Aggregator, inputStorages...)
	multiplexerConfig.Output.StartOutput(&multiplexerConfig.Aggregator)
	wg.Add(1)
	wg.Wait()
}

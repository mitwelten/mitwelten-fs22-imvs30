package multiplexer

import (
	"log"
	"mjpeg_multiplexer/src/aggregator"
	"mjpeg_multiplexer/src/global"
	"mjpeg_multiplexer/src/input"
	"mjpeg_multiplexer/src/output"
	"mjpeg_multiplexer/src/utils"
	"sync"
)

type MultiplexerConfig struct {
	Inputs     []input.Input
	Output     output.Output
	Aggregator aggregator.Aggregator
}

func Multiplexer(multiplexerConfig MultiplexerConfig) {
	log.Println("Running the MJPEG-multiFLEXer")
	var wg sync.WaitGroup

	if global.Config.UseAuth {
		utils.ParseAuthenticationFile()
	}

	for _, inputConnection := range multiplexerConfig.Inputs {
		wg.Add(1)
		input.StartInput(inputConnection)
	}

	aggregator.StartAggregator(&multiplexerConfig.Aggregator, multiplexerConfig.Inputs...)
	multiplexerConfig.Output.StartOutput(&multiplexerConfig.Aggregator)
	wg.Add(1)
	wg.Wait()
}

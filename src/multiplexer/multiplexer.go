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

const Version = "1.0"
const authenticationFileLocation = "authentication.json"

// MultiplexerConfig config struct for specifying the Multiplexer configuration
type MultiplexerConfig struct {
	Inputs       []input.Input
	Output       output.Output
	Aggregator   aggregator.Aggregator
	InputConfigs []global.InputConfig
}

// Multiplexer starts the application with the given config
// will not return
func Multiplexer(multiplexerConfig MultiplexerConfig) {
	log.Printf("Running mjpeg_multiplexer version %v\n", Version)
	var wg sync.WaitGroup

	if global.Config.UseAuth {
		urls := make([]string, len(multiplexerConfig.InputConfigs))
		for i := 0; i < len(urls); i++ {
			urls[i] = multiplexerConfig.InputConfigs[i].Url
		}

		authentications, err := utils.ParseAuthenticationFile(urls, authenticationFileLocation)
		if err != nil {
			log.Fatalln(err.Error())
		}

		for i := 0; i < len(authentications); i++ {
			multiplexerConfig.InputConfigs[i].Authentication = authentications[i]
		}
	}

	global.Config.InputConfigs = multiplexerConfig.InputConfigs

	for _, inputConnection := range multiplexerConfig.Inputs {
		wg.Add(1)
		input.StartInput(inputConnection)
	}

	aggregator.StartAggregator(&multiplexerConfig.Aggregator, multiplexerConfig.Inputs...)
	multiplexerConfig.Output.StartOutput(&multiplexerConfig.Aggregator)
	wg.Add(1)
	wg.Wait()
}

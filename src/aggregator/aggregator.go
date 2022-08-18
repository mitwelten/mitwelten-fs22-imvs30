package aggregator

import (
	"log"
	"mjpeg_multiplexer/src/global"
	"mjpeg_multiplexer/src/input"
	"mjpeg_multiplexer/src/mjpeg"
	"mjpeg_multiplexer/src/utils"
	"reflect"
	"time"
)

const defaultDuration = 15 * time.Second

// Aggregator aggregates multiple frame storages to one frame storage
// takes multiple frameStorages (one frame storage for each input Connection) and
// aggregates them to one frameStorage
type Aggregator interface {
	//Setup initiates all the aggregator's fields
	Setup(storages ...*mjpeg.FrameStorage)
	//GetAggregatorData is the getter for the AggregatorData struct
	GetAggregatorData() *AggregatorData
	//aggregate Combines the input frames into a combined mjpeg-jpeg frame
	aggregate(channels ...*mjpeg.FrameStorage) *mjpeg.MjpegFrame
}

type AggregatorData struct {
	passthrough bool
	//Enabled           bool
	AggregatorStorage *mjpeg.FrameStorage
}

// StartAggregator starts the aggregator loop for the passed aggregator in a separate go routine
func StartAggregator(agg *Aggregator, inputs ...input.Input) {
	var inputStorages []*mjpeg.FrameStorage
	for _, in := range inputs {
		inputStorages = append(inputStorages, in.GetInputData().InputStorage)
	}
	// Setup
	aggregator := *agg
	aggregator.Setup(inputStorages...)
	aggregatorData := aggregator.GetAggregatorData()
	aggregatorData.AggregatorStorage = mjpeg.NewFrameStorage()
	condition := mjpeg.CreateUpdateCondition(inputStorages...)

	lastUpdate := time.Unix(0, 0)
	currentFPS := 0

	passthroughMode := false

	if !global.DecodingNecessary() && reflect.TypeOf(aggregator) == reflect.TypeOf(&AggregatorCarousel{}) {
		passthroughMode = true
	} else if !global.DecodingNecessary() && reflect.TypeOf(aggregator) == reflect.TypeOf(&AggregatorGrid{}) && len(inputStorages) == 1 {
		passthroughMode = true
	}

	if passthroughMode {
		log.Printf("Passthorugh-Mode activate - Each frame will be directly passed to the output without any decoding and encoding\n")
	} else {
		log.Printf("Passthorugh-Mode not active - Each frame will be decoded and encoded, which may be slow.\n=> Activate it by either using the 'carousel' mode or by only using one input source and by removing the width, height, quality and show_label options.\n")
	}

	go func() {
		for {
			condition.Wait()

			// not enabled => no client is connected and no work has to be done
			if !global.Config.AlwaysActive {

				global.Config.AggregatorEnabledMutex.RLock()
				if !global.Config.AggregatorEnabled {
					global.Config.AggregatorEnabledMutex.RUnlock()
					time.Sleep(1 * time.Second)
					continue
				}
				global.Config.AggregatorEnabledMutex.RUnlock()
			}

			if global.Config.OutputFramerate != -1 && time.Since(lastUpdate).Seconds() < (1.0/global.Config.OutputFramerate) {
				continue
			}

			// process the new frame...
			var frame *mjpeg.MjpegFrame
			if aggregatorData.passthrough && len(inputStorages) == 1 && passthroughMode {
				frame = inputStorages[0].GetFrame()
			} else {
				frame = aggregator.aggregate(inputStorages...)
			}

			if frame == nil {
				continue
			}

			// and store it
			if aggregatorData.AggregatorStorage != nil {
				aggregatorData.AggregatorStorage.Store(frame)
			}

			currentFPS++
			secondsNow := time.Now().Second()
			if lastUpdate.Second() != secondsNow {
				if global.Config.LogFPS {
					//we wanna either log
					// X FPS if above zero
					// or 0.XX FPS if below

					//0:50 - 0:05 = 0:45, but should be 0:15
					//fix this by turning 0:05 into 1:05
					if secondsNow < lastUpdate.Second() {
						secondsNow += 60
					}
					deltaSeconds := utils.Abs(lastUpdate.Second() - secondsNow)
					fps := float64(currentFPS) / float64(deltaSeconds)
					if fps >= 1 {
						log.Printf("%.0f FPS\n", fps)
					} else {
						log.Printf("%.2f FPS\n", fps)
					}
				}
				currentFPS = 0
			}

			lastUpdate = time.Now()
		}
	}()

}

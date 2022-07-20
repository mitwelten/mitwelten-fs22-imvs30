package aggregator

import (
	"log"
	"mjpeg_multiplexer/src/global"
	"mjpeg_multiplexer/src/input"
	"mjpeg_multiplexer/src/mjpeg"
	"reflect"
	"time"
)

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
	passthrough       bool
	Enabled           bool
	AggregatorStorage *mjpeg.FrameStorage
}

// StartAggregator starts the aggregator loop for the passed aggregator in a separate go routine
func StartAggregator(agg *Aggregator, inputs ...input.Input) {
	var inputStorages []*mjpeg.FrameStorage
	for _, input := range inputs {
		inputStorages = append(inputStorages, input.GetInputData().InputStorage)
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
	} else if !global.DecodingNecessary() && len(inputStorages) == 1 {
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
			if !global.Config.AlwaysActive && !aggregator.GetAggregatorData().Enabled {
				time.Sleep(1 * time.Second)
				continue
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
			if lastUpdate.Second() != time.Now().Second() {
				if global.Config.LogFPS {
					log.Printf("%v FPS\n", currentFPS)
				}
				currentFPS = 0
			}

			lastUpdate = time.Now()
		}
	}()

}

package aggregator

import (
	"fmt"
	"log"
	"mjpeg_multiplexer/src/global"
	"mjpeg_multiplexer/src/mjpeg"
	"reflect"
	"sync"
	"time"
)

// Aggregator aggregates multiple frame storages to one frame storage
// takes multiple frameStorages (one frame storage for each input Connection) and
// aggregates them to one frameStorage
type Aggregator interface {
	aggregate(channels ...*mjpeg.FrameStorage) *mjpeg.MjpegFrame
	Setup(storages ...*mjpeg.FrameStorage)
	GetAggregatorData() *AggregatorData
}

type AggregatorData struct {
	passthrough     bool
	Enabled         bool
	OutputStorage   *mjpeg.FrameStorage
	OutputCondition *sync.Cond
}

func Aggregate(aggregatorPtr *Aggregator, storages ...*mjpeg.FrameStorage) {
	aggregator := *aggregatorPtr
	aggregator.Setup(storages...)
	aggregatorData := aggregator.GetAggregatorData()
	aggregatorData.OutputStorage = mjpeg.NewFrameStorage()
	condition := setupCondition(storages...)

	lastUpdate := time.Unix(0, 0)
	FPS := 0

	passthroughMode := false

	if !global.DecodingNecessary() && reflect.TypeOf(aggregator) == reflect.TypeOf(&AggregatorCarousel{}) {
		passthroughMode = true
	} else if !global.DecodingNecessary() && len(storages) == 1 {
		passthroughMode = true
	}

	if passthroughMode {
		fmt.Printf("Passthorugh-Mode activate - Each frame will be directly passed to the output without any decoding and encoding\n")
	} else {
		fmt.Printf("Passthorugh-Mode not active - Each frame will be decoded and encoded, which may be slow.\n=> Activate it by either using the 'carousel' mode or by only using one input source and by removing the width, height, quality and show_label options.\n")
	}

	go func() {
		for {
			condition.Wait()

			//todo wait with a condition here?
			if !aggregator.GetAggregatorData().Enabled {
				time.Sleep(1 * time.Second)
				continue
			}

			if global.Config.OutputFramerate != -1 && time.Since(lastUpdate).Seconds() < (1.0/global.Config.OutputFramerate) {
				continue
			}

			// get frame
			var frame *mjpeg.MjpegFrame
			if aggregatorData.passthrough && len(storages) == 1 && passthroughMode {
				frame = storages[0].GetLatestPtr()
			} else {
				frame = aggregator.aggregate(storages...)
			}

			if frame == nil {
				continue
			}

			if aggregatorData.OutputStorage != nil {
				aggregatorData.OutputStorage.StorePtr(frame)
				aggregatorData.OutputCondition.Signal()
			}

			FPS++
			if lastUpdate.Second() != time.Now().Second() {
				if global.Config.LogFPS {
					log.Printf("%v FPS\n", FPS)
				}
				FPS = 0
			}

			lastUpdate = time.Now()
		}
	}()

}

func setupCondition(storages ...*mjpeg.FrameStorage) *sync.Cond {
	// init the lock and condition object to notify the aggregator when a new frame has been stored
	lock := sync.Mutex{}
	lock.Lock()
	condition := sync.NewCond(&lock)
	for _, el := range storages {
		el.AggregatorCondition = condition
	}

	return condition
}

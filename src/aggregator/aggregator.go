package aggregator

import (
	"log"
	"mjpeg_multiplexer/src/global"
	"mjpeg_multiplexer/src/mjpeg"
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

	go func() {
		for {
			condition.Wait()

			if global.Config.OutputFramerate != -1 && time.Since(lastUpdate).Seconds() < (1.0/global.Config.OutputFramerate) {
				continue
			}

			// start time
			var s time.Time
			if global.Config.LogTime {
				s = time.Now()
			}

			// get frame
			var frame *mjpeg.MjpegFrame
			if aggregatorData.passthrough && len(storages) == 1 {
				frame = storages[0].GetLatestPtr()
			} else {
				frame = aggregator.aggregate(storages...)
			}

			// stop time
			if global.Config.LogTime {
				log.Printf("Aggregate with %v images: %vms\n", len(storages), time.Since(s).Milliseconds())
			}

			outputCondition := aggregatorData.OutputStorage
			if outputCondition != nil {
				aggregatorData.OutputStorage.StorePtr(frame)
				aggregatorData.OutputCondition.Signal()
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

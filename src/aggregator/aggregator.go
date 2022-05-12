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
	aggregatorData := aggregator.GetAggregatorData()
	aggregatorData.OutputStorage = mjpeg.NewFrameStorage()
	condition := setupCondition(storages...)
	go func() {
		for {
			condition.Wait()

			var frame *mjpeg.MjpegFrame
			if aggregatorData.passthrough && len(storages) == 1 {
				frame = storages[0].GetLatestPtr()
			} else {
				var s time.Time
				if global.Config.LogTime {
					s = time.Now()
				}

				frame = aggregator.aggregate(storages...)

				if global.Config.LogTime {
					log.Printf("Aggregate with %v images: %vms\n", len(storages), time.Since(s).Milliseconds())
				}
			}

			outputCondition := aggregatorData.OutputStorage
			if outputCondition != nil {
				aggregatorData.OutputStorage.StorePtr(frame)
				aggregatorData.OutputCondition.Signal()
			}
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

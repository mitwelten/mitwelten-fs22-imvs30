package aggregator

import (
	"log"
	"mjpeg_multiplexer/src/communication"
	"mjpeg_multiplexer/src/global"
	"mjpeg_multiplexer/src/imageUtils"
	"mjpeg_multiplexer/src/mjpeg"
	"sync"
	"time"
)

type AggregatorGrid struct {
	Row             int
	Col             int
	OutputStorage   *communication.FrameStorage
	OutputCondition *sync.Cond
}

func (aggregator *AggregatorGrid) SetOutputCondition(cond *sync.Cond) {
	aggregator.OutputCondition = cond
}
func (aggregator *AggregatorGrid) GetStorage() *communication.FrameStorage {
	return aggregator.OutputStorage
}

func (aggregator *AggregatorGrid) Aggregate(storages ...*communication.FrameStorage) {
	aggregator.OutputStorage = communication.NewFrameStorage()

	// init the lock and condition object to notify the aggregator when a new frame has been stored
	lock := sync.Mutex{}
	lock.Lock()
	condition := sync.NewCond(&lock)
	for _, el := range storages {
		el.AggregatorCondition = condition
	}

	go func() {
		for {
			condition.Wait()

			var frames []*mjpeg.MjpegFrame

			for i := 0; i < len(storages); i++ {
				frame := storages[i]
				frames = append(frames, frame.GetLatestPtr())
			}

			//todo remove me
			doPassthrough := true
			var frame *mjpeg.MjpegFrame
			//do pass through on only 1 source
			if len(storages) == 1 && doPassthrough {
				frame = frames[0]
			} else {

				var s time.Time
				if global.Config.LogTime {
					s = time.Now()
				}

				frame = imageUtils.Grid(aggregator.Row, aggregator.Col, frames...)

				if global.Config.LogTime {
					log.Printf("Grid with %v images: %vms\n", len(frames), time.Since(s).Milliseconds())
				}
			}

			if aggregator.OutputCondition != nil {
				aggregator.OutputStorage.StorePtr(frame)
				aggregator.OutputCondition.Signal()
			}
		}
	}()
}

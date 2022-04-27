package aggregator

import (
	"mjpeg_multiplexer/src/communication"
	"mjpeg_multiplexer/src/image"
	"mjpeg_multiplexer/src/mjpeg"
	"sync"
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

			var frames []mjpeg.MjpegFrame

			for i := 0; i < len(storages); i++ {
				frame := storages[i]
				frames = append(frames, frame.GetLatest())
			}

			var frame mjpeg.MjpegFrame
			if len(storages) == 1 {
				frame = frames[0]
			} else {
				frame = image.Grid(aggregator.Row, aggregator.Col, frames...)
			}

			if aggregator.OutputCondition != nil {
				aggregator.OutputStorage.Store(frame)
				aggregator.OutputCondition.Signal()
			}
		}
	}()
}

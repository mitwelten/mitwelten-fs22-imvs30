package aggregator

import (
	"mjpeg_multiplexer/src/communication"
	"mjpeg_multiplexer/src/image"
	"mjpeg_multiplexer/src/mjpeg"
	"sync"
)

type AggregatorGrid struct {
	Row int
	Col int
}

func (grid AggregatorGrid) Aggregate(storages ...*communication.FrameStorage) *communication.FrameStorage {
	storage := communication.FrameStorage{}

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

			frame := image.Grid(grid.Row, grid.Col, frames...)
			storage.Store(frame)
		}
	}()

	return &storage
}

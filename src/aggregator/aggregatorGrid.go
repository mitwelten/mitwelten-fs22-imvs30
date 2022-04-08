package aggregator

import (
	"mjpeg_multiplexer/src/communication"
	"mjpeg_multiplexer/src/image"
	"mjpeg_multiplexer/src/mjpeg"
)

type AggregatorGrid struct {
	Row int
	Col int
}

func (grid AggregatorGrid) Aggregate(storages ...*communication.FrameStorage) *communication.FrameStorage {
	storage := communication.FrameStorage{}
	go func() {
		for {
			var frames []mjpeg.MjpegFrame
			for i := 0; i < len(storages); i++ {
				frame := storages[i]
				frames = append(frames, frame.Get())
			}

			//			start := time.Now()

			frame := image.Grid(grid.Row, grid.Col, frames...)
			storage.Store(frame)
			/*
				t := time.Now()
				elapsed := t.Sub(start)
				log.Println(elapsed.Milliseconds(), "ms for image merging grid")
			*/
		}
	}()

	return &storage
}

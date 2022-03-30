package aggregator

import (
	"log"
	"mjpeg_multiplexer/src/communication"
	"mjpeg_multiplexer/src/image"
	"mjpeg_multiplexer/src/mjpeg"
	"time"
)

type AggregatorGrid struct {
	Row int
	Col int
}

func (grid AggregatorGrid) Aggregate(channels ...*communication.FrameData) chan mjpeg.Frame {
	var out = make(chan mjpeg.Frame)
	go func() {
		for {
			var frames []mjpeg.Frame
			for i := 0; i < len(channels); i++ {
				frame := channels[i]
				frames = append(frames, frame.Get())
			}

			start := time.Now()
			out <- image.Grid(grid.Row, grid.Col, frames...)
			t := time.Now()
			elapsed := t.Sub(start)
			log.Println(elapsed.Milliseconds(), "ms for image merging grid")
		}
	}()

	return out
}

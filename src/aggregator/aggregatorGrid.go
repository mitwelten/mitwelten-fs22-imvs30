package aggregator

import (
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
			var n = false
			for i := 0; i < len(channels); i++ {
				frame := channels[i]
				if frame.Get().Body == nil {
					n = true
				}
				frames = append(frames, frame.Get())
			}

			if n == true {
				continue
			}

			start := time.Now()
			out <- image.Grid(grid.Row, grid.Col, frames...)
			t := time.Now()
			elapsed := t.Sub(start)
			println(elapsed.Milliseconds(), "ms for image merging grid")
		}
	}()

	return out
}

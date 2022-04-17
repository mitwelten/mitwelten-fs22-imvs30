package aggregator

import (
	"mjpeg_multiplexer/src/communication"
	"mjpeg_multiplexer/src/image"
	"mjpeg_multiplexer/src/mjpeg"
	"reflect"
)

type AggregatorGrid struct {
	Row int
	Col int
}

func (grid AggregatorGrid) Aggregate(storages ...*communication.FrameStorage) *communication.FrameStorage {
	storage := communication.FrameStorage{}
	go func() {
		var prev mjpeg.MjpegFrame
		prev = mjpeg.MjpegFrame{Body: mjpeg.Init()}
		for {
			var frames []mjpeg.MjpegFrame

			for i := 0; i < len(storages); i++ {
				frame := storages[i]

				if reflect.DeepEqual(frame.Get(), prev.Body) {
					continue
				}

				frames = append(frames, frame.Get())
				//s := time.Now()
				//imgDiff := image.GetImgDiff(frame.Get(), prev)
				imgDiff := image.GetImgDiffResize(frame.Get(), prev)
				//println(time.Since(s).Milliseconds())

				frames = append(frames, imgDiff)
				prev = frame.Get()
				/*				prev = mjpeg.MjpegFrame{}
								copy(prev.Body, frame.Get().Body)
				*/
			}

			frame := image.Grid(grid.Row, grid.Col, frames...)
			storage.Store(frame)
		}
	}()

	return &storage
}

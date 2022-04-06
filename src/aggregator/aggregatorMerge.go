package aggregator

import (
	"mjpeg_multiplexer/src/image"
	"mjpeg_multiplexer/src/mjpeg"
	"time"
)

type AggregatorMerge struct {
}

func (merge AggregatorMerge) MergeImages(channels ...chan mjpeg.MjpegFrame) chan mjpeg.MjpegFrame {
	var combine = AggregatorCombine{}
	var channel = combine.Aggregate(channels...)
	var out = make(chan mjpeg.MjpegFrame)
	go func(channel_ chan mjpeg.MjpegFrame) {
		for {
			var f1 = <-channel_
			var f2 = <-channel_
			start := time.Now()
			out <- image.MergeImages(f1, f2)
			t := time.Now()
			elapsed := t.Sub(start)
			println(elapsed.Milliseconds(), "ms for image merging")
		}
	}(channel)

	return out
}

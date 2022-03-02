package aggregator

import (
	"mjpeg_multiplexer/src/image"
	"mjpeg_multiplexer/src/mjpeg"
	"time"
)

func CombineChannels(channels []chan mjpeg.Frame) chan mjpeg.Frame {
	agg := make(chan mjpeg.Frame)
	for _, channel := range channels {
		go func(channel_ chan mjpeg.Frame) {
			for msg := range channel_ {
				agg <- msg
			}
		}(channel)
	}

	return agg
}

func MergeImages(channels []chan mjpeg.Frame) chan mjpeg.Frame {
	var channel = CombineChannels(channels)
	var out = make(chan mjpeg.Frame)
	go func(channel_ chan mjpeg.Frame) {
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

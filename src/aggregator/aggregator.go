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
func MergeImagesGrid4(channels []chan mjpeg.Frame) chan mjpeg.Frame {
	var channel = CombineChannels(channels)
	var out = make(chan mjpeg.Frame)
	go func(channel_ chan mjpeg.Frame) {
		for {
			var f1 = <-channel_
			var f2 = <-channel_
			var f3 = <-channel_
			var f4 = <-channel_
			start := time.Now()
			out <- image.Grid4(f1, f2, f3, f4)
			t := time.Now()
			elapsed := t.Sub(start)
			println(elapsed.Milliseconds(), "ms for image merging grid4")
		}
	}(channel)

	return out
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

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Merge2Images(channels []chan mjpeg.Frame) chan mjpeg.Frame {
	if len(channels) != 2 {
		panic("Must specify 2 input channels")
	}
	var out = make(chan mjpeg.Frame)
	go func(c1 chan mjpeg.Frame, c2 chan mjpeg.Frame) {
		for {
			start := time.Now()
			var f1 = <-c1
			var f2 = <-c2
			t := time.Now()
			elapsed := t.Sub(start)
			println(elapsed.Milliseconds(), "ms for image aggregation")

			start = time.Now()
			out <- image.MergeImages(f1, f2)
			t = time.Now()
			elapsed = t.Sub(start)
			println(elapsed.Milliseconds(), "ms for image merging")
		}
	}(channels[0], channels[1])

	return out
}

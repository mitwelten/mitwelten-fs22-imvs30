package connection

import "mjpeg_multiplexer/src/mjpeg"

type Sink interface {
	ProcessFrame(frame mjpeg.Frame)
}

func RunSink(sink Sink, channel chan mjpeg.Frame) {
	go func(agg chan mjpeg.Frame) {
		for {
			frame := <-channel
			sink.ProcessFrame(frame)
		}
	}(channel)
}

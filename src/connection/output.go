package connection

import "mjpeg_multiplexer/src/mjpeg"

type Output interface {
	ProcessFrame(frame mjpeg.Frame)
}

func RunOutput(sink Output, channel chan mjpeg.Frame) {
	go func(agg chan mjpeg.Frame) {
		for {
			frame := <-channel
			sink.ProcessFrame(frame)
		}
	}(channel)
}

package connection

import "mjpeg_multiplexer/src/mjpeg"

type Output interface {
	SendFrame(frame mjpeg.Frame) error
}

func RunOutput(sink Output, channel chan mjpeg.Frame) {
	go func(agg chan mjpeg.Frame) {
		for {
			frame := <-channel
			err := sink.SendFrame(frame)
			if err != nil {
				println("Error while trying to send frame to output")
				println(err.Error())
				continue
			}
		}
	}(channel)
}

package connection

import (
	"mjpeg_multiplexer/src/communication"
	"mjpeg_multiplexer/src/mjpeg"
)

type Output interface {
	SendFrame(frame mjpeg.MjpegFrame) error
}

func RunOutput(sink Output, channel chan mjpeg.MjpegFrame) {
	go func(agg chan mjpeg.MjpegFrame) {
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

func (output OutputHTTP) run(storage *communication.FrameStorage) {
	go func(storage_ *communication.FrameStorage) {
		for {
			frame := storage_.Get()
			err := output.SendFrame(frame)
			if err != nil {
				println("Error while trying to send frame to output")
				println(err.Error())
				continue
			}
		}
	}(storage)
}

func (output OutputFile) run(storage *communication.FrameStorage) {
	go func(storage_ *communication.FrameStorage) {
		for {
			frame := storage_.Get()
			err := output.SendFrame(frame)
			if err != nil {
				println("Error while trying to send frame to output")
				println(err.Error())
				continue
			}
		}
	}(storage)
}

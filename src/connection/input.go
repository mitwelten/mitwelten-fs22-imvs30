package connection

import (
	"mjpeg_multiplexer/src/communication"
	"mjpeg_multiplexer/src/mjpeg"
)

type Input interface {
	ReceiveFrame() (mjpeg.Frame, error)
}

func ListenToInput(input Input) *communication.FrameData {
	var frame = mjpeg.Frame{}
	frame.Body = mjpeg.Init()

	var frameData = communication.FrameData{}
	frameData.Store(frame) // init with a black frame

	go func() {
		for {
			var frame, err = input.ReceiveFrame()
			if err != nil {
				println("Error while trying to read frame from input")
				println(err.Error())
				continue
			}
			frameData.Store(frame)
		}
	}()
	return &frameData
}

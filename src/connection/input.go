package connection

import (
	"mjpeg_multiplexer/src/communication"
	"mjpeg_multiplexer/src/mjpeg"
)

type Input interface {
	ReceiveFrame() (mjpeg.MjpegFrame, error)
	Start()
}

func ListenToInput(input Input) *communication.FrameStorage {
	var frame = mjpeg.MjpegFrame{}
	frame.Body = mjpeg.Init()

	var frameData = communication.FrameStorage{}
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

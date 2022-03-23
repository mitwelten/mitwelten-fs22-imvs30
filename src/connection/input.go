package connection

import (
	args "mjpeg_multiplexer/src/communication"
	"mjpeg_multiplexer/src/mjpeg"
)

type Input interface {
	ReceiveFrame() (mjpeg.Frame, error)
}

func ListenToInput(input Input) *args.FrameData {
	var frameData = args.FrameData{}
	frameData.Init() // init with a black frame
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

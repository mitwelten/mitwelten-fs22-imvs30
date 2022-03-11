package connection

import "mjpeg_multiplexer/src/mjpeg"

type Input interface {
	ReceiveFrame() (mjpeg.Frame, error)
}

func ListenToInput(input Input) chan mjpeg.Frame {
	var channel = make(chan mjpeg.Frame)
	go func() {
		for {
			var frame, err = input.ReceiveFrame()
			if err != nil {
				println("Error while trying to read frame from input")
				println(err.Error())
				continue
			}

			//channel <- frame

			//Skip current frame if channel is not being read
			select {
			case channel <- frame:
			default: //skip frame
			}
		}
	}()
	return channel
}
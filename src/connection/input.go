package connection

import (
	"log"
	"mjpeg_multiplexer/src/communication"
	"mjpeg_multiplexer/src/mjpeg"
	"time"
)

type Input interface {
	ReceiveFrame() (mjpeg.MjpegFrame, error)
	Start() error
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
				log.Println("error " + err.Error())

				// TODO proper error handling if input is no longer reachable
				time.Sleep(1 * time.Second)
				continue
			}
			frameData.Store(frame)
		}
	}()
	return &frameData
}

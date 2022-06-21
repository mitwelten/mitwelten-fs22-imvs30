package connection

import (
	"errors"
	"log"
	"mjpeg_multiplexer/src/customErrors"
	"mjpeg_multiplexer/src/mjpeg"
	"time"
)

type Input interface {
	ReceiveFrame() (mjpeg.MjpegFrame, error)
	ReceiveFrameFast() (mjpeg.MjpegFrame, error)
	Start() error
	Info() string
}

func reconnectInput(input Input) {
	for {
		err := input.Start()
		if err == nil {
			log.Printf("Successfully reconnected to %s\n", input.Info())
			return
		}

		log.Printf("Could not reconnect to %s\n", input.Info())
		time.Sleep(5 * time.Minute)
	}

}

func ListenToInput(input Input) *mjpeg.FrameStorage {

	frameData := mjpeg.NewFrameStorage()

	go func() {
		for {
			var frame, err = input.ReceiveFrameFast()

			if errors.As(err, &customErrors.ErrInvalidFrame{}) {
				// retry when receiving invalid frame
				log.Printf("Invalid frame read from socket %s: %s\n", input.Info(), err.Error())
				continue
			} else if err != nil {
				// reconnect on read error
				log.Printf("Could not read from socket %s: %s\n", input.Info(), err.Error())
				reconnectInput(input)
				continue
			}
			frameData.Store(frame)
		}
	}()
	return frameData
}

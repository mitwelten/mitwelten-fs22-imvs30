package connection

import (
	"log"
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
		log.Printf("Retrying to connect to %s...\n", input.Info())
		err := input.Start()
		if err == nil {
			log.Printf("Successfully reconnected to %s\n", input.Info())
			return
		}

		log.Printf("Could not reconnect to %s\n", input.Info())
		time.Sleep(1 * time.Minute)
	}

}

func ListenToInput(input Input) *mjpeg.FrameStorage {

	frameData := mjpeg.NewFrameStorage()

	go func() {
		for {
			var frame, err = input.ReceiveFrameFast()
			if err != nil {
				log.Printf("error %s\n", err.Error())
				reconnectInput(input)
				continue
			}
			frameData.Store(frame)
		}
	}()
	return frameData
}

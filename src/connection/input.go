package connection

import (
	"log"
	"mjpeg_multiplexer/src/global"
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

	lastReceive := time.Unix(0, 0)
	go func() {
		for {
			if time.Since(lastReceive).Seconds() < (1.0 / global.Config.InputFramerate) {
				//todo optimize: This always reads and parses a full frame. Is a more efficient alternative possible and useful?
				_, _ = input.ReceiveFrameFast()
				continue
			}
			lastReceive = time.Now()
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

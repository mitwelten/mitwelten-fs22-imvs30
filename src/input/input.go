package input

import (
	"errors"
	"log"
	"mjpeg_multiplexer/src/customErrors"
	"mjpeg_multiplexer/src/imageUtils"
	"mjpeg_multiplexer/src/mjpeg"
	"time"
)

type Input interface {
	ReceiveFrame() (mjpeg.MjpegFrame, error)
	Init() error
	Info() string
	GetInputData() *InputData
}

type InputData struct {
	InputStorage *mjpeg.FrameStorage
}

func reconnectInput(input Input) {
	for {
		err := input.Init()
		if err == nil {
			log.Printf("Successfully reconnected to %s\n", input.Info())
			return
		}

		log.Printf("Could not reconnect to %s\n", input.Info())
		time.Sleep(5 * time.Minute)
	}

}

//StartInput starts the input source by calling the Init() method and running ReceiveFrame() in a loop
func StartInput(input Input) {

	inputData := input.GetInputData()
	inputData.InputStorage = mjpeg.NewFrameStorage()

	go func() {
		err := input.Init()
		if err != nil {
			log.Fatalf("Can't open input stream: %s", err.Error())
		}
		frame, err := input.ReceiveFrame()
		// store and encode the first frame to get information about its size
		if err == nil {
			inputData.InputStorage.Store(&frame)
			imageUtils.Decode(inputData.InputStorage)
		}

		for {
			var frame, err = input.ReceiveFrame()

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
			inputData.InputStorage.Store(&frame)
		}
	}()
}

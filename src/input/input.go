package input

import (
	"errors"
	"log"
	"mjpeg_multiplexer/src/customErrors"
	"mjpeg_multiplexer/src/mjpeg"
	"time"
)

type Input interface {
	//ReceiveFrame receives a single frame from the input
	ReceiveFrame() (mjpeg.MjpegFrame, error)
	//Init starts the connection to the input
	Init() error
	//GetInfo returns a description to of the connection
	GetInfo() string
	//GetInputData returns all other input data of the input
	GetInputData() *InputData
}

type InputData struct {
	InputStorage *mjpeg.FrameStorage
}

func reconnectInput(input Input) {
	for {
		err := input.Init()
		if err == nil {
			log.Printf("Successfully reconnected to %s\n", input.GetInfo())
			return
		}

		log.Printf("Could not reconnect to %s\n", input.GetInfo())
		time.Sleep(5 * time.Minute)
	}

}

//StartInput starts the input source by calling the Init() method and running ReceiveFrame() in a loop
func StartInput(input Input) {
	inputData := input.GetInputData()

	go func() {
		err := input.Init()
		if err != nil {
			log.Fatalf("Can't open input stream for %v: %s", input.GetInfo(), err.Error())
		}

		for {
			var frame, err = input.ReceiveFrame()

			if errors.As(err, &customErrors.ErrInvalidFrame{}) {
				// retry when receiving invalid frame
				log.Printf("Invalid frame read from socket %s: %s\n", input.GetInfo(), err.Error())
				continue
			} else if err != nil {
				// reconnect on read error
				log.Printf("Could not read from socket %s: %s\n", input.GetInfo(), err.Error())
				reconnectInput(input)
				continue
			}
			inputData.InputStorage.Store(&frame)
		}
	}()
}

package connection

import (
	"errors"
	"mjpeg_multiplexer/src/communication"
	"mjpeg_multiplexer/src/mjpeg"
	"os"
)

type OutputFile struct {
	filePath string
}

func NewOutputFile(filePath string) (sink OutputFile) {
	return OutputFile{filePath}
}

// SendFrame todo TEST: check if file has been created and matches
func (output OutputFile) SendFrame(frame mjpeg.MjpegFrame) error {
	fh, err := os.OpenFile(output.filePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return errors.New("cannot write to file")
	}
	_, err = fh.Write(frame.Body)
	if err != nil {
		return errors.New("cannot write to file")
	}
	return nil
}

func (output OutputFile) Run(storage *communication.FrameStorage) {
	go func(storage_ *communication.FrameStorage) {
		for {
			frame := storage_.Get()
			err := output.SendFrame(frame)
			if err != nil {
				println("Error while trying to send frame to output")
				println(err.Error())
				continue
			}
		}
	}(storage)
}

package connection

import (
	"errors"
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

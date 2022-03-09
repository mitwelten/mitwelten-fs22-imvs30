package connection

import (
	"mjpeg_multiplexer/src/mjpeg"
	"os"
)

type OutputFile struct {
	filePath string
}

func NewOutputFile(filePath string) (sink OutputFile) {
	return OutputFile{filePath}
}

func (sink OutputFile) ProcessFrame(frame mjpeg.Frame) {
	fh, err := os.OpenFile(sink.filePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		panic("Can't write file")
	}
	_, err = fh.Write(frame.Body)
	if err != nil {
		return
	}
}

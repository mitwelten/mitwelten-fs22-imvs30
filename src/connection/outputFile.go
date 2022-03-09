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
	err := os.WriteFile(sink.filePath, frame.Body, 0644)
	if err != nil {
		panic("Can't write file")
	}
}

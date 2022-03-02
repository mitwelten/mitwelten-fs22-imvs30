package connection

import (
	"mjpeg_multiplexer/src/mjpeg"
	"os"
)

type FileSink struct {
	filePath string
}

func NewFileSink(filePath string) (sink FileSink) {
	return FileSink{filePath}
}

func (sink FileSink) ProcessFrame(frame mjpeg.Frame) {
	os.WriteFile(sink.filePath, frame.Body, 0644)
}

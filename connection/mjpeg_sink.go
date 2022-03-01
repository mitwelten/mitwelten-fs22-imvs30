package connection

import "mjpeg"

import (
	"os"
)

type Sink struct {
	filePath string
}

func NewSink(filePath string) (sink Sink) {
	return Sink{filePath}
}

func (sink Sink) Run(sources []chan mjpeg.Frame) {
	agg := make(chan mjpeg.Frame)
	for _, ch := range sources {
		go func(c chan mjpeg.Frame) {
			for msg := range c {
				agg <- msg
			}
		}(ch)
	}

	go func(agg chan mjpeg.Frame) {
		for {
			frame := <-agg
			os.WriteFile(sink.filePath, frame.Body, 0644)
		}
	}(agg)
}

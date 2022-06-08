package connection

import (
	"log"
	"mjpeg_multiplexer/src/aggregator"
	"mjpeg_multiplexer/src/mjpeg"
	"os"
)

type OutputFile struct {
	filePath   string
	aggregator aggregator.Aggregator
}

func NewOutputFile(filePath string, aggregator aggregator.Aggregator) Output {
	return &OutputFile{filePath, aggregator}
}

// SendFrame todo TEST: check if file has been created and matches
func (output *OutputFile) SendFrame(frame *mjpeg.MjpegFrame) {
	fh, err := os.OpenFile(output.filePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Println("error: cannot write to file")
		//return &customErrors.ErrIOWrite{}
	}
	_, err = fh.Write(frame.Body)
	if err != nil {
		log.Println("error: cannot write to file")
		//return &customErrors.ErrIOWrite{}
	}
	//return nil
}

func (output *OutputFile) Run() {
	go func(storage_ *mjpeg.FrameStorage) {
		for {
			frame := storage_.GetLatestPtr()
			output.SendFrame(frame)
		}
	}(output.aggregator.GetAggregatorData().OutputStorage)
}

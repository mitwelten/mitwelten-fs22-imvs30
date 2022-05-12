package connection

import (
	"log"
	"mjpeg_multiplexer/src/aggregator"
	"mjpeg_multiplexer/src/customErrors"
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
		log.Println("error: cannot write to file")
		return &customErrors.ErrIOWrite{}
	}
	_, err = fh.Write(frame.Body)
	if err != nil {
		log.Println("error: cannot write to file")
		return &customErrors.ErrIOWrite{}
	}
	return nil
}

func (output OutputFile) Run(aggregator aggregator.Aggregator) {
	go func(storage_ *mjpeg.FrameStorage) {
		for {
			frame := storage_.GetLatest()
			err := output.SendFrame(frame)
			if err != nil {
				log.Println("error: while trying to send frame to output")
				println(err.Error())
				continue
			}
		}
	}(aggregator.GetAggregatorData().OutputStorage)
}

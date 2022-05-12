package aggregator

import (
	"log"
	"mjpeg_multiplexer/src/imageUtils"
	"mjpeg_multiplexer/src/mjpeg"
)

type AggregatorGrid struct {
	data AggregatorData
	Row  int
	Col  int
}

func (aggregator *AggregatorGrid) Setup(storages ...*mjpeg.FrameStorage) {
	//ensure that enough space is available:
	var nCells = aggregator.Row * aggregator.Col
	var nFrames = len(storages)
	if nFrames > nCells {
		log.Fatalf("Too many frames for this grid configuartion: row %v col %v, but %v frames to compute\n", aggregator.Row, aggregator.Col, nFrames)
	}

	aggregator.data.passthrough = true
}

func (aggregator *AggregatorGrid) GetAggregatorData() *AggregatorData {
	return &aggregator.data
}

func (aggregator *AggregatorGrid) aggregate(storages ...*mjpeg.FrameStorage) *mjpeg.MjpegFrame {
	var frames []*mjpeg.MjpegFrame

	for i := 0; i < len(storages); i++ {
		frame := storages[i]
		frames = append(frames, frame.GetLatestPtr())
	}

	return imageUtils.Grid(aggregator.Row, aggregator.Col, frames...)
}

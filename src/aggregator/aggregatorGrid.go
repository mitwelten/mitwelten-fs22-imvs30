package aggregator

import (
	"log"
	"mjpeg_multiplexer/src/imageUtils"
	"mjpeg_multiplexer/src/mjpeg"
)

// AggregatorGrid shows all images on the screen, filling them up from top left to bottom right.
// Empty spaces will remain black.
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
	return imageUtils.Grid(aggregator.Row, aggregator.Col, storages...)
}

package aggregator

import (
	"mjpeg_multiplexer/src/imageUtils"
	"mjpeg_multiplexer/src/mjpeg"
)

type AggregatorGrid struct {
	data AggregatorData
	Row  int
	Col  int
}

func (aggregator *AggregatorGrid) Setup(_ ...*mjpeg.FrameStorage) {
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

	frame := imageUtils.Grid(aggregator.Row, aggregator.Col, frames...)

	return frame
}

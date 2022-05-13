package aggregator

import (
	"mjpeg_multiplexer/src/mjpeg"
	"time"
)

type AggregatorCarousel struct {
	data         AggregatorData
	Duration     time.Duration
	lastSwitch   time.Time
	currentIndex int
}

func (aggregator *AggregatorCarousel) Setup(_ ...*mjpeg.FrameStorage) {
	aggregator.data.passthrough = true
	aggregator.lastSwitch = time.Now()
	aggregator.currentIndex = 0
}

func (aggregator *AggregatorCarousel) GetAggregatorData() *AggregatorData {
	return &aggregator.data
}

func (aggregator *AggregatorCarousel) aggregate(storages ...*mjpeg.FrameStorage) *mjpeg.MjpegFrame {
	//todo what about images with different resolutions? Downscale to the smallest? Use the OutputMaxWidth config?

	if time.Since(aggregator.lastSwitch) > aggregator.Duration {
		aggregator.currentIndex = (aggregator.currentIndex + 1) % len(storages)
		aggregator.lastSwitch = time.Now()
	}
	return storages[aggregator.currentIndex].GetLatestPtr()
}

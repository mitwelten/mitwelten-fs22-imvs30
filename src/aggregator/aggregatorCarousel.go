package aggregator

import (
	"log"
	"mjpeg_multiplexer/src/global"
	"mjpeg_multiplexer/src/imageUtils"
	"mjpeg_multiplexer/src/mjpeg"
	"mjpeg_multiplexer/src/motionDetection"
	"time"
)

const minWaitBetweenChanges = 2000 * time.Millisecond

type AggregatorCarousel struct {
	data           AggregatorData
	Duration       time.Duration
	lastSwitch     time.Time
	currentIndex   int
	motionDetector *motionDetection.MotionDetector
}

func (aggregator *AggregatorCarousel) Setup(storages ...*mjpeg.FrameStorage) {
	aggregator.data.passthrough = false
	aggregator.lastSwitch = time.Now()
	aggregator.currentIndex = 0
	if global.Config.UseMotion {
		aggregator.motionDetector = motionDetection.NewMotionDetector(storages...)
	}
}

func (aggregator *AggregatorCarousel) GetAggregatorData() *AggregatorData {
	return &aggregator.data
}

func (aggregator *AggregatorCarousel) aggregate(storages ...*mjpeg.FrameStorage) *mjpeg.MjpegFrame {
	//todo what about images with different resolutions? Downscale to the smallest? Use the OutputMaxWidth config?
	newIndex := aggregator.motionDetector.GetMostActiveIndex()
	if newIndex == -1 && time.Since(aggregator.lastSwitch) >= aggregator.Duration {
		// duration update
		aggregator.currentIndex = (aggregator.currentIndex + 1) % len(storages)
		log.Printf("Switched due to DURATION\n")
		aggregator.lastSwitch = time.Now()
	} else if newIndex != -1 && time.Since(aggregator.lastSwitch) >= minWaitBetweenChanges {
		//  motion update
		aggregator.currentIndex = newIndex
		log.Printf("Switched due to MOTION\n")
		aggregator.lastSwitch = time.Now()
	}

	return imageUtils.Transform(storages[aggregator.currentIndex])
	//return imageUtils.Encode(aggregator.motionDetector.GetMostActiveImage())
}

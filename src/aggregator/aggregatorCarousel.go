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
	previousIndex  int
	motionDetector *motionDetection.MotionDetector

	lastFrame       *mjpeg.MjpegFrame
	lastFrameUpdate time.Time
}

func (aggregator *AggregatorCarousel) Setup(storages ...*mjpeg.FrameStorage) {
	if global.DecodingNecessary() {
		log.Printf("Carousel: Each frame will have to be decoded and encoded, which may be slow. Remove the width, height, quality and show_label paramter to activate the passthrough mode for more frames per second.\n")
	} else {
		log.Printf("Carousel: Passthrough mode activated - each frame will be directly passed to the output without any decoding and encoding\n")
	}

	aggregator.data.passthrough = false
	aggregator.lastSwitch = time.Now()
	aggregator.currentIndex = 0
	aggregator.previousIndex = 0

	aggregator.lastFrameUpdate = time.Now()

	if global.Config.UseMotion {
		aggregator.motionDetector = motionDetection.NewMotionDetector(storages...)
	}
}

func (aggregator *AggregatorCarousel) GetAggregatorData() *AggregatorData {
	return &aggregator.data
}

func (aggregator *AggregatorCarousel) aggregate(storages ...*mjpeg.FrameStorage) *mjpeg.MjpegFrame {
	index := -1
	if aggregator.motionDetector != nil {
		index = aggregator.motionDetector.GetMostActiveIndex()
	}
	if index == -1 && time.Since(aggregator.lastSwitch) >= aggregator.Duration {
		// duration update
		aggregator.currentIndex = (aggregator.currentIndex + 1) % len(storages)
		aggregator.lastSwitch = time.Now()
	} else if index != -1 && time.Since(aggregator.lastSwitch) >= minWaitBetweenChanges {
		//  motion update
		aggregator.currentIndex = index
		aggregator.lastSwitch = time.Now()
	}

	if aggregator.previousIndex == aggregator.currentIndex && storages[aggregator.currentIndex].LastUpdated == aggregator.lastFrameUpdate {
		var frame *mjpeg.MjpegFrame = nil
		// chrome bug: Because the stream lags 1 frame behind, we resend the last frame before stopping
		// link: https://bugs.chromium.org/p/chromium/issues/detail?id=527446
		if aggregator.lastFrame != nil {
			frame = aggregator.lastFrame
			aggregator.lastFrame = nil
		}
		return frame
	}

	aggregator.previousIndex = aggregator.currentIndex
	aggregator.lastFrameUpdate = storages[aggregator.currentIndex].LastUpdated
	// save the last frame to resend it later on
	frame := imageUtils.Carousel(storages[aggregator.currentIndex], aggregator.currentIndex)
	aggregator.lastFrame = frame
	return frame
}

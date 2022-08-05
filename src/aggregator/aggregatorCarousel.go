package aggregator

import (
	"mjpeg_multiplexer/src/activityDetection"
	"mjpeg_multiplexer/src/global"
	"mjpeg_multiplexer/src/imageUtils"
	"mjpeg_multiplexer/src/mjpeg"
	"time"
)

//minimum time used to prevent changes too quickly, thus confusing the user
const minWaitBetweenChanges = 3000 * time.Millisecond

//AggregatorCarousel shows one images of the other.
//It switches between these images after the by the --duration parameter specified time.
//Can be used with --activity
type AggregatorCarousel struct {
	data                    AggregatorData
	Duration                time.Duration
	lastSwitch              time.Time
	lastMotionInActiveFrame time.Time
	CurrentIndex            int
	previousIndex           int
	activityDetector        *activityDetection.ActivityDetector

	lastFrame       *mjpeg.MjpegFrame
	lastFrameUpdate time.Time
}

func (aggregator *AggregatorCarousel) Setup(storages ...*mjpeg.FrameStorage) {
	if aggregator.Duration <= 0 {
		aggregator.Duration = defaultDuration
	}

	aggregator.data.passthrough = false
	aggregator.lastSwitch = time.Now()
	aggregator.lastMotionInActiveFrame = time.Now()
	aggregator.CurrentIndex = 0
	aggregator.previousIndex = 0

	aggregator.lastFrameUpdate = time.Now()

	if global.Config.UseActivity {
		aggregator.activityDetector = activityDetection.NewActivityDetector(storages...)
	}
}

func (aggregator *AggregatorCarousel) GetAggregatorData() *AggregatorData {
	return &aggregator.data
}

func (aggregator *AggregatorCarousel) aggregate(storages ...*mjpeg.FrameStorage) *mjpeg.MjpegFrame {
	index := -1
	if aggregator.activityDetector != nil {
		index = aggregator.activityDetector.GetMostActiveIndex()
	}
	if index == -1 && time.Since(aggregator.lastSwitch) >= aggregator.Duration && time.Since(aggregator.lastMotionInActiveFrame) >= minWaitBetweenChanges {
		// duration update
		aggregator.CurrentIndex = (aggregator.CurrentIndex + 1) % len(storages)
		aggregator.lastSwitch = time.Now()
	} else if index != -1 && index != aggregator.CurrentIndex && time.Since(aggregator.lastSwitch) >= minWaitBetweenChanges {
		//  motion update
		aggregator.CurrentIndex = index
		aggregator.lastSwitch = time.Now()
	} else if index != -1 && index == aggregator.CurrentIndex {
		//motion in the same frame
		aggregator.lastMotionInActiveFrame = time.Now()
	}

	if aggregator.previousIndex == aggregator.CurrentIndex && storages[aggregator.CurrentIndex].LastUpdated == aggregator.lastFrameUpdate {
		var frame *mjpeg.MjpegFrame = nil
		// chrome bug: Because the stream lags 1 frame behind, we resend the last frame before stopping
		// link: https://bugs.chromium.org/p/chromium/issues/detail?id=527446
		if aggregator.lastFrame != nil {
			frame = aggregator.lastFrame
			aggregator.lastFrame = nil
		}
		return frame
	}

	aggregator.previousIndex = aggregator.CurrentIndex
	aggregator.lastFrameUpdate = storages[aggregator.CurrentIndex].LastUpdated
	// save the last frame to resend it later on
	frame := imageUtils.Carousel(storages[aggregator.CurrentIndex], aggregator.CurrentIndex)
	aggregator.lastFrame = frame
	return frame
}

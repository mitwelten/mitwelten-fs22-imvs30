package aggregator

import (
	"mjpeg_multiplexer/src/global"
	"mjpeg_multiplexer/src/imageUtils"
	"mjpeg_multiplexer/src/mjpeg"
	"mjpeg_multiplexer/src/motionDetection"
	"time"
)

type AggregatorPanel struct {
	data           AggregatorData
	layout         imageUtils.PanelLayout
	CycleFrames    bool
	Duration       time.Duration
	lastSwitch     time.Time
	currentIndex   int
	motionDetector *motionDetection.MotionDetector
}

func (aggregator *AggregatorPanel) Setup(storages ...*mjpeg.FrameStorage) {
	aggregator.data.passthrough = false //only valid for 3 or more inputs
	aggregator.lastSwitch = time.Now()
	aggregator.currentIndex = 0

	nStorages := len(storages)
	//todo input validation in parser
	if nStorages > 6 {
		aggregator.layout = imageUtils.Slots8
	} else if nStorages > 4 {
		aggregator.layout = imageUtils.Slots6
	} else if nStorages > 3 {
		aggregator.layout = imageUtils.Slots4
	} else {
		aggregator.layout = imageUtils.Slots3
	}

	if global.Config.UseMotion {
		aggregator.motionDetector = motionDetection.NewMotionDetector(storages...)
	}
}

func (aggregator *AggregatorPanel) GetAggregatorData() *AggregatorData {
	return &aggregator.data
}

func (aggregator *AggregatorPanel) aggregate(storages ...*mjpeg.FrameStorage) *mjpeg.MjpegFrame {
	newIndex := aggregator.motionDetector.GetMostActiveIndex()
	if newIndex == -1 && time.Since(aggregator.lastSwitch) >= aggregator.Duration {
		// duration update
		aggregator.currentIndex = (aggregator.currentIndex + 1) % len(storages)
		aggregator.lastSwitch = time.Now()
	} else if newIndex != -1 && time.Since(aggregator.lastSwitch) >= minWaitBetweenChanges {
		//  motion update
		aggregator.currentIndex = newIndex
		aggregator.lastSwitch = time.Now()
	}

	return imageUtils.Panel(aggregator.layout, aggregator.currentIndex, storages...)
}

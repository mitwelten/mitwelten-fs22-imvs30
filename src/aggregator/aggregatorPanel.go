package aggregator

import (
	"mjpeg_multiplexer/src/activityDetection"
	"mjpeg_multiplexer/src/global"
	"mjpeg_multiplexer/src/imageUtils"
	"mjpeg_multiplexer/src/mjpeg"
	"time"
)

//AggregatorPanel shows all images on the screen in predefined layouts.
//A fitting layout will be selected on the Setup call
//--panel_cycle: Will switch the slots position over time
//--activity: move most active panel into the panel top left
//see also aggregatorCarousel
type AggregatorPanel struct {
	data                    AggregatorData
	layout                  imageUtils.PanelLayout
	CycleFrames             bool
	Duration                time.Duration
	lastSwitch              time.Time
	lastMotionInActiveFrame time.Time
	currentIndex            int
	activityDetector        *activityDetection.ActivityDetector
}

func (aggregator *AggregatorPanel) Setup(storages ...*mjpeg.FrameStorage) {
	if aggregator.Duration <= 0 {
		aggregator.Duration = defaultDuration
	}

	aggregator.data.passthrough = false
	aggregator.lastSwitch = time.Now()
	aggregator.lastMotionInActiveFrame = time.Now()
	aggregator.currentIndex = 0

	nStorages := len(storages)
	if nStorages > 6 {
		aggregator.layout = imageUtils.Slots8
	} else if nStorages > 4 {
		aggregator.layout = imageUtils.Slots6
	} else if nStorages > 3 {
		aggregator.layout = imageUtils.Slots4
	} else {
		aggregator.layout = imageUtils.Slots3
	}

	if global.Config.UseActivity {
		aggregator.activityDetector = activityDetection.NewActivityDetector(storages...)
	}
}

func (aggregator *AggregatorPanel) GetAggregatorData() *AggregatorData {
	return &aggregator.data
}

func (aggregator *AggregatorPanel) aggregate(storages ...*mjpeg.FrameStorage) *mjpeg.MjpegFrame {
	index := -1
	if aggregator.activityDetector != nil {
		index = aggregator.activityDetector.GetMostActiveIndex()
	}

	if aggregator.CycleFrames && index == -1 && time.Since(aggregator.lastSwitch) >= aggregator.Duration && time.Since(aggregator.lastMotionInActiveFrame) >= minWaitBetweenChanges {
		// duration update
		aggregator.currentIndex = (aggregator.currentIndex + 1) % len(storages)
		aggregator.lastSwitch = time.Now()
	} else if index != -1 && index != aggregator.currentIndex && time.Since(aggregator.lastSwitch) >= minWaitBetweenChanges {
		//motion update
		aggregator.currentIndex = index
		aggregator.lastSwitch = time.Now()
	} else if index != -1 && index == aggregator.currentIndex {
		//motion in the same frame
		aggregator.lastMotionInActiveFrame = time.Now()

	}

	//return imageUtils.Encode(activityDetection.FrameDifferenceImage(imageUtils.DecodeAt(storages[0], 0), imageUtils.DecodeAt(storages[0], 1)))
	return imageUtils.Panel(aggregator.layout, aggregator.currentIndex, storages...)
}

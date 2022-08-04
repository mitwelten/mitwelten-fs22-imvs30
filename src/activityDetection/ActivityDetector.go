package activityDetection

import (
	"fmt"
	"image"
	"mjpeg_multiplexer/src/global"
	"mjpeg_multiplexer/src/imageUtils"
	"mjpeg_multiplexer/src/mjpeg"
	"mjpeg_multiplexer/src/utils"
	"time"
)

// ActivityDetector aggregates multiple frame storages, calculates score to detect most attractive input and sets output
type ActivityDetector struct {
	storages        []*mjpeg.FrameStorage
	previousScores  []utils.RingBuffer[float64]
	previousFrames  []image.Image
	lastScoreUpdate time.Time
}

const updateDelay = 1000 * time.Millisecond
const minScore = 0.05
const nPreviousScores = 5

//NewActivityDetector creates a new instances and allocates all needed memory
func NewActivityDetector(storages ...*mjpeg.FrameStorage) *ActivityDetector {
	activityDetector := ActivityDetector{}
	activityDetector.storages = storages
	activityDetector.previousScores = make([]utils.RingBuffer[float64], len(storages))
	activityDetector.previousFrames = make([]image.Image, len(storages))

	for i := range storages {
		// init the buffers for the average change scores
		activityDetector.previousScores[i] = utils.NewRingBuffer[float64](nPreviousScores)
	}

	activityDetector.lastScoreUpdate = time.Now()

	return &activityDetector
}

//UpdateScores will update the scores if enough time has passed since the last update (more than 'updateDelay')
func (activityDetector *ActivityDetector) UpdateScores() {
	if time.Since(activityDetector.lastScoreUpdate) < updateDelay {
		return
	}

	for i, storage := range activityDetector.storages {
		currentFrame := imageUtils.Decode(storage)
		ratio := float64(currentFrame.Bounds().Dy()) / float64(currentFrame.Bounds().Dx())
		height := int(imageSize * ratio)
		currentFrame = imageUtils.Resize(currentFrame, imageSize, height)

		previousFrame := activityDetector.previousFrames[i]

		if previousFrame != nil {
			score := FrameDifferenceScore(currentFrame, previousFrame)
			activityDetector.previousScores[i].Push(score)
		}
		activityDetector.previousFrames[i] = currentFrame
	}

	activityDetector.lastScoreUpdate = time.Now()
}

// GetMostActiveIndex calls `UpdateScores` and then returns the index of the most active stream
// A small threshold has to be reached for a stream to be considered active, if no stream is considered active `-1` is returned
func (activityDetector *ActivityDetector) GetMostActiveIndex() int {
	activityDetector.UpdateScores()

	scores := make([]float64, len(activityDetector.previousScores))
	for i, el := range activityDetector.previousScores {
		data, size := el.GetData()
		if size == 0 {
			continue
		}
		scores[i] = average(*data, size)

		if global.Config.Debug {
			global.Config.InputConfigs[i].Label = fmt.Sprintf("%.3f", scores[i])
		}
	}

	index, score := argmax(scores)

	if score >= minScore {

		if global.Config.Debug {
			global.Config.InputConfigs[index].Label = fmt.Sprintf("!%.3f", scores[index])
		}

		return index
	} else {
		return -1
	}
}

//average takes the average of a float array, only considering the first `size` elements
func average(arr []float64, size int) float64 {
	sum := 0.0
	for i := 0; i < size; i++ {
		sum += arr[i]
	}
	return sum / float64(size)

}

// argmax returns index of max int value in given array
func argmax(data []float64) (int, float64) {
	max := data[0]
	maxIndex := 0

	for i := 1; i < len(data); i++ {
		if data[i] > max {
			max = data[i]
			maxIndex = i
		}
	}
	return maxIndex, max
}

//GetMostActiveImage todo REMOVE
func (activityDetector *ActivityDetector) GetMostActiveImage() image.Image {
	index := activityDetector.GetMostActiveIndex()
	if index == -1 {
		return mjpeg.NewMJPEGFrame().CachedImage
	}

	return FrameDifferenceImage(imageUtils.Decode(activityDetector.storages[index]), imageUtils.Decode(activityDetector.storages[index]))

}

package motionDetection

import (
	"fmt"
	"image"
	"mjpeg_multiplexer/src/global"
	"mjpeg_multiplexer/src/imageUtils"
	"mjpeg_multiplexer/src/mjpeg"
	"mjpeg_multiplexer/src/utils"
	"time"
)

// MotionDetector aggregates multiple frame storages, calculates score to detect most attractive input and sets output
type MotionDetector struct {
	storages        []*mjpeg.FrameStorage
	previousScores  []utils.RingBuffer[float64]
	previousFrames  []image.Image
	lastScoreUpdate time.Time
}

const updateDelay = 750 * time.Millisecond
const minScore = 0.001
const nPreviousScores = 3

func NewMotionDetector(storages ...*mjpeg.FrameStorage) *MotionDetector {
	motionDetector := MotionDetector{}
	motionDetector.storages = storages
	motionDetector.previousScores = make([]utils.RingBuffer[float64], len(storages))
	motionDetector.previousFrames = make([]image.Image, len(storages))

	for i := range storages {
		// init the buffers for the average change scores
		motionDetector.previousScores[i] = utils.NewRingBuffer[float64](nPreviousScores)
	}

	motionDetector.lastScoreUpdate = time.Now()

	return &motionDetector
}

func (motionDetector *MotionDetector) UpdateScores() {
	if time.Since(motionDetector.lastScoreUpdate) < updateDelay {
		return
	}

	for i, storage := range motionDetector.storages {
		currentFrame := imageUtils.Decode(storage)
		ratio := float64(currentFrame.Bounds().Dy()) / float64(currentFrame.Bounds().Dx())
		height := int(imageSize * ratio)
		currentFrame = imageUtils.Resize(currentFrame, imageSize, height)

		previousFrame := motionDetector.previousFrames[i]

		if previousFrame != nil {
			score := FrameDifferenceScore(currentFrame, previousFrame)
			motionDetector.previousScores[i].Push(score)
		}
		motionDetector.previousFrames[i] = currentFrame
	}

	motionDetector.lastScoreUpdate = time.Now()
}

func (motionDetector *MotionDetector) GetMostActiveIndex() int {
	motionDetector.UpdateScores()

	scores := make([]float64, len(motionDetector.previousScores))
	for i, el := range motionDetector.previousScores {
		data, size := el.GetData()
		if size == 0 {
			continue
		}
		scores[i] = averageScore(*data, size)

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

func averageScore(arr []float64, size int) float64 {
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

func (motionDetector *MotionDetector) GetMostActiveImage() image.Image {

	index := motionDetector.GetMostActiveIndex()
	if index == -1 {
		return mjpeg.NewMJPEGFrame().CachedImage
	}

	return FrameDifferenceImage(imageUtils.DecodeAt(motionDetector.storages[index], 0), imageUtils.DecodeAt(motionDetector.storages[index], 4))

}

package motionDetection

import (
	"fmt"
	"image"
	"mjpeg_multiplexer/src/mjpeg"
	"mjpeg_multiplexer/src/utils"
	"time"
)

// MotionDetection aggregates multiple frame storages, calculates score to detect most attractive input and sets output
type MotionDetection struct {
	storages        []*mjpeg.FrameStorage
	previousImages  []utils.RingBuffer[image.Image]
	previousScores  []utils.RingBuffer[float64]
	lastScoreUpdate time.Time
}

const updateDelay = 1000 * time.Millisecond
const nPreviousScores = 5

func (motionDetector *MotionDetection) Setup(storages ...*mjpeg.FrameStorage) {
	motionDetector.storages = storages
	motionDetector.previousImages = make([]utils.RingBuffer[image.Image], len(storages))
	motionDetector.previousScores = make([]utils.RingBuffer[float64], len(storages))

	for i, _ := range storages {
		// init the buffers for the average change scores
		motionDetector.previousImages[i] = utils.NewRingBuffer[image.Image](nPreviousScores)
		motionDetector.previousScores[i] = utils.NewRingBuffer[float64](nPreviousScores)
	}

	motionDetector.lastScoreUpdate = time.Now()

}

func (motionDetector *MotionDetection) UpdateScores(images ...image.Image) {
	if time.Since(motionDetector.lastScoreUpdate) < updateDelay {
		return
	}

	for i, img := range images {
		if len(motionDetector.previousImages[i].GetAllPtr()) != 0 {
			score := FrameDifferenceScore(motionDetector.previousImages[i].Peek(), img)
			motionDetector.previousScores[i].Push(score)
		}
		motionDetector.previousImages[i].Push(img)
	}

	motionDetector.lastScoreUpdate = time.Now()
}

func (motionDetector *MotionDetection) GetMostActiveIndex() int {
	scores := make([]float64, len(motionDetector.previousScores))
	for i, el := range motionDetector.previousScores {
		data, size := el.GetData()
		if size == 0 {
			continue
		}
		scores[i] = averageScore(*data, size)
		fmt.Printf("Frame %v scores: %v\n", i, data)
	}

	index, _ := argmax(scores)
	return index
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

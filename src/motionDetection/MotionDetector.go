package motionDetection

import (
	"image"
	"mjpeg_multiplexer/src/imageUtils"
	"mjpeg_multiplexer/src/mjpeg"
	"mjpeg_multiplexer/src/utils"
	"time"
)

// MotionDetector aggregates multiple frame storages, calculates score to detect most attractive input and sets output
type MotionDetector struct {
	storages        []*mjpeg.FrameStorage
	previousScores  []utils.RingBuffer[float64]
	lastScoreUpdate time.Time
}

const updateDelay = 1000 * time.Millisecond
const minScore = 0.001
const nPreviousScores = 3

func NewMotionDetector(storages ...*mjpeg.FrameStorage) *MotionDetector {
	motionDetector := MotionDetector{}
	motionDetector.storages = storages
	motionDetector.previousScores = make([]utils.RingBuffer[float64], len(storages))

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
		if len(storage.GetAllPtr()) >= 2 {
			//todo cache last frame to reuse it next second?
			score := FrameDifferenceScore(imageUtils.Decode(storage.GetAllPtr()[0]), imageUtils.Decode(storage.GetAllPtr()[1]))
			motionDetector.previousScores[i].Push(score)
		}
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
		//fmt.Printf("Frame %v scores: %v\n", i, scores[i])
	}

	index, score := argmax(scores)
	//fmt.Printf("Best score: %v from index %v\n", score, index)

	if score >= minScore {
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

	return FrameDifferenceImage(imageUtils.Decode(motionDetector.storages[index].GetAllPtr()[0]), imageUtils.Decode(motionDetector.storages[index].GetAllPtr()[4]))

}

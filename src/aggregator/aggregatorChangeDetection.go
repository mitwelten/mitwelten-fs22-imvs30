package aggregator

import (
	"fmt"
	"log"
	"mjpeg_multiplexer/src/changeDetection"
	"mjpeg_multiplexer/src/global"
	"mjpeg_multiplexer/src/mjpeg"
	"mjpeg_multiplexer/src/utils"
	"time"
)

// AggregatorChange aggregates multiple frame storages, calculates score to detect most attractive input and sets output
type AggregatorChange struct {
	data              AggregatorData
	scorer            changeDetection.PixelDifferenceScorer
	previousScores    []utils.RingBuffer[float64]
	frameStorageIndex int
	lastScoreUpdate   time.Time
	lastFrameUpdate   time.Time
}

const delay = 2000 * time.Millisecond
const nPreviousScores = 20

func (aggregator *AggregatorChange) Setup(storages ...*mjpeg.FrameStorage) {
	aggregator.data.passthrough = true
	aggregator.scorer = changeDetection.PixelDifferenceScorer{}
	aggregator.previousScores = make([]utils.RingBuffer[float64], len(storages))

	for i, _ := range storages {
		// init the buffers for the average change scores
		aggregator.previousScores[i] = utils.NewRingBuffer[float64](nPreviousScores)
	}

	aggregator.frameStorageIndex = -1

	aggregator.lastScoreUpdate = time.Now()
	aggregator.lastFrameUpdate = time.Now()

}

func (aggregator *AggregatorChange) GetAggregatorData() *AggregatorData {
	return &aggregator.data
}

func (aggregator *AggregatorChange) aggregate(FrameStorages ...*mjpeg.FrameStorage) *mjpeg.MjpegFrame {
	// calculate the new scores if the storage was updated
	for i := 0; i < len(FrameStorages); i++ {
		frameStorage := FrameStorages[i]
		if frameStorage.LastUpdated.Before(aggregator.lastFrameUpdate) {
			continue
		}

		var s time.Time
		if global.Config.LogTime {
			s = time.Now()
		}

		score := aggregator.scorer.Score(frameStorage.GetAllPtr())
		aggregator.previousScores[i].Push(score)

		if global.Config.LogTime {
			log.Printf("AggregatorChangeDetection (index %v): %vms\n", i, time.Since(s).Milliseconds())
		}
	}

	// change the index to the newest active frame
	if aggregator.frameStorageIndex == -1 || time.Since(aggregator.lastScoreUpdate) > delay {
		scores := make([]float64, len(aggregator.previousScores))
		for i, el := range aggregator.previousScores {
			data, size := el.GetData()
			if size == 0 {
				continue
			}
			scores[i] = averageScore(*data, size)
			fmt.Printf("Frame %v scores: %v\n", i, data)
		}

		aggregator.frameStorageIndex, _ = argmax(scores)
		fmt.Printf("Biggest score index is %d\n", aggregator.frameStorageIndex)
		aggregator.lastScoreUpdate = time.Now()
	}

	return FrameStorages[aggregator.frameStorageIndex].GetLatestPtr()
	//return aggregator.scorer.Diff(FrameStorages[aggregator.frameStorageIndex].GetAll()))
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

package aggregator

import (
	"fmt"
	"log"
	"mjpeg_multiplexer/src/changeDetection"
	"mjpeg_multiplexer/src/communication"
	"mjpeg_multiplexer/src/global"
	"mjpeg_multiplexer/src/utils"
	"sync"
	"time"
)

// AggregatorChange aggregates multiple frame storages, calculates score to detect most attractive input and sets output
type AggregatorChange struct {
	OutputCondition *sync.Cond
	OutputStorage   *communication.FrameStorage
}

const delay = 2000 * time.Millisecond
const nPreviousScores = 20

func (aggregator *AggregatorChange) SetOutputCondition(cond *sync.Cond) {
	aggregator.OutputCondition = cond
}

func (aggregator *AggregatorChange) GetStorage() *communication.FrameStorage {
	return aggregator.OutputStorage
}

func (aggregator *AggregatorChange) Aggregate(FrameStorages ...*communication.FrameStorage) {
	aggregator.OutputStorage = communication.NewFrameStorage()
	scorer := changeDetection.PixelDifferenceScorer{}

	// init the lock and condition object to notify the aggregator when a new frame has been stored
	lock := sync.Mutex{}
	lock.Lock()
	condition := sync.NewCond(&lock)

	// buffer for the average change values - one for each storage
	previousScores := make([]utils.RingBuffer[float64], len(FrameStorages))

	for i, el := range FrameStorages {
		// set the condition for the wakeup signal
		el.AggregatorCondition = condition

		// init the buffers for the average change scores
		previousScores[i] = utils.NewRingBuffer[float64](nPreviousScores)
	}

	frameStorageIndex := -1

	lastScoreUpdate := time.Now()
	lastFrameUpdate := time.Now()

	go func() {
		for {
			condition.Wait()

			// calculate the new scores if the storage was updated
			for i := 0; i < len(FrameStorages); i++ {
				frameStorage := FrameStorages[i]
				if frameStorage.LastUpdated.Before(lastFrameUpdate) {
					continue
				}

				var s time.Time
				if global.Config.LogTime {
					s = time.Now()
				}

				score := scorer.Score(frameStorage.GetAllPtr())
				previousScores[i].Push(score)

				if global.Config.LogTime {
					log.Printf("AggregatorChangeDetection (index %v): %vms\n", i, time.Since(s).Milliseconds())
				}
			}

			// change the index to the newest active frame
			if frameStorageIndex == -1 || time.Since(lastScoreUpdate) > delay {
				scores := make([]float64, len(previousScores))
				for i, el := range previousScores {
					data, size := el.GetData()
					if size == 0 {
						continue
					}
					scores[i] = averageScore(*data, size)
					fmt.Printf("Frame %v scores: %v\n", i, data)
				}

				frameStorageIndex, _ = argmax(scores)
				fmt.Printf("Biggest score index is %d\n", frameStorageIndex)
				lastScoreUpdate = time.Now()
			}

			if aggregator.OutputCondition != nil {
				//fmt.Printf("update index is %d\n", frameStorageIndex)
				aggregator.OutputStorage.Store(FrameStorages[frameStorageIndex].GetLatest())
				//s := time.Now()
				//aggregator.OutputStorage.Store(scorer.Diff(FrameStorages[frameStorageIndex].GetAll()))
				//fmt.Printf("%v\n", time.Since(s))

				/*				aggregator.OutputStorage.Store(scorer.Diff(FrameStorages[frameStorageIndex].GetAll()))
								fmt.Printf("%v\n", time.Since(s))
				*/
				aggregator.OutputCondition.Signal()
			}
			lastFrameUpdate = time.Now()
		}
	}()
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

package aggregator

import (
	"fmt"
	"mjpeg_multiplexer/src/changeDetection"
	"mjpeg_multiplexer/src/communication"
	"sync"
	"time"
)

// AggregatorChange aggregates multiple frame storages, calculates score to detect most attractive input and sets output
type AggregatorChange struct {
	OutputCondition *sync.Cond
}

const delay = 5000 * time.Millisecond
const averageSize = 20

func (aggregator AggregatorChange) Aggregate(storages ...*communication.FrameStorage) *communication.FrameStorage {
	storage := communication.FrameStorage{}

	// init the lock and condition object to notify the aggregator when a new frame has been stored
	lock := sync.Mutex{}
	lock.Lock()
	condition := sync.NewCond(&lock)
	for _, el := range storages {
		el.AggregatorCondition = condition
	}

	scorer := changeDetection.PixelDifferenceScorer{}

	lastIndex := -1
	lastTimestamp := time.Now()

	go func() {
		for {
			condition.Wait()

			if lastIndex == -1 || time.Since(lastTimestamp) > delay {

				scores := make([]int, len(storages))

				for i := 0; i < len(storages); i++ {
					frame := storages[i]
					score := scorer.Score(frame.GetAll())
					scores[i] = score
				}
				fmt.Printf("%v\n", scores)
				index, _ := argmax(scores)

				lastIndex = index
				lastTimestamp = time.Now()

				storage.Store(storages[index].GetLatest())
			} else {
				storage.Store(storages[lastIndex].GetLatest())
			}

			if aggregator.OutputCondition != nil {
				aggregator.OutputCondition.Signal()
			}
		}
	}()

	return &storage
}

// argmax returns index of max int value in given array
func argmax(data []int) (int, int) {
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

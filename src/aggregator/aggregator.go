package aggregator

import (
	"mjpeg_multiplexer/src/communication"
)

type Aggregator interface {
	Aggregate(channels ...*communication.FrameStorage) *communication.FrameStorage
}

package aggregator

import (
	"mjpeg_multiplexer/src/communication"
)

// Aggregator aggregates multiple frame storages to one frame storage
// takes multiple frameStorages (one frame storage for each input Connection) and
// aggregates them to one frameStorage
type Aggregator interface {
	Aggregate(channels ...*communication.FrameStorage) *communication.FrameStorage
}

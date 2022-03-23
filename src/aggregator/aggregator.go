package aggregator

import (
	"mjpeg_multiplexer/src/communication"
	"mjpeg_multiplexer/src/mjpeg"
)

type Aggregator interface {
	Aggregate(channels ...*communication.FrameData) chan mjpeg.Frame
}

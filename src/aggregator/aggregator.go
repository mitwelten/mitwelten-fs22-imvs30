package aggregator

import (
	"mjpeg_multiplexer/src/mjpeg"
)

type Aggregator interface {
	Aggregate(channels ...chan mjpeg.Frame) chan mjpeg.Frame
}

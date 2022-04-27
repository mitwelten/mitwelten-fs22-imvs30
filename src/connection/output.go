package connection

import (
	"mjpeg_multiplexer/src/aggregator"
	"mjpeg_multiplexer/src/mjpeg"
)

type Output interface {
	SendFrame(frame mjpeg.MjpegFrame) error
	Run(aggregator aggregator.Aggregator)
}

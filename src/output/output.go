package output

import (
	"mjpeg_multiplexer/src/aggregator"
	"mjpeg_multiplexer/src/mjpeg"
)

type Output interface {
	//SendFrame Sends a single frame to the output
	SendFrame(frame *mjpeg.MjpegFrame)
	//StartOutput Starts the Output
	StartOutput(agg *aggregator.Aggregator)
}

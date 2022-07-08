package connection

import (
	"mjpeg_multiplexer/src/aggregator"
	"mjpeg_multiplexer/src/mjpeg"
	"mjpeg_multiplexer/src/utils"
	"testing"
	"time"
)

func TestInAndOutput(t *testing.T) {

	var agg aggregator.Aggregator = &aggregator.AggregatorCarousel{}
	output := NewOutputHTTP("8123", agg)
	input := NewInputHTTP(nil, "localhost:8123")
	storage := ListenToInput(input)

	aggregator.Aggregate(&agg, storage)
	output.Run()

	frame := mjpeg.MjpegFrame{}
	frame.Body = []byte{1, 2, 3}

	utils.Assert(t, frame.Body, storage.GetLatestPtr().Body, false)
	storage.Store(frame)

	time.Sleep(200 * time.Millisecond)

	utils.Assert(t, frame.Body, storage.GetLatestPtr().Body, true)
}

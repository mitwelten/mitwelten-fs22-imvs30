package connection

import (
	"mjpeg_multiplexer/src/aggregator"
	input2 "mjpeg_multiplexer/src/input"
	"mjpeg_multiplexer/src/mjpeg"
	output2 "mjpeg_multiplexer/src/output"
	"mjpeg_multiplexer/src/utils"
	"testing"
	"time"
)

func TestInAndOutput(t *testing.T) {

	var agg aggregator.Aggregator = &aggregator.AggregatorCarousel{}
	output := output2.NewOutputHTTP("8123")
	input := input2.NewInputHTTP(nil, "localhost:8123")
	input2.StartInput(input)

	aggregator.StartAggregator(&agg, input)
	output.StartOutput(&agg)

	frame := mjpeg.MjpegFrame{}
	frame.Body = []byte{1, 2, 3}

	utils.Assert(t, frame.Body, input.GetInputData().InputStorage.GetFrame().Body, false)
	input.GetInputData().InputStorage.Store(&frame)

	time.Sleep(200 * time.Millisecond)

	utils.Assert(t, frame.Body, input.GetInputData().InputStorage.GetFrame().Body, true)
}

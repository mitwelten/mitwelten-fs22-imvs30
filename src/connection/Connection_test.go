package connection

import (
	_ "embed"
	"mjpeg_multiplexer/src/aggregator"
	"mjpeg_multiplexer/src/global"
	"mjpeg_multiplexer/src/input"
	"mjpeg_multiplexer/src/mjpeg"
	"mjpeg_multiplexer/src/output"
	"mjpeg_multiplexer/src/utils"
	"testing"
	"time"
)

//go:embed img1.jpg
var img []byte

//test layout:
//outTest --(:8120)--> multiplexer --(:8121)--> inTest
func TestInToOutput(t *testing.T) {
	global.SetupInitialConfig()

	//outTest: dummy output to send frames
	dummyIn := input.NewInputHTTP(0, "")
	var dummyAgg aggregator.Aggregator = &aggregator.AggregatorCarousel{}
	outTest := output.NewOutputHTTP("8120")
	aggregator.StartAggregator(&dummyAgg, dummyIn)
	outTest.StartOutput(&dummyAgg)

	//multiplexer:
	in := input.NewInputHTTP(1, "localhost:8120")
	var agg aggregator.Aggregator = &aggregator.AggregatorCarousel{}
	out := output.NewOutputHTTP("8121")
	input.StartInput(in)
	aggregator.StartAggregator(&agg, in)
	out.StartOutput(&agg)

	//inTest: listen to 8121 - output of multiplexer
	inTest := input.NewInputHTTP(2, "localhost:8121")
	input.StartInput(inTest)
	time.Sleep(250 * time.Millisecond)

	//setup a dummy frame
	frame := mjpeg.MjpegFrame{}
	frame.Body = img

	//inTest should not have received the dummy frame yet
	utils.Assert(t, frame.Body, inTest.GetInputData().InputStorage.GetFrame().Body, false)
	utils.Assert(t, mjpeg.NewMJPEGFrame().Body, inTest.GetInputData().InputStorage.GetFrame().Body, true)

	//send the dummy frame through to the multiplexer via the input...
	outTest.SendFrame(&frame)

	time.Sleep(250 * time.Millisecond)

	//inTest should have received the dummy frame as output from the multiplexer
	utils.Assert(t, frame.Body, inTest.GetInputData().InputStorage.GetFrame().Body, true)
}

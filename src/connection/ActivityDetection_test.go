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
var img1 []byte

//go:embed img2.jpg
var img2 []byte

//test layout:
// outTest1 --(:8220)-->
//                       multiplexer --(:8221)--> inTest
// outTest2 --(:8222)-->

func TestActivation(t *testing.T) {
	global.SetupInitialConfig()
	global.Config.UseMotion = true

	//outTest: dummy output to send frames
	dummyIn := input.NewInputHTTP(0, "")
	var dummyAgg aggregator.Aggregator = &aggregator.AggregatorCarousel{}

	outTest1 := output.NewOutputHTTP("8220")
	outTest2 := output.NewOutputHTTP("8222")
	aggregator.StartAggregator(&dummyAgg, dummyIn)
	outTest1.StartOutput(&dummyAgg)
	outTest2.StartOutput(&dummyAgg)

	//multiplexer:
	in1 := input.NewInputHTTP(1, "localhost:8220")
	in2 := input.NewInputHTTP(1, "localhost:8222")
	var agg aggregator.Aggregator = &aggregator.AggregatorCarousel{}
	out := output.NewOutputHTTP("8221")
	input.StartInput(in1)
	input.StartInput(in2)
	aggregator.StartAggregator(&agg, in1, in2)
	out.StartOutput(&agg)

	//inTest: listen to 8121 - output of multiplexer
	inTest := input.NewInputHTTP(2, "localhost:8221")
	input.StartInput(inTest)
	time.Sleep(250 * time.Millisecond)

	//setup a dummy frames
	frame1 := mjpeg.MjpegFrame{}
	frame1.Body = img1

	frame2 := mjpeg.MjpegFrame{}
	frame2.Body = img2

	//inTest should not have received the dummy frame yet
	utils.Assert(t, frame1.Body, inTest.GetInputData().InputStorage.GetFrame().Body, false)
	utils.Assert(t, frame2.Body, inTest.GetInputData().InputStorage.GetFrame().Body, false)
	utils.Assert(t, mjpeg.NewMJPEGFrame().Body, inTest.GetInputData().InputStorage.GetFrame().Body, true)

	agg_ := agg.(*aggregator.AggregatorCarousel)
	//activity on stream 1:
	//5sec
	for i := 0; i < 5; i++ {
		outTest1.SendFrame(&frame1)
		outTest2.SendFrame(&frame1)
		time.Sleep(1000 * time.Millisecond)

		outTest1.SendFrame(&frame2)
		outTest2.SendFrame(&frame1)
		time.Sleep(1000 * time.Millisecond)
	}

	utils.Assert(t, agg_.CurrentIndex, 0, true)

	//activity on stream 2:
	//5sec
	for i := 0; i < 100; i++ {
		outTest1.SendFrame(&frame2)
		outTest2.SendFrame(&frame2)
		time.Sleep(25 * time.Millisecond)

		outTest1.SendFrame(&frame2)
		outTest2.SendFrame(&frame1)
		time.Sleep(25 * time.Millisecond)
	}
	utils.Assert(t, agg_.CurrentIndex, 1, true)
}

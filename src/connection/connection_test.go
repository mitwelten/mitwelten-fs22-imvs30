package connection

import (
	"bytes"
	"io/ioutil"
	"mjpeg_multiplexer/src/mjpeg"
	"os"
	"testing"
)

//todo change hardcoded file path to something more robust?
var imageLocation = "../../resources/images1/image1.jpg"

func mockFrame() mjpeg.Frame {
	fh, err := os.OpenFile(imageLocation, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		panic("cant find file " + imageLocation)
	}
	b, err := ioutil.ReadAll(fh)
	if err != nil {
		panic("cant read from file " + imageLocation)
	}

	err = fh.Close()
	if err != nil {
		panic("cant close file handler for " + imageLocation)
	}

	return mjpeg.Frame{Body: b}

}

func TestSendAndReceive(t *testing.T) {
	output, err := NewOutputHTTP("8081")
	if err != nil {
		panic("Can't open output")
	}

	input := NewInputHTTP("localhost:8081")

	frame := mockFrame()
	err = output.SendFrame(frame)
	if err != nil {
		panic("Can't send frame")
	}

	frameOut, err := input.ReceiveFrame()
	if err != nil {
		panic("Can't receive frame")
	}

	if bytes.Compare(frame.Body, frameOut.Body) != 0 {
		t.Errorf("Frames do not match")
	}
}
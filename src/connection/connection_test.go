package connection

import (
	"bytes"
	"io/ioutil"
	"mjpeg_multiplexer/src/mjpeg"
	"os"
	"testing"
	"time"
)

//todo change hardcoded file path to something more robust?
var imageLocation = "../../resources/images1/image1.jpg"

func mockFrame() mjpeg.MjpegFrame {
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

	return mjpeg.MjpegFrame{Body: b}
}

func TestSendAndReceive(t *testing.T) {
	go func() {
		time.Sleep(2 * time.Second)
		t.Error("frame not received within time limit")
		os.Exit(1)
	}()

	output, err := NewOutputHTTP("8081")
	if err != nil {
		panic("Can't open output")
	}

	input := NewInputHTTP("localhost:8081")

	// Wait a short amount to make sure that input and output are ready
	time.Sleep(250 * time.Millisecond)

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

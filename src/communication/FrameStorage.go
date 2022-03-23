package communication

import (
	"io/ioutil"
	"mjpeg_multiplexer/src/mjpeg"
	"os"
	"sync"
)

type FrameData struct {
	mu    sync.RWMutex
	frame mjpeg.Frame
}

func (frameData *FrameData) Init() {
	//fixme fragile hardcoded path
	var imageLocation = "resources/black.jpg"

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

	frameData.Store(mjpeg.Frame{Body: b})
}

func (frameData *FrameData) Store(data mjpeg.Frame) {
	frameData.mu.Lock()
	defer frameData.mu.Unlock()

	frameData.frame = data
}

func (frameData *FrameData) Get() mjpeg.Frame {
	frameData.mu.RLock()
	defer frameData.mu.RUnlock()

	return frameData.frame
}

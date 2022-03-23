package communication

import (
	"mjpeg_multiplexer/src/mjpeg"
	"sync"
)

type FrameData struct {
	mu    sync.RWMutex
	frame mjpeg.Frame
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

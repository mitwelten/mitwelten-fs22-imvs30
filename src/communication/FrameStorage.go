package communication

import (
	"mjpeg_multiplexer/src/mjpeg"
	"sync"
)

type FrameStorage struct {
	mu    sync.RWMutex
	frame mjpeg.MjpegFrame
}

func (frameData *FrameStorage) Store(data mjpeg.MjpegFrame) {
	frameData.mu.Lock()
	defer frameData.mu.Unlock()

	frameData.frame = data
}

func (frameData *FrameStorage) Get() mjpeg.MjpegFrame {
	frameData.mu.RLock()
	defer frameData.mu.RUnlock()

	return frameData.frame
}

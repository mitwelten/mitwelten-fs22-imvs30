package communication

import (
	"mjpeg_multiplexer/src/mjpeg"
	"sync"
)

type FrameStorage struct {
	mu                  sync.RWMutex
	frame               mjpeg.MjpegFrame
	AggregatorCondition *sync.Cond
}

func NewFrameStorage() *FrameStorage {
	frame := mjpeg.MjpegFrame{}
	frame.Body = mjpeg.Init()

	frameStorage := FrameStorage{}
	frameStorage.Store(frame) // init with a black frame

	return &frameStorage
}

func (frameData *FrameStorage) Store(data mjpeg.MjpegFrame) {
	frameData.mu.Lock()
	defer frameData.mu.Unlock()

	frameData.frame = data

	if frameData.AggregatorCondition != nil {
		frameData.AggregatorCondition.Signal()
	}
}

func (frameData *FrameStorage) Get() mjpeg.MjpegFrame {
	frameData.mu.RLock()
	defer frameData.mu.RUnlock()

	return frameData.frame
}

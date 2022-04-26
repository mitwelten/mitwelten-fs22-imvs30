package communication

import (
	"mjpeg_multiplexer/src/mjpeg"
	"mjpeg_multiplexer/src/utils"
	"sync"
	"time"
)

const nOfStoredFrames = 5

// FrameStorage stores multiple MJPEG frames
type FrameStorage struct {
	mu                  sync.RWMutex
	AggregatorCondition *sync.Cond
	buffer              utils.RingBuffer[mjpeg.MjpegFrame]
	LastUpdated         time.Time
}

// NewFrameStorage FrameStorage ctor
func NewFrameStorage() *FrameStorage {
	frame := mjpeg.MjpegFrame{}
	frame.Body = mjpeg.Init()

	frameStorage := FrameStorage{}
	frameStorage.buffer = utils.NewRingBuffer[mjpeg.MjpegFrame](nOfStoredFrames)
	frameStorage.buffer.Push(frame)
	frameStorage.LastUpdated = time.Now()

	return &frameStorage
}

// Store stores a MjpegFrame into the storage
func (frameStorage *FrameStorage) Store(frame mjpeg.MjpegFrame) {
	frameStorage.mu.RLock()
	defer frameStorage.mu.RUnlock()

	frameStorage.buffer.Push(frame)
	frameStorage.LastUpdated = time.Now()

	if frameStorage.AggregatorCondition != nil {
		frameStorage.AggregatorCondition.Signal()
	}
}

// GetLatest returns the newest frame inserted into the storage
func (frameStorage *FrameStorage) GetLatest() mjpeg.MjpegFrame {
	frameStorage.mu.RLock()
	defer frameStorage.mu.RUnlock()

	return frameStorage.buffer.Peek()
}

// GetAll returns all frames in storage
func (frameStorage *FrameStorage) GetAll() []mjpeg.MjpegFrame {
	frameStorage.mu.RLock()
	defer frameStorage.mu.RUnlock()

	return frameStorage.buffer.GetAll()
}

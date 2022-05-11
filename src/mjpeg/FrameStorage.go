package mjpeg

import (
	"mjpeg_multiplexer/src/utils"
	"sync"
	"time"
)

const nOfStoredFrames = 5

// FrameStorage stores multiple MJPEG frames
type FrameStorage struct {
	mu                  sync.RWMutex
	AggregatorCondition *sync.Cond
	buffer              utils.RingBuffer[MjpegFrame]
	LastUpdated         time.Time
}

// NewFrameStorage FrameStorage ctor
func NewFrameStorage() *FrameStorage {
	frame := MjpegFrame{}
	frame.Body = Init()

	frameStorage := FrameStorage{}
	frameStorage.buffer = utils.NewRingBuffer[MjpegFrame](nOfStoredFrames)
	frameStorage.buffer.Push(frame)
	frameStorage.LastUpdated = time.Now()

	return &frameStorage
}

// Store stores a MjpegFrame into the storage
func (frameStorage *FrameStorage) Store(frame MjpegFrame) {
	frameStorage.mu.RLock()
	defer frameStorage.mu.RUnlock()

	frameStorage.buffer.Push(frame)
	frameStorage.LastUpdated = time.Now()

	if frameStorage.AggregatorCondition != nil {
		frameStorage.AggregatorCondition.Signal()
	}
}

// Store stores a MjpegFrame into the storage
func (frameStorage *FrameStorage) StorePtr(frame *MjpegFrame) {
	frameStorage.mu.RLock()
	defer frameStorage.mu.RUnlock()

	frameStorage.buffer.PushPtr(frame)
	frameStorage.LastUpdated = time.Now()

	if frameStorage.AggregatorCondition != nil {
		frameStorage.AggregatorCondition.Signal()
	}
}

// GetLatest returns the newest frame inserted into the storage
func (frameStorage *FrameStorage) GetLatest() MjpegFrame {
	frameStorage.mu.RLock()
	defer frameStorage.mu.RUnlock()

	return frameStorage.buffer.Peek()
}

// GetLatest returns the newest frame inserted into the storage
func (frameStorage *FrameStorage) GetLatestPtr() *MjpegFrame {
	frameStorage.mu.RLock()
	defer frameStorage.mu.RUnlock()

	return frameStorage.buffer.PeekPtr()
}

// GetAll returns all frames in storage
func (frameStorage *FrameStorage) GetAll() []MjpegFrame {
	frameStorage.mu.RLock()
	defer frameStorage.mu.RUnlock()

	return frameStorage.buffer.GetAll()
}

// GetAll returns all frames in storage
func (frameStorage *FrameStorage) GetAllPtr() []*MjpegFrame {
	frameStorage.mu.RLock()
	defer frameStorage.mu.RUnlock()

	return frameStorage.buffer.GetAllPtr()
}

package mjpeg

import (
	"sync"
	"time"
)

// FrameStorage stores multiple MJPEG frames
type FrameStorage struct {
	mu                     sync.RWMutex
	StorateChangeCondition *sync.Cond
	frame                  MjpegFrame
	//	buffer                 utils.RingBuffer[MjpegFrame]
	LastUpdated time.Time

	imageWidth  int
	imageHeight int
}

func (frameStorage *FrameStorage) GetImageSize() (int, int) {
	frameStorage.mu.RLock()
	defer frameStorage.mu.RUnlock()
	return frameStorage.imageWidth, frameStorage.imageHeight
}

func (frameStorage *FrameStorage) SetImageSize(width, height int) {
	frameStorage.mu.Lock()
	defer frameStorage.mu.Unlock()
	frameStorage.imageWidth = width
	frameStorage.imageHeight = height
}

func CreateUpdateCondition(storages ...*FrameStorage) *sync.Cond {
	lock := sync.Mutex{}
	lock.Lock()
	condition := sync.NewCond(&lock)
	for _, storage := range storages {
		storage.StorateChangeCondition = condition
	}

	return condition
}

// NewFrameStorage FrameStorage ctor
func NewFrameStorage() *FrameStorage {
	frame := NewMJPEGFrame()

	frameStorage := FrameStorage{}
	//frameStorage.buffer = utils.NewRingBuffer[MjpegFrame](nOfStoredFrames)
	//frameStorage.buffer.Push(frame)
	frameStorage.frame = frame
	frameStorage.LastUpdated = time.Now()
	frameStorage.imageWidth = 0
	frameStorage.imageHeight = 0

	return &frameStorage
}

// Store stores a MjpegFrame into the storage
func (frameStorage *FrameStorage) Store(frame *MjpegFrame) {
	frameStorage.mu.RLock()
	defer frameStorage.mu.RUnlock()

	//frameStorage.buffer.Push(frame)
	frameStorage.frame = *frame
	frameStorage.LastUpdated = time.Now()

	if frameStorage.StorateChangeCondition != nil {
		frameStorage.StorateChangeCondition.Signal()
	}
}

// GetFrame returns the newest frame inserted into the storage
func (frameStorage *FrameStorage) GetFrame() *MjpegFrame {
	frameStorage.mu.RLock()
	defer frameStorage.mu.RUnlock()

	return &frameStorage.frame
	//return frameStorage.buffer.PeekPtr()
}

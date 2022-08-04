package mjpeg

import (
	"sync"
	"time"
)

// FrameStorage stores an MJPEG-frame and Signals its condition on update
type FrameStorage struct {
	mu                     sync.RWMutex
	StorageChangeCondition *sync.Cond
	frame                  MjpegFrame
	LastUpdated            time.Time

	imageWidth  int
	imageHeight int
}

//GetImageSize gets the image size of the frames stored in this storage
func (frameStorage *FrameStorage) GetImageSize() (int, int) {
	frameStorage.mu.RLock()
	defer frameStorage.mu.RUnlock()
	return frameStorage.imageWidth, frameStorage.imageHeight
}

//SetImageSize updates the image size of the frames stored in this storage
func (frameStorage *FrameStorage) SetImageSize(width, height int) {
	frameStorage.mu.Lock()
	defer frameStorage.mu.Unlock()
	frameStorage.imageWidth = width
	frameStorage.imageHeight = height
}

//CreateUpdateCondition creates one condition for all the passed storages
func CreateUpdateCondition(storages ...*FrameStorage) *sync.Cond {
	lock := sync.Mutex{}
	lock.Lock()
	condition := sync.NewCond(&lock)
	for _, storage := range storages {
		storage.StorageChangeCondition = condition
	}

	return condition
}

// NewFrameStorage FrameStorage ctor
func NewFrameStorage() *FrameStorage {
	frame := NewMJPEGFrame()

	frameStorage := FrameStorage{}
	frameStorage.frame = frame
	frameStorage.LastUpdated = time.Now()
	frameStorage.imageWidth = 0
	frameStorage.imageHeight = 0

	return &frameStorage
}

// Store stores a MjpegFrame into the storage
func (frameStorage *FrameStorage) Store(frame *MjpegFrame) {
	frameStorage.mu.Lock()
	defer frameStorage.mu.Unlock()

	frameStorage.frame = *frame
	frameStorage.LastUpdated = time.Now()

	if frameStorage.StorageChangeCondition != nil {
		frameStorage.StorageChangeCondition.Signal()
	}
}

// GetFrame returns the newest frame inserted into the storage
func (frameStorage *FrameStorage) GetFrame() *MjpegFrame {
	frameStorage.mu.RLock()
	defer frameStorage.mu.RUnlock()

	return &frameStorage.frame
}

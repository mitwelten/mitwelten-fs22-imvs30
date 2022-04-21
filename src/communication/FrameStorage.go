package communication

import (
	"mjpeg_multiplexer/src/mjpeg"
	"sync"
)

const nOfStoredFrames = 5

type FrameStorage struct {
	mu                  sync.RWMutex
	frames              [nOfStoredFrames]mjpeg.MjpegFrame
	AggregatorCondition *sync.Cond
	currentFramePointer int
}

func NewFrameStorage() *FrameStorage {
	frame := mjpeg.MjpegFrame{}
	frame.Body = mjpeg.Init()

	frameStorage := FrameStorage{}
	frameStorage.currentFramePointer = 0 // init pointer at position 0
	frameStorage.Store(frame)            // init with a black frame

	return &frameStorage
}

func (frameStorage *FrameStorage) Store(frame mjpeg.MjpegFrame) {
	frameStorage.set(frame)

	if frameStorage.AggregatorCondition != nil {
		frameStorage.AggregatorCondition.Signal()
	}
}

func (frameStorage *FrameStorage) GetLatest() mjpeg.MjpegFrame {
	frameStorage.mu.RLock()
	defer frameStorage.mu.RUnlock()

	index := (frameStorage.currentFramePointer + nOfStoredFrames - 1) % nOfStoredFrames
	return frameStorage.frames[index]
}

func (frameStorage *FrameStorage) GetAll() []mjpeg.MjpegFrame {
	frameStorage.mu.RLock()
	defer frameStorage.mu.RUnlock()

	currentFrame := frameStorage.frames[frameStorage.currentFramePointer]

	var nFrames int
	if len(currentFrame.Body) > 0 {
		nFrames = nOfStoredFrames
	} else {
		nFrames = frameStorage.currentFramePointer
	}

	outputFrames := make([]mjpeg.MjpegFrame, nFrames)

	index := (frameStorage.currentFramePointer + nOfStoredFrames - 1) % nOfStoredFrames
	for i := 0; i < nFrames; i++ {
		outputFrames[i] = frameStorage.frames[index]
		index = (index + nOfStoredFrames - 1) % nOfStoredFrames
	}

	return outputFrames
}

func (frameStorage *FrameStorage) set(frame mjpeg.MjpegFrame) {
	frameStorage.mu.RLock()
	defer frameStorage.mu.RUnlock()

	frameStorage.frames[frameStorage.currentFramePointer] = frame
	frameStorage.currentFramePointer = (frameStorage.currentFramePointer + 1) % nOfStoredFrames
}

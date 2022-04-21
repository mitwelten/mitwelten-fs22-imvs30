package communication

import (
	"mjpeg_multiplexer/src/mjpeg"
	"mjpeg_multiplexer/src/tests"
	"testing"
)

// proxies requests to the golang.org playground service.
func getFrame(content int) mjpeg.MjpegFrame {
	return mjpeg.MjpegFrame{Body: []byte{byte(content)}}
}

// TestFrameStorage_GetLast_Overflow creates an overflow in data structure and tests correct behavior
func TestFrameStorage_GetLatest(t *testing.T) {
	frameStorage := NewFrameStorage()

	frameStorage.Store(getFrame(0))

	for i := 0; i < nOfStoredFrames; i++ {
		frameStorage.Store(getFrame(i))
		tests.Assert(t, getFrame(i), frameStorage.GetLatest(), true)
	}

}

//TestFrameStorage_GetAll tests if the correct amount of frames are returned in the correct order
func TestFrameStorage_GetAll(t *testing.T) {
	// when
	frameStorage := NewFrameStorage()
	// then expect default frame
	tests.Assert(t, 1, len(frameStorage.GetAll()), true)

	// when
	for i := 1; i < nOfStoredFrames; i++ {
		frameStorage.Store(getFrame(i))
		// then expect frames stored in framestorage
		tests.Assert(t, i+1, len(frameStorage.GetAll()), true)
	}

	// expect GetAll() to return the newest frames first ('newest' = added last)
	frames := make([]mjpeg.MjpegFrame, nOfStoredFrames)
	for i := 0; i < nOfStoredFrames; i++ {
		frameStorage.Store(getFrame(i))
		frames[nOfStoredFrames-i-1] = getFrame(i)
	}

	tests.Assert(t, frames, frameStorage.GetAll(), true)
}

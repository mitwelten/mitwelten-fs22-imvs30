package mjpeg

import (
	"mjpeg_multiplexer/src/utils"
	"testing"
)

// proxies requests to the golang.org playground service.
func getFrame(content int) *MjpegFrame {
	return &MjpegFrame{Body: []byte{byte(content)}}
}

// TestFrameStorage_GetLast_Overflow creates an overflow in data structure and util correct behavior
func TestFrameStorage_GetLatest(t *testing.T) {
	frameStorage := NewFrameStorage()

	for i := 0; i < 10; i++ {
		frameStorage.Store(getFrame(i))
		frameStorage.Store(getFrame(i))
		utils.Assert(t, getFrame(i), frameStorage.GetFrame(), true)
	}

}

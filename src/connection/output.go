package connection

import (
	"mjpeg_multiplexer/src/communication"
	"mjpeg_multiplexer/src/mjpeg"
)

type Output interface {
	SendFrame(frame mjpeg.MjpegFrame) error
	Run(storage *communication.FrameStorage)
}

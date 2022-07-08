package connection

import (
	"mjpeg_multiplexer/src/mjpeg"
)

type Output interface {
	SendFrame(frame *mjpeg.MjpegFrame)
	StartOutput()
}

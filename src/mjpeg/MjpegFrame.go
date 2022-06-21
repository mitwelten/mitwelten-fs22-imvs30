package mjpeg

import (
	_ "embed"
	"image"
)

var JPEG_PREFIX = []byte("\xff\xd8")
var FRAME_DELIM = []byte("--boundarydonotcross\r\n")

type MjpegFrame struct {
	Body  []byte
	Empty bool

	CachedImage image.Image
}

//go:embed black.jpg
var blackJPG []byte

func NewMJPEGFrame() MjpegFrame {
	frame := MjpegFrame{}
	frame.Body = blackJPG
	frame.Empty = true
	return frame
}

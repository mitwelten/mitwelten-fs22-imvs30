package mjpeg

import (
	_ "embed"
	"image"
)

// MjpegFrame is a MJPEG frame, containing the jpeg bytes and a possible cached decoded image
type MjpegFrame struct {
	Body []byte
	//is the image in the Body the black placeholder image?
	Empty       bool
	CachedImage image.Image
}

//go:embed black.jpg
var blackJPG []byte

// NewMJPEGFrame creates a new mjpeg with a placeholder frame
func NewMJPEGFrame() MjpegFrame {
	frame := MjpegFrame{}
	frame.Body = blackJPG
	frame.Empty = true
	return frame
}

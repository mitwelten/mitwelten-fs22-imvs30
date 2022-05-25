package mjpeg

import (
	"bytes"
	_ "embed"
	"image"
)

var JPEG_PREFIX = []byte("\xff\xd8")
var FRAME_DELIM = []byte("--boundarydonotcross\r\n")

type MjpegFrame struct {
	Body           []byte
	CachedImage    image.Image
	Resized        bool
	OriginalHeight int
	OriginalWidth  int
}

//go:embed black.jpg
var blackJPG []byte

func Init() []byte {
	return blackJPG
}

func parse_frame(data []byte) (frame MjpegFrame) {
	for i := 0; i < len(data); i++ {
		if bytes.Compare(data[i:i+2], JPEG_PREFIX) == 0 {
			return MjpegFrame{data[i:], nil, false, -1, -1}
		}
	}
	panic("Can't parse frame")
}

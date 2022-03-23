package mjpeg

import (
	"bytes"
	_ "embed"
)

var JPEG_PREFIX = []byte("\xff\xd8")
var FRAME_DELIM = []byte("--boundarydonotcross\r\n")

type Frame struct {
	Body []byte
}

//go:embed black.jpg
var blackJPG []byte

func Init() []byte {
	return blackJPG
}

func parse_frame(data []byte) (frame Frame) {
	for i := 0; i < len(data); i++ {
		if bytes.Compare(data[i:i+2], JPEG_PREFIX) == 0 {
			return Frame{data[i:]}
		}
	}
	panic("Can't parse frame")
}

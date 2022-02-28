package mjpeg

import (
	"bytes"
)

var JPEG_PREFIX = []byte("\xff\xd8")
var FRAME_DELIM = []byte("--boundarydonotcross\r\n")

type Frame struct {
	Header []byte
	Body   []byte
}

func parse_frame(data []byte) (frame Frame) {
	for i := 0; i < len(data); i++ {
		if bytes.Compare(data[i:i+2], JPEG_PREFIX) == 0 {
			return Frame{data[:i], data[i:]}
		}
	}
	panic("Can't parse frame")
}

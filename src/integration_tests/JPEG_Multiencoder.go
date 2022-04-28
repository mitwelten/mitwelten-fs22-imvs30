package integration_tests

import (
	"bytes"
	_ "embed"
	"image"
	"image/draw"
	"image/jpeg"
	"os"
	"testing"
)

//go:embed file.jpg
var file []byte

func EncodeAndDecode(in []byte) []byte {
	img, _, _ := image.Decode(bytes.NewReader(in))
	rect := img.Bounds()
	imageOut := image.NewRGBA(rect)
	draw.Draw(imageOut, rect, img, image.Point{}, draw.Src)

	buff := bytes.NewBuffer([]byte{})
	options := jpeg.Options{Quality: 100}
	err := jpeg.Encode(buff, imageOut, &options)
	if err != nil {
		panic("can't encode jpg")
	}

	return buff.Bytes()
}

func TestJPEG(t *testing.T) {
	data := file
	for i := 0; i < 1_000; i++ {
		data = EncodeAndDecode(data)
	}

	fh, err := os.OpenFile("out.jpg", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		panic("can't encode jpg")

	}
	_, err = fh.Write(data)
	if err != nil {
		panic("can't encode jpg")
	}
	fh.Close()
}

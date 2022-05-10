package image

// Reference: https://go.dev/blog/image-draw

import (
	"bytes"
	"github.com/pixiv/go-libjpeg/jpeg"
	"image"
	"image/draw"
	"mjpeg_multiplexer/src/mjpeg"
)

var DecodeOptions = jpeg.DecoderOptions{ScaleTarget: image.Rectangle{}, DCTMethod: jpeg.DCTIFast, DisableFancyUpsampling: false, DisableBlockSmoothing: false}
var EncodingOptions = jpeg.EncoderOptions{Quality: 100, OptimizeCoding: false, ProgressiveMode: false, DCTMethod: jpeg.DCTISlow}

func Decode(frames ...mjpeg.MjpegFrame) []image.Image {
	var images []image.Image
	for _, frame := range frames {
		var img, _ = jpeg.Decode(bytes.NewReader(frame.Body), &DecodeOptions)
		images = append(images, img)
	}
	return images
}

func Encode(images ...image.Image) []mjpeg.MjpegFrame {
	var frames []mjpeg.MjpegFrame
	for _, img := range images {
		buff := bytes.NewBuffer([]byte{})
		err := jpeg.Encode(buff, img, &EncodingOptions)

		if err != nil {
			panic("can't encode jpg")
		}

		frames = append(frames, mjpeg.MjpegFrame{Body: buff.Bytes()})
	}

	return frames

}

func Grid(row int, col int, frames ...mjpeg.MjpegFrame) mjpeg.MjpegFrame {
	var nCells = row * col
	var nFrames = len(frames)

	if nFrames > nCells {
		panic("Too many frames")
	}

	if nFrames == 0 {
		panic("At least one frame needed")
	}
	var images = Decode(frames...)

	// rectangle
	var i0 = images[0]
	var pointMax = image.Point{X: i0.Bounds().Dx() * col, Y: i0.Bounds().Dy() * row}
	var rectangle = image.Rectangle{Min: image.Point{}, Max: pointMax}

	// image
	var imageOut = image.NewRGBA(rectangle)

	for i := 0; i < nCells; i++ {
		var row_ = i / col
		var col_ = i % col

		if i >= nFrames {
			break
		}

		var sp = image.Point{X: i0.Bounds().Dx() * col_, Y: i0.Bounds().Dy() * row_}
		var r = image.Rectangle{Min: sp, Max: sp.Add(images[i].Bounds().Size())}
		draw.Draw(imageOut, r, images[i], image.Point{}, draw.Src)
	}

	return Encode(imageOut)[0]
}

package imageUtils

// Reference: https://go.dev/blog/image-draw

import (
	"bytes"
	"github.com/pixiv/go-libjpeg/jpeg"
	"golang.org/x/image/draw"
	"image"
	"log"
	"mjpeg_multiplexer/src/global"
	"mjpeg_multiplexer/src/mjpeg"
	"mjpeg_multiplexer/src/utils"
)

var DecodeOptions = jpeg.DecoderOptions{ScaleTarget: image.Rectangle{}, DCTMethod: jpeg.DCTIFast, DisableFancyUpsampling: false, DisableBlockSmoothing: false}
var EncodingOptions = jpeg.EncoderOptions{Quality: 100, OptimizeCoding: false, ProgressiveMode: false, DCTMethod: jpeg.DCTISlow}

func DecodeAll(frames ...*mjpeg.MjpegFrame) []image.Image {
	var images []image.Image
	for _, frame := range frames {
		var img = Decode(frame)
		images = append(images, img)
	}
	return images
}

func Decode(frame *mjpeg.MjpegFrame) image.Image {
	if frame.Image == nil {
		img, err := jpeg.Decode(bytes.NewReader(frame.Body), &DecodeOptions)
		//todo: this assigns to a copy
		frame.Image = img

		if err != nil {
			panic("can't decode jpg")
		}

		return img
	} else {
		return frame.Image
	}
}

func EncodeAll(images ...image.Image) []*mjpeg.MjpegFrame {
	var frames []*mjpeg.MjpegFrame
	for _, img := range images {
		imageOut := Encode(img)
		frames = append(frames, imageOut)
	}

	return frames
}

func Encode(image image.Image) *mjpeg.MjpegFrame {
	buff := bytes.NewBuffer([]byte{})
	err := jpeg.Encode(buff, image, &EncodingOptions)

	if err != nil {
		panic("can't encode jpg")
	}

	return &mjpeg.MjpegFrame{Body: buff.Bytes()}
}

func Grid(row int, col int, frames ...*mjpeg.MjpegFrame) *mjpeg.MjpegFrame {
	var nCells = row * col
	var nFrames = len(frames)

	if nFrames > nCells {
		log.Fatalf("Too many frames for this grid configuartion: row %v col %v, but %v frames to compute\n", row, col, len(frames))
	}

	if nFrames == 0 {
		log.Fatalf("At least one frame needed\n")
	}
	var images = DecodeAll(frames...)

	//calculate the final grid size
	//todo invalid in certain cases
	height := 0
	width := 0
	for i := 0; i < nCells; i++ {
		var row_ = i / col
		var col_ = i % col
		if row_ == 0 {
			if frames[i].Resized {
				width += frames[i].OriginalWidth
			} else {
				width += images[i].Bounds().Dx()
			}
		}
		if col_ == 0 {
			if frames[i].Resized {
				height += frames[i].OriginalHeight
			} else {
				height += images[i].Bounds().Dy()
			}
		}
	}

	// check the max width and height
	if (global.Config.MaxHeight != -1 && global.Config.MaxHeight < height) || (global.Config.MaxWidth != -1 && global.Config.MaxWidth < width) {
		deltaW := 1
		deltaH := 1

		if global.Config.MaxWidth != -1 {
			deltaW = width / global.Config.MaxWidth
		}

		if global.Config.MaxHeight != -1 {
			deltaH = height / global.Config.MaxHeight
		}

		targetWidth := (width / deltaW) / col
		targetHeight := (height / deltaH) / row

		// resize all images if needed
		for i, _ := range images {
			if images[i].Bounds().Dx() != targetWidth || images[i].Bounds().Dy() != targetHeight {
				frames[i].OriginalWidth = images[i].Bounds().Dx()
				frames[i].OriginalHeight = images[i].Bounds().Dy()
				frames[i].Resized = true

				images[i] = Resize(images[i], targetWidth, targetHeight)
				frames[i].Image = images[i]
			}
		}
	}

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

	return Encode(imageOut)
}

func Resize(img image.Image, width int, height int) image.Image {
	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.NearestNeighbor.Scale(dst, dst.Rect, img, img.Bounds(), draw.Over, nil)
	return dst
}

func GetFrameStorageSize(input *mjpeg.FrameStorage) (int, int, image.Image) {
	width, height := input.GetImageSize()

	// check if height or width has been set in the storage
	var img image.Image = nil
	if width == -1 || height == -1 {
		img = Decode(input.GetLatestPtr())
		input.SetImageSize(img.Bounds().Dx(), img.Bounds().Dy())
		width = img.Bounds().Dx()
		height = img.Bounds().Dy()
	}
	return width, height, img
}

func Transform(input *mjpeg.FrameStorage) *mjpeg.MjpegFrame {
	width, height, img := GetFrameStorageSize(input)

	if global.Config.MaxWidth == -1 && global.Config.MaxHeight == -1 {
		// no global maxHeight or maxWidth set
		return input.GetLatestPtr()
	}

	// check if resizing is needed
	if (global.Config.MaxWidth < width) || (global.Config.MaxHeight < height) {
		if img == nil {
			img = Decode(input.GetLatestPtr())
		}

		resized := Resize(img, utils.Min(width, global.Config.MaxWidth), utils.Min(height, global.Config.MaxHeight))
		return Encode(resized)
	}

	// else just return the original frame
	return input.GetLatestPtr()
}

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

type PanelLayout struct {
	FirstWidth  float64
	FirstHeight float64

	ChildrenWidth  float64
	ChildrenHeight float64

	ChildrenPositions []utils.FloatPoint
}

var Slots3 = PanelLayout{
	FirstWidth:     float64(2) / 3.0,
	FirstHeight:    1,
	ChildrenWidth:  float64(1) / 3,
	ChildrenHeight: float64(1) / 2,
	ChildrenPositions: []utils.FloatPoint{
		{X: float64(2) / 3, Y: 0},
		{X: float64(2) / 3, Y: float64(1) / 2},
	},
}

var Slots4 = PanelLayout{
	FirstWidth:     float64(3) / 4.0,
	FirstHeight:    1,
	ChildrenWidth:  float64(1) / 4,
	ChildrenHeight: float64(1) / 3,
	ChildrenPositions: []utils.FloatPoint{
		{X: float64(3) / 4, Y: 0},
		{X: float64(3) / 4, Y: float64(1) / 3},
		{X: float64(3) / 4, Y: float64(2) / 3},
	},
}

var Slots6 = PanelLayout{
	FirstWidth:     float64(2) / 3.0,
	FirstHeight:    float64(2) / 3.0,
	ChildrenWidth:  float64(1) / 3,
	ChildrenHeight: float64(1) / 3,
	ChildrenPositions: []utils.FloatPoint{
		{X: 0, Y: float64(2) / 3},
		{X: float64(1) / 3, Y: float64(2) / 3},
		{X: float64(2) / 3, Y: float64(2) / 3},
		{X: float64(2) / 3, Y: float64(1) / 3},
		{X: float64(2) / 3, Y: 0},
	},
}

var Slots8 = PanelLayout{
	FirstWidth:     float64(3) / 4.0,
	FirstHeight:    float64(3) / 4.0,
	ChildrenWidth:  float64(1) / 4,
	ChildrenHeight: float64(1) / 4,
	ChildrenPositions: []utils.FloatPoint{
		{X: 0, Y: float64(3) / 4},
		{X: float64(1) / 4, Y: float64(3) / 4},
		{X: float64(2) / 4, Y: float64(3) / 4},
		{X: float64(3) / 4, Y: float64(3) / 4},
		{X: float64(3) / 4, Y: float64(2) / 4},
		{X: float64(3) / 4, Y: float64(1) / 4},
		{X: float64(3) / 4, Y: 0},
	},
}

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
	config := EncodingOptions
	config.Quality = global.Config.EncodeQuality
	err := jpeg.Encode(buff, image, &config)

	if err != nil {
		panic("can't encode jpg")
	}

	return &mjpeg.MjpegFrame{Body: buff.Bytes()}
}

func Panel(layout PanelLayout, startIndex int, storages ...*mjpeg.FrameStorage) *mjpeg.MjpegFrame {
	//todo optimize by caching the result and painting over it?
	var frames []*mjpeg.MjpegFrame
	for i := 0; i < len(storages); i++ {
		frames = append(frames, storages[i].GetLatestPtr())
	}
	var images = DecodeAll(frames...)

	firstWidthInitial, firstHeightInitial, _ := GetFrameStorageSize(storages[0])

	totalWidth := int(float64(firstWidthInitial) / layout.FirstWidth)
	totalHeight := int(float64(firstHeightInitial) / layout.FirstHeight)

	if global.Config.MaxWidth != -1 {
		totalWidth = utils.Min(totalWidth, global.Config.MaxWidth)
	}

	if global.Config.MaxHeight != -1 {
		totalHeight = utils.Min(totalHeight, global.Config.MaxHeight)
	}

	//output img
	rectangle := image.Rectangle{Max: image.Point{X: totalWidth, Y: totalHeight}}
	imageOut := image.NewRGBA(rectangle)

	//main character image
	sp := image.Point{}
	r := image.Rectangle{Min: sp, Max: image.Point{X: int(float64(totalWidth) * layout.FirstWidth), Y: int(float64(totalHeight) * layout.FirstHeight)}}
	draw.NearestNeighbor.Scale(imageOut, r, images[startIndex], images[startIndex].Bounds(), draw.Over, nil)

	for i, child := range layout.ChildrenPositions {
		if i+1 >= len(storages) {
			break
		}
		index := (startIndex + i + 1) % len(storages)
		sp = image.Point{X: int(float64(totalWidth) * child.X), Y: int(float64(totalHeight) * child.Y)}
		delta := image.Point{X: int(float64(totalWidth) * layout.ChildrenWidth), Y: int(float64(totalHeight) * layout.ChildrenHeight)}
		r = image.Rectangle{Min: sp, Max: sp.Add(delta)}
		draw.NearestNeighbor.Scale(imageOut, r, images[index], images[index].Bounds(), draw.Over, nil)
	}

	return Encode(imageOut)
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

	// with default settings (no resize, no quality change) the image doesn't need to be decoded and encoded at all
	if !global.DecodingNecessary() {
		return input.GetLatestPtr()
	}

	// check if resizing is needed
	if (global.Config.MaxWidth < width) || (global.Config.MaxHeight < height) {
		if img == nil {
			img = Decode(input.GetLatestPtr())
		}

		outputWidth := width
		if global.Config.MaxWidth != -1 && global.Config.MaxWidth < width {
			outputWidth = global.Config.MaxWidth
		}

		outputHeight := height
		if global.Config.MaxHeight != -1 && global.Config.MaxHeight < height {
			outputHeight = global.Config.MaxHeight
		}

		resized := Resize(img, outputWidth, outputHeight)
		return Encode(resized)
	}

	// Else just decode and encode to adjust the quality
	return Encode(Decode(input.GetLatestPtr()))

}

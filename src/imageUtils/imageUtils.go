package imageUtils

// Reference: https://go.dev/blog/image-draw

import (
	"bytes"
	_ "embed"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/pixiv/go-libjpeg/jpeg"
	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"log"
	"mjpeg_multiplexer/src/global"
	"mjpeg_multiplexer/src/mjpeg"
	"mjpeg_multiplexer/src/utils"
)

//go:embed arial.ttf
var arial []byte

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

const borderFactor = 0.0025

func DecodeAll(frames ...*mjpeg.MjpegFrame) []image.Image {
	var images []image.Image
	for _, frame := range frames {
		var img = Decode(frame)
		images = append(images, img)
	}
	return images
}

func Decode(frame *mjpeg.MjpegFrame) image.Image {
	if frame.CachedImage == nil {
		img, err := jpeg.Decode(bytes.NewReader(frame.Body), &DecodeOptions)
		frame.CachedImage = img

		if err != nil {
			panic("can't decode jpg")
		}

		return img
	} else {
		return frame.CachedImage
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

	totalWidth, totalHeight = GetFinalImageSize(totalWidth, totalHeight)

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

func Grid(nRows int, nCols int, storages ...*mjpeg.FrameStorage) *mjpeg.MjpegFrame {
	var nCells = nRows * nCols
	var nFrames = len(storages)

	if nFrames > nCells {
		log.Fatalf("Too many frames for this grid configuartion: nRows %v nCols %v, but %v frames to compute\n", nRows, nCols, len(storages))
	}

	if nFrames == 0 {
		log.Fatalf("At least one frame needed\n")
	}

	var frames []*mjpeg.MjpegFrame
	for i := 0; i < len(storages); i++ {
		frames = append(frames, storages[i].GetLatestPtr())
	}

	var images = DecodeAll(frames...)
	firstWidthInitial, firstHeightInitial, _ := GetFrameStorageSize(storages[0])
	totalWidth := firstWidthInitial * nCols
	totalHeight := firstHeightInitial * nRows

	var cellWidth int
	var cellHeight int
	// check the max totalWidth and totalHeight

	totalWidth, totalHeight = GetFinalImageSize(totalWidth, totalHeight)

	cellWidth = totalWidth / nCols
	cellHeight = totalHeight / nRows

	// resize all images if needed
	for i, _ := range images {
		if images[i].Bounds().Dx() != cellWidth || images[i].Bounds().Dy() != cellHeight {
			images[i] = Resize(images[i], cellWidth, cellHeight)
			//frames[i].CachedImage = images[i]
		}
	}

	// rectangle
	var pointMax = image.Point{X: totalWidth, Y: totalHeight}
	var rectangle = image.Rectangle{Min: image.Point{}, Max: pointMax}

	// image
	var imageOut = image.NewRGBA(rectangle)

	var marginStartPoints []image.Point

	for i := 0; i < nCells; i++ {
		var row_ = i / nCols
		var col_ = i % nCols

		if i >= nFrames {
			break
		}

		var sp = image.Point{X: cellWidth * col_, Y: cellHeight * row_}
		//grid lines:

		if row_ == 0 || col_ == 0 {
			marginStartPoints = append(marginStartPoints, sp)
		}

		var r = image.Rectangle{Min: sp, Max: sp.Add(images[i].Bounds().Size())}
		draw.Draw(imageOut, r, images[i], image.Point{}, draw.Src)
	}

	//place labels
	if global.Config.ShowInputLabel {

		for i := 0; i < nCells; i++ {
			var row_ = i / nCols
			var col_ = i % nCols

			if i >= nFrames {
				break
			}

			var sp = image.Point{X: cellWidth * col_, Y: cellHeight * row_}
			// adjust the start point relative to the border size
			x := sp.X
			y := sp.Y
			if col_ != 0 {
				x += int(float64(totalWidth)*borderFactor) / 2
			}
			if row_ != 0 {
				y += int(float64(totalWidth)*borderFactor) / 2
			}

			addLabel(imageOut, x, y, global.Config.InputConfigs[i].Url)

		}
	}

	// draw border
	if global.Config.Border {
		border := int(float64(totalWidth) * borderFactor)
		borderVertical := image.Rectangle{Min: image.Point{X: -border / 2}, Max: image.Point{X: border / 2, Y: imageOut.Bounds().Dy()}}
		borderHorizontal := image.Rectangle{Min: image.Point{Y: -border / 2}, Max: image.Point{X: imageOut.Bounds().Dx(), Y: border / 2}}
		for i, p := range marginStartPoints {
			//ignore first point to avoid border lines
			if i == 0 {
				continue
			}

			if p.Y == 0 {
				//vertical
				draw.Draw(imageOut, borderVertical.Add(p), image.Black, image.Point{}, draw.Src)
			} else {
				//horizontal
				draw.Draw(imageOut, borderHorizontal.Add(p), image.Black, image.Point{}, draw.Src)
			}
		}
	}

	return Encode(imageOut)
}

//ResizeStorage resizes and image and saved the resized image to the storage
func ResizeStorage(frame *mjpeg.MjpegFrame, img image.Image, width int, height int) image.Image {
	frame.OriginalWidth = img.Bounds().Dx()
	frame.OriginalHeight = img.Bounds().Dy()
	frame.Resized = true

	return Resize(img, width, height)
}

//Resize resizes an image
func Resize(img image.Image, width int, height int) image.Image {
	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.NearestNeighbor.Scale(dst, dst.Rect, img, img.Bounds(), draw.Over, nil)
	return dst
}

//GetFinalImageSize returns the size of the image in regard to the globally set MaxWidth and MinWidth
func GetFinalImageSize(currentWidth int, currentHeight int) (int, int) {
	if global.Config.MaxWidth != -1 && global.Config.MaxHeight != -1 {
		// both dimensions may be resized to the desired size
		currentWidth = utils.Min(currentWidth, global.Config.MaxWidth)
		currentHeight = utils.Min(currentHeight, global.Config.MaxHeight)
	} else if global.Config.MaxWidth != -1 && global.Config.MaxHeight == -1 {
		// reduce width and scale height
		if currentWidth > global.Config.MaxWidth {
			currentHeight = int(float64(currentHeight) * (float64(global.Config.MaxWidth) / float64(currentWidth)))
			currentWidth = global.Config.MaxWidth
		}
	} else if global.Config.MaxWidth == -1 && global.Config.MaxHeight != -1 {
		// reduce height and scale width
		if currentHeight > global.Config.MaxHeight {
			currentWidth = int(float64(currentWidth) * (float64(global.Config.MaxHeight) / float64(currentHeight)))
			currentHeight = global.Config.MaxHeight
		}
	}

	return currentWidth, currentHeight
}

func GetFrameStorageSize(input *mjpeg.FrameStorage) (int, int, image.Image) {
	width, height := input.GetImageSize()

	// check if height or width has been set in the storage
	var img image.Image = nil
	if width == -1 || height == -1 {
		img = Decode(input.GetLatestPtr())
		width = img.Bounds().Dx()
		height = img.Bounds().Dy()

		// store the result, but only if it's not the empty image
		if !input.GetLatestPtr().Empty {
			input.SetImageSize(img.Bounds().Dx(), img.Bounds().Dy())
		}
	}
	return width, height, img
}

func Transform(input *mjpeg.FrameStorage) *mjpeg.MjpegFrame {
	// todo more caching here?

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

		outputWidth, outputHeight := GetFinalImageSize(width, height)

		resized := ResizeStorage(input.GetLatestPtr(), img, outputWidth, outputHeight)
		return Encode(resized)
	}

	// Else just decode and encode to adjust the quality
	return Encode(Decode(input.GetLatestPtr()))

}

const labelSize = 36
const padding = 8

var labelSrc = image.NewUniform(color.RGBA{R: 255, G: 255, B: 255, A: 255})
var f, _ = freetype.ParseFont(arial)
var face = truetype.NewFace(f, &truetype.Options{
	Size: labelSize,
})

func addLabel(img *image.RGBA, x, y int, label string) {
	point := fixed.Point26_6{X: fixed.I(x + padding), Y: fixed.I(y + labelSize)}

	d := &font.Drawer{
		Dst:  img,
		Src:  labelSrc,
		Face: face,
		Dot:  point,
	}

	textWidth := d.MeasureString(label).Round() + 2*padding
	textHeight := labelSize + padding

	// draw the black border...
	var r = image.Rectangle{Min: image.Point{X: x, Y: y}, Max: image.Point{X: x + textWidth, Y: y + textHeight}}
	draw.Draw(img, r, image.Black, image.Point{}, draw.Src)

	// and the text
	d.DrawString(label)
}

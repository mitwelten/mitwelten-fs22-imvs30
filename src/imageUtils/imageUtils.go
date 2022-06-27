package imageUtils

// Reference: https://go.dev/blog/image-draw

import (
	"bytes"
	_ "embed"
	"image"
	"image/color"
	"log"
	"math"
	"mjpeg_multiplexer/src/global"
	"mjpeg_multiplexer/src/mjpeg"
	"mjpeg_multiplexer/src/utils"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/pixiv/go-libjpeg/jpeg"
	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

//go:embed arial.ttf
var arial []byte

var DecodeOptions = jpeg.DecoderOptions{ScaleTarget: image.Rectangle{}, DCTMethod: jpeg.DCTIFast, DisableFancyUpsampling: false, DisableBlockSmoothing: false}
var EncodingOptions = jpeg.EncoderOptions{Quality: 100, OptimizeCoding: false, ProgressiveMode: false, DCTMethod: jpeg.DCTIFast}

type PanelLayout struct {
	FirstWidth  float64
	FirstHeight float64

	ChildrenWidth  float64
	ChildrenHeight float64

	ChildrenPositions []utils.FloatPoint

	VerticalBorderPoints   []utils.Tuple[utils.FloatPoint]
	HorizontalBorderPoints []utils.Tuple[utils.FloatPoint]
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

	VerticalBorderPoints: []utils.Tuple[utils.FloatPoint]{
		{T1: utils.FloatPoint{X: float64(2) / 3, Y: 0}, T2: utils.FloatPoint{X: float64(2) / 3, Y: 1}},
	},
	HorizontalBorderPoints: []utils.Tuple[utils.FloatPoint]{
		{T1: utils.FloatPoint{X: float64(2) / 3, Y: 0.5}, T2: utils.FloatPoint{X: 1, Y: 0.5}},
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
	VerticalBorderPoints: []utils.Tuple[utils.FloatPoint]{
		{T1: utils.FloatPoint{X: float64(3) / 4, Y: 0}, T2: utils.FloatPoint{X: float64(3) / 4, Y: 1}},
	},
	HorizontalBorderPoints: []utils.Tuple[utils.FloatPoint]{
		{T1: utils.FloatPoint{X: float64(3) / 4, Y: float64(1) / 3}, T2: utils.FloatPoint{X: 1, Y: float64(1) / 3}},
		{T1: utils.FloatPoint{X: float64(3) / 4, Y: float64(2) / 3}, T2: utils.FloatPoint{X: 1, Y: float64(2) / 3}},
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
	VerticalBorderPoints: []utils.Tuple[utils.FloatPoint]{
		{T1: utils.FloatPoint{X: float64(1) / 3, Y: float64(2) / 3}, T2: utils.FloatPoint{X: float64(1) / 3, Y: 1}},
		{T1: utils.FloatPoint{X: float64(2) / 3, Y: 0}, T2: utils.FloatPoint{X: float64(2) / 3, Y: 1}},
	},
	HorizontalBorderPoints: []utils.Tuple[utils.FloatPoint]{
		{T1: utils.FloatPoint{X: float64(2) / 3, Y: float64(1) / 3}, T2: utils.FloatPoint{X: 1, Y: float64(1) / 3}},
		{T1: utils.FloatPoint{X: 0, Y: float64(2) / 3}, T2: utils.FloatPoint{X: 1, Y: float64(2) / 3}},
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
	VerticalBorderPoints: []utils.Tuple[utils.FloatPoint]{
		{T1: utils.FloatPoint{X: float64(1) / 4, Y: float64(3) / 4}, T2: utils.FloatPoint{X: float64(1) / 4, Y: 1}},
		{T1: utils.FloatPoint{X: float64(2) / 4, Y: float64(3) / 4}, T2: utils.FloatPoint{X: float64(2) / 4, Y: 1}},
		{T1: utils.FloatPoint{X: float64(3) / 4, Y: 0}, T2: utils.FloatPoint{X: float64(3) / 4, Y: 1}},
	},
	HorizontalBorderPoints: []utils.Tuple[utils.FloatPoint]{
		{T1: utils.FloatPoint{X: float64(3) / 4, Y: float64(1) / 4}, T2: utils.FloatPoint{X: 1, Y: float64(1) / 4}},
		{T1: utils.FloatPoint{X: float64(3) / 4, Y: float64(2) / 4}, T2: utils.FloatPoint{X: 1, Y: float64(2) / 4}},
		{T1: utils.FloatPoint{X: 0, Y: float64(3) / 4}, T2: utils.FloatPoint{X: 1, Y: float64(3) / 4}},
	},
}

func DecodeAll(storages ...*mjpeg.FrameStorage) []image.Image {
	var images []image.Image
	for _, storage := range storages {
		var img = Decode(storage)
		images = append(images, img)
	}
	return images
}

func DecodeAt(storage *mjpeg.FrameStorage, index int) image.Image {
	frame := storage.GetAllPtr()[index]

	//try to read from the cache
	if frame.CachedImage == nil {
		img, err := jpeg.DecodeIntoRGBA(bytes.NewReader(frame.Body), &DecodeOptions)
		frame.CachedImage = img

		if err != nil {
			log.Printf("Unable to decode jpg: %v\n", err.Error())
			// Create a fallback black image
			rectangle := image.Rectangle{Min: image.Point{}, Max: image.Point{X: 640, Y: 360}}
			return image.NewRGBA(rectangle)
		}

		// update the width and height of the storage
		width, height := storage.GetImageSize()
		if img.Bounds().Dx() != width || img.Bounds().Dy() != height {
			storage.SetImageSize(img.Bounds().Dx(), img.Bounds().Dy())
		}

		return img
	} else {
		return frame.CachedImage
	}
}

func Decode(storage *mjpeg.FrameStorage) image.Image {
	return DecodeAt(storage, 0)
}

func Encode(image image.Image) *mjpeg.MjpegFrame {
	buff := bytes.NewBuffer([]byte{})
	config := EncodingOptions
	config.Quality = global.Config.EncodeQuality
	err := jpeg.Encode(buff, image, &config)

	if err != nil {
		log.Printf("Unable to encode image: %v\n", err.Error())
		// create a new MJPEGframe as fallback
		frame := mjpeg.NewMJPEGFrame()
		return &frame
	}

	return &mjpeg.MjpegFrame{Body: buff.Bytes()}
}

func Panel(layout PanelLayout, startIndex int, storages ...*mjpeg.FrameStorage) *mjpeg.MjpegFrame {
	//todo optimize by caching the result and painting over it?
	var images = DecodeAll(storages...)

	firstWidthInitial, firstHeightInitial := storages[0].GetImageSize()

	totalWidth := int(float64(firstWidthInitial) / layout.FirstWidth)
	totalHeight := int(float64(firstHeightInitial) / layout.FirstHeight)

	totalWidth, totalHeight = GetFinalImageSize(totalWidth, totalHeight)

	//output img
	imageOut := getImageOut(totalWidth, totalHeight)

	//main character image
	sp := image.Point{}
	r := image.Rectangle{Min: sp, Max: image.Point{X: int(float64(totalWidth) * layout.FirstWidth), Y: int(float64(totalHeight) * layout.FirstHeight)}}

	imageIn := ResizeOutputFrame(images[startIndex], r.Dx(), r.Dy())
	draw.NearestNeighbor.Scale(imageOut, r, imageIn, imageIn.Bounds(), draw.Over, nil)
	addLabel(imageOut, sp.X, sp.Y, global.Config.InputConfigs[startIndex].Label)

	for i, child := range layout.ChildrenPositions {
		if i+1 >= len(storages) {
			break
		}
		index := (startIndex + i + 1) % len(storages)
		sp = image.Point{X: int(float64(totalWidth) * child.X), Y: int(float64(totalHeight) * child.Y)}
		delta := image.Point{X: int(float64(totalWidth) * layout.ChildrenWidth), Y: int(float64(totalHeight) * layout.ChildrenHeight)}
		r = image.Rectangle{Min: sp, Max: sp.Add(delta)}

		imageIn := ResizeOutputFrame(images[index], r.Dx(), r.Dy())
		draw.NearestNeighbor.Scale(imageOut, r, imageIn, imageIn.Bounds(), draw.Over, nil)

		if global.Config.ShowInputLabel {
			// add offset to avoid overlap with the borders
			offsetW := 0
			offsetH := 0

			//labels to the left / on top don't need an offset
			if global.Config.ShowBorder && child.X != 0 {
				offsetW = getBorderSize(totalWidth) / 2
			}
			if global.Config.ShowBorder && child.Y != 0 {
				offsetH = getBorderSize(totalWidth) / 2
			}
			addLabel(imageOut, sp.X+offsetW, sp.Y+offsetH, global.Config.InputConfigs[index].Label)
		}
	}

	// draw border
	if global.Config.ShowBorder {
		border := getBorderSize(totalWidth)

		//vertical lines
		for _, el := range layout.VerticalBorderPoints {
			// calculate width and height
			rectangleWidth := border
			rectangleHeight := int((el.T2.Y*float64(totalHeight))-(el.T1.Y*float64(totalHeight))) + border/2
			// start point is t1 - borderWidth
			sp := image.Point{X: int(el.T1.X*float64(totalWidth)) - border/2, Y: int(el.T1.Y * float64(totalHeight))}
			borderVertical := image.Rectangle{
				Min: sp,
				Max: sp.Add(image.Point{X: rectangleWidth, Y: rectangleHeight}),
			}
			draw.Draw(imageOut, borderVertical, image.Black, image.Point{}, draw.Src)

		}

		//horizontal lines
		for _, el := range layout.HorizontalBorderPoints {
			// calculate width and height
			rectangleWidth := int((el.T2.X*float64(totalWidth))-(el.T1.X*float64(totalWidth))) + border/2
			rectangleHeight := border
			// start point is t1 - borderHeight
			sp := image.Point{X: int(el.T1.X*float64(totalWidth)) - border/2, Y: int(el.T1.Y * float64(totalHeight))}
			borderVertical := image.Rectangle{
				Min: sp,
				Max: sp.Add(image.Point{X: rectangleWidth, Y: rectangleHeight}),
			}
			draw.Draw(imageOut, borderVertical, image.Black, image.Point{}, draw.Src)
		}
	}

	return Encode(imageOut)
}

var imageOut *image.RGBA

func getImageOut(width int, height int) *image.RGBA {
	if imageOut == nil || imageOut.Rect.Max.X != width || imageOut.Rect.Max.Y != height {
		pointMax := image.Point{X: width, Y: height}
		rectangle := image.Rectangle{Min: image.Point{}, Max: pointMax}
		imageOut = image.NewRGBA(rectangle)
	}
	return imageOut
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
	var images = DecodeAll(storages...)
	firstWidthInitial, firstHeightInitial := storages[0].GetImageSize()
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
			images[i] = ResizeOutputFrame(images[i], cellWidth, cellHeight)
			//frames[i].CachedImage = images[i]
		}
	}

	// rectangle
	imageOut := getImageOut(totalWidth, totalHeight)

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
				x += getBorderSize(totalWidth) / 2
			}
			if row_ != 0 {
				y += getBorderSize(totalWidth) / 2
			}

			addLabel(imageOut, x, y, global.Config.InputConfigs[i].Label)

		}
	}

	// draw border
	if global.Config.ShowBorder {
		border := getBorderSize(totalWidth)
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

//ResizeOutputFrame resizes an image with regard to letterbox
func ResizeOutputFrame(img image.Image, width int, height int) image.Image {
	if !global.Config.IgnoreAspectRatio {

		deltaW := float64(width) / float64(img.Bounds().Dx())
		deltaH := float64(height) / float64(img.Bounds().Dy())

		factor := math.Min(deltaW, deltaH)

		outputWidth := int(float64(img.Bounds().Dx()) * factor)
		outputHeight := int(float64(img.Bounds().Dy()) * factor)

		offsetW := 0
		offsetH := 0
		if deltaW > deltaH {
			//letterbox left and right
			offsetW = (width - outputWidth) / 2

		} else if deltaW < deltaH {
			//letterbox top and bottom
			offsetH = (height - outputHeight) / 2
		}

		dst := image.NewRGBA(image.Rect(0, 0, width, height))
		draw.NearestNeighbor.Scale(dst, image.Rect(offsetW, offsetH, width-offsetW, height-offsetH), img, img.Bounds(), draw.Over, nil)
		return dst
	} else {
		dst := image.NewRGBA(image.Rect(0, 0, width, height))
		draw.NearestNeighbor.Scale(dst, dst.Rect, img, img.Bounds(), draw.Over, nil)
		return dst
	}
}

//Resize resizes an image
func Resize(img image.Image, width int, height int) image.Image {

	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.NearestNeighbor.Scale(dst, dst.Rect, img, img.Bounds(), draw.Over, nil)
	return dst
}

//GetFinalImageSize returns the size of the image in regard to the globally set Width and MinWidth
func GetFinalImageSize(currentWidth int, currentHeight int) (int, int) {
	if global.Config.Width != -1 && global.Config.Height != -1 {
		// both dimensions may be resized to the desired size
		currentWidth = global.Config.Width
		currentHeight = global.Config.Height
	} else if global.Config.Width != -1 && global.Config.Height == -1 {
		// reduce width and scale height
		currentHeight = int(float64(currentHeight) * (float64(global.Config.Width) / float64(currentWidth)))
		currentWidth = global.Config.Width
	} else if global.Config.Width == -1 && global.Config.Height != -1 {
		// reduce height and scale width
		currentWidth = int(float64(currentWidth) * (float64(global.Config.Height) / float64(currentHeight)))
		currentHeight = global.Config.Height
	}

	return currentWidth, currentHeight
}

func Carousel(input *mjpeg.FrameStorage, index int) *mjpeg.MjpegFrame {
	// with default settings (no resize, no quality change) the image doesn't need to be decoded and encoded at all
	if !global.DecodingNecessary() {
		return input.GetLatestPtr()
	}

	img := Decode(input)
	width, height := input.GetImageSize()

	// resizing
	if (global.Config.Width != width) || (global.Config.Height != height) {
		outputWidth, outputHeight := GetFinalImageSize(width, height)
		img = ResizeOutputFrame(img, outputWidth, outputHeight)
	}

	// label
	if global.Config.ShowInputLabel {
		addLabel(img.(*image.RGBA), 0, 0, global.Config.InputConfigs[index].Label)
	}

	// quality gets adjusted here
	return Encode(img)

}

var labelSrc = image.NewUniform(color.RGBA{R: 255, G: 255, B: 255, A: 255})
var f, _ = freetype.ParseFont(arial)

func addLabel(img *image.RGBA, x, y int, label string) {
	padding := global.Config.InputLabelFontSize / 4
	point := fixed.Point26_6{X: fixed.I(x + padding), Y: fixed.I(y + global.Config.InputLabelFontSize)}

	face := truetype.NewFace(f, &truetype.Options{
		Size: float64(global.Config.InputLabelFontSize),
	})
	d := &font.Drawer{
		Dst:  img,
		Src:  labelSrc,
		Face: face,
		Dot:  point,
	}

	textWidth := d.MeasureString(label).Round() + 2*padding
	textHeight := global.Config.InputLabelFontSize + padding

	// draw the black border...
	r := image.Rectangle{Min: image.Point{X: x, Y: y}, Max: image.Point{X: x + textWidth, Y: y + textHeight}}
	draw.Draw(img, r, image.Black, image.Point{}, draw.Src)

	// and the text
	d.DrawString(label)
}

const borderFactor = 0.0025

func getBorderSize(totalWidth int) int {
	if global.Config.ShowBorder {
		return utils.Max(int(float64(totalWidth)*borderFactor), 2)
	}
	return 0
}

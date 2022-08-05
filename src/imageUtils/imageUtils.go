package imageUtils

// Reference: https://go.dev/blog/image-draw

import (
	"bytes"
	_ "embed"
	"github.com/TobiasKunzFHNW/go-libjpeg/jpeg"
	"image"
	"image/color"
	"log"
	"math"
	"mjpeg_multiplexer/src/global"
	"mjpeg_multiplexer/src/mjpeg"
	"mjpeg_multiplexer/src/utils"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

//go:embed arial.ttf
var arial []byte

var imageOut *image.RGBA
var imageInContainers []*image.RGBA

var DecodeOptions = jpeg.DecoderOptions{ScaleTarget: image.Rectangle{}, DCTMethod: jpeg.DCTIFast, DisableFancyUpsampling: true, DisableBlockSmoothing: true}
var EncodingOptions = jpeg.EncoderOptions{Quality: 80, OptimizeCoding: false, ProgressiveMode: false, DCTMethod: jpeg.DCTIFast}

// Decode decodes an mjpeg image using the DecodeOptions (if not already decoded), caches the result and return the image
func Decode(storage *mjpeg.FrameStorage) image.Image {
	return decode(storage, nil)
}

// DecodeContainer is a more optimized version of Decode, which doesn't allocate a new image.Image instance
func DecodeContainer(storage *mjpeg.FrameStorage, imageContainer *image.RGBA) image.Image {
	return decode(storage, imageContainer)
}

func decode(storage *mjpeg.FrameStorage, imageContainer *image.RGBA) image.Image {
	frame := storage.GetFrame()

	//try to read from the cache
	if frame.CachedImage == nil {
		var img *image.RGBA
		var err error

		//try to use the optimized version if an imageContainer is provided
		if imageContainer != nil {
			err = jpeg.DecodeIntoRGBA2(&imageContainer, bytes.NewReader(frame.Body), &DecodeOptions)
			img = imageContainer
		} else {
			img, err = jpeg.DecodeIntoRGBA(bytes.NewReader(frame.Body), &DecodeOptions)
		}
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

// Encode encodes the image using the EncodingOptions
func Encode(image image.Image) *mjpeg.MjpegFrame {
	buff := bytes.NewBuffer([]byte{})
	config := EncodingOptions
	if global.Config.EncodeQuality != -1 {
		config.Quality = global.Config.EncodeQuality
	}
	err := jpeg.Encode(buff, image, &config)

	if err != nil {
		log.Printf("Unable to encode image: %v\n", err.Error())
		// create a new MJPEGframe as fallback
		frame := mjpeg.NewMJPEGFrame()
		return &frame
	}

	return &mjpeg.MjpegFrame{Body: buff.Bytes()}
}

// getImageOut creates / returns the reference to the image used to draw the results into
func getImageOut(width int, height int) *image.RGBA {
	if imageOut == nil || imageOut.Rect.Max.X != width || imageOut.Rect.Max.Y != height {
		pointMax := image.Point{X: width, Y: height}
		rectangle := image.Rectangle{Min: image.Point{}, Max: pointMax}
		imageOut = image.NewRGBA(rectangle)
	}
	return imageOut
}

// ResizeOutputFrame resizes an image with regard to letterbox
func ResizeOutputFrame(img image.Image, width int, height int) image.Image {
	if !global.Config.IgnoreAspectRatio {
		//resizing using letterbox, this means we have to draw the black borders left + right or top + bottom

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
		// just resize the image
		return Resize(img, width, height)
	}
}

// Resize resizes an image
func Resize(img image.Image, width int, height int) image.Image {
	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.NearestNeighbor.Scale(dst, dst.Rect, img, img.Bounds(), draw.Src, nil)
	return dst
}

// GetFinalImageSize returns the size of the image in regard to the globally set Width and MinWidth
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

var labelSrc = image.NewUniform(color.RGBA{R: 255, G: 255, B: 255, A: 255})
var f, _ = freetype.ParseFont(arial)

// addLabel places the given label string onto the image at location x,y
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

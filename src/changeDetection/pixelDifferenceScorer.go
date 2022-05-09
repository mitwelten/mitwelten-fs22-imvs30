package changeDetection

import (
	"bytes"
	"github.com/pixiv/go-libjpeg/jpeg"
	"image"
	"image/color"
	"math"
	"mjpeg_multiplexer/src/mjpeg"
)

// PixelDifferenceScorer simple pixel scorer struct
type PixelDifferenceScorer struct{}

var DecodeOptions = jpeg.DecoderOptions{ScaleTarget: image.Rectangle{}, DCTMethod: jpeg.DCTIFast, DisableFancyUpsampling: false, DisableBlockSmoothing: false}
var EncodingOptions = jpeg.EncoderOptions{Quality: 100, OptimizeCoding: false, ProgressiveMode: false, DCTMethod: jpeg.DCTISlow}

// Score implements Scorer.Score Method
func (s *PixelDifferenceScorer) Diff(frames []mjpeg.MjpegFrame) mjpeg.MjpegFrame {
	if len(frames) < 2 {
		return mjpeg.MjpegFrame{Body: mjpeg.Init()}
	}
	img1, _ := jpeg.Decode(bytes.NewReader(frames[0].Body), &DecodeOptions)
	img2, _ := jpeg.Decode(bytes.NewReader(frames[1].Body), &DecodeOptions)

	//return kernelDiff(img1, img2) / (img1.Bounds().Dx() * img1.Bounds().Dy())
	img := kernelPixelChangedThresholdImg(img1, img2)
	buff := bytes.NewBuffer([]byte{})
	//options := jpeg.Options{Quality: 100}
	_ = jpeg.Encode(buff, img, &EncodingOptions)
	return mjpeg.MjpegFrame{buff.Bytes()}
}

// Score implements Scorer.Score Method
func (s *PixelDifferenceScorer) Score(frames []mjpeg.MjpegFrame) int {
	if len(frames) < 2 {
		return -1
	}
	img1, _ := jpeg.Decode(bytes.NewReader(frames[0].Body), &DecodeOptions)
	img2, _ := jpeg.Decode(bytes.NewReader(frames[1].Body), &DecodeOptions)

	//return kernelDiff(img1, img2) / (img1.Bounds().Dx() * img1.Bounds().Dy())
	return kernelPixelChangedThreshold(img1, img2)
}

func kernelPixelChangedThresholdImg(img1 image.Image, img2 image.Image) image.Image {
	img := image.NewGray(img1.Bounds())

	for y := 0; y < img1.Bounds().Dy(); y++ {
		for x := 0; x < img1.Bounds().Dx(); x++ {
			r1, g1, b1, _ := img1.At(x, y).RGBA()
			v1 := (int(r1/255) + int(g1/255) + int(b1/255)) / 3
			r2, g2, b2, _ := img2.At(x, y).RGBA()
			v2 := (int(r2/255) + int(g2/255) + int(b2/255)) / 3

			//img.Set(x, y, color.Gray{Y: uint8(v1)})
			if math.Abs(float64(v1-v2)) > 20 {
				img.Set(x, y, color.White)
			} else {
				img.Set(x, y, color.Black)
			}

		}
	}
	return img
}
func kernelPixelChangedThreshold(img1 image.Image, img2 image.Image) int {
	score := 0
	for y := 0; y < img1.Bounds().Dy(); y++ {
		for x := 0; x < img1.Bounds().Dx(); x++ {
			r1, g1, b1, _ := img1.At(x, y).RGBA()
			v1 := (int(r1/255) + int(g1/255) + int(b1/255)) / 3
			r2, g2, b2, _ := img2.At(x, y).RGBA()
			v2 := (int(r2/255) + int(g2/255) + int(b2/255)) / 3

			if math.Abs(float64(v1-v2)) > 20 {
				score++
			}
		}
	}
	return score
}

// kernelDiff returns kernel Difference of two images as int values
// difference = sum of all greyscale pixel differences
func kernelDiff(img1 image.Image, img2 image.Image) int {
	score := 0
	for y := 0; y < img1.Bounds().Dy(); y++ {
		for x := 0; x < img1.Bounds().Dx(); x++ {
			r1, g1, b1, _ := img1.At(x, y).RGBA()
			v1 := int(r1) + int(g1) + int(b1)
			r2, g2, b2, _ := img2.At(x, y).RGBA()
			v2 := int(r2) + int(g2) + int(b2)
			d := uint8(math.Abs(float64(v1 - v2)))
			score += int(d)
		}
	}
	return score
}

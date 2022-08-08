package activityDetection

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"io/ioutil"
	"math"
	"mjpeg_multiplexer/src/imageUtils"
	"strconv"
	"time"
)

const threshold = 20

// FrameDifferenceScore compares the two images by comparing  pixel and counting how many pixels are different.
// Different means that the grayscale value has a difference of `threshold` or more.
// Img1 and Img2 must have the same size
func FrameDifferenceScore(img1 image.Image, img2 image.Image) float64 {
	nPixels := float64(img1.Bounds().Dx() * img1.Bounds().Dy())
	return float64(kernelPixelChangedThreshold(img1, img2)) / nPixels
}

// kernelPixelChangedThreshold calculates the score by comparing the images pixel by pixel
// returns the number of pixel which have a difference of lambda or more
func kernelPixelChangedThreshold(img1 image.Image, img2 image.Image) int {
	score := 0

	for y := 0; y < img1.Bounds().Dy(); y++ {
		for x := 0; x < img1.Bounds().Dx(); x++ {
			r1, g1, b1, _ := img1.At(x, y).RGBA()
			v1 := (int(r1/255) + int(g1/255) + int(b1/255)) / 3
			r2, g2, b2, _ := img2.At(x, y).RGBA()
			v2 := (int(r2/255) + int(g2/255) + int(b2/255)) / 3

			if math.Abs(float64(v1-v2)) >= threshold {
				score++
			}
		}
	}
	return score
}

// todo remove rest down below
const imageSize = 128

func FrameDifferenceImage(img1 image.Image, img2 image.Image) image.Image {
	ratio := float64(img1.Bounds().Dy()) / float64(img1.Bounds().Dx())
	height := int(imageSize * ratio)
	img1 = imageUtils.Resize(img1, imageSize, height)
	img2 = imageUtils.Resize(img2, imageSize, height)

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	buff := bytes.NewBuffer([]byte{})
	_ = jpeg.Encode(buff, img1, nil)
	_ = ioutil.WriteFile("temp/"+timestamp+"img1.benchmark_jpeg_libraries", buff.Bytes(), 0644)

	buff = bytes.NewBuffer([]byte{})
	_ = jpeg.Encode(buff, img2, nil)
	_ = ioutil.WriteFile("temp/"+timestamp+"img2.benchmark_jpeg_libraries", buff.Bytes(), 0644)

	img3 := kernelPixelChangedThresholdImage(img1, img2)

	buff = bytes.NewBuffer([]byte{})
	_ = jpeg.Encode(buff, img3, nil)
	_ = ioutil.WriteFile("temp/"+timestamp+"img3.benchmark_jpeg_libraries", buff.Bytes(), 0644)

	return img3
}

func kernelPixelChangedThresholdImage(img1 image.Image, img2 image.Image) image.Image {
	img := image.NewGray(img1.Bounds())
	for y := 0; y < img1.Bounds().Dy(); y++ {
		for x := 0; x < img1.Bounds().Dx(); x++ {
			r1, g1, b1, _ := img1.At(x, y).RGBA()
			v1 := (int(r1/255) + int(g1/255) + int(b1/255)) / 3
			r2, g2, b2, _ := img2.At(x, y).RGBA()
			v2 := (int(r2/255) + int(g2/255) + int(b2/255)) / 3

			if math.Abs(float64(v1-v2)) >= threshold {
				img.Set(x, y, color.White)
			} else {
				img.Set(x, y, color.Black)
			}
		}
	}
	return img
}

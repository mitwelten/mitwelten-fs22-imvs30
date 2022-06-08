package motionDetection

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
const imageSize = 100

// FrameDifferenceScore implements Scorer.FrameDifferenceScore Method
func FrameDifferenceScore(img1 image.Image, img2 image.Image) float64 {
	img1 = imageUtils.Resize(img1, imageSize, imageSize)
	img2 = imageUtils.Resize(img2, imageSize, imageSize)

	nPixels := float64(img1.Bounds().Dx() * img1.Bounds().Dy())
	return kernelPixelChangedThreshold(img1, img2) / nPixels
}

func kernelPixelChangedThreshold(img1 image.Image, img2 image.Image) float64 {
	score := 0.0

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

func FrameDifferenceImage(img1 image.Image, img2 image.Image) image.Image {
	/*	img1 = imageUtils.Resize(img1, 100, 100)
		img2 = imageUtils.Resize(img2, 100, 100)
	*/
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	buff := bytes.NewBuffer([]byte{})
	_ = jpeg.Encode(buff, img1, nil)
	_ = ioutil.WriteFile("temp/"+timestamp+"img1.jpg", buff.Bytes(), 0644)

	buff = bytes.NewBuffer([]byte{})
	_ = jpeg.Encode(buff, img2, nil)
	_ = ioutil.WriteFile("temp/"+timestamp+"img2.jpg", buff.Bytes(), 0644)

	img3 := kernelPixelChangedThresholdImage(img1, img2)

	buff = bytes.NewBuffer([]byte{})
	_ = jpeg.Encode(buff, img3, nil)
	_ = ioutil.WriteFile("temp/"+timestamp+"img3.jpg", buff.Bytes(), 0644)

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

package activityDetection

import (
	"image"
	"math"
)

const threshold = 20
const imageSize = 128

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

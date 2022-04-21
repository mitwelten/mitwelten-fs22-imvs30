package changeDetection

import (
	"bytes"
	"image"
	"math"
	"mjpeg_multiplexer/src/mjpeg"
)

// PixelDifferenceScorer simple pixel scorer struct
type PixelDifferenceScorer struct{}

// Score implements Scorer.Score Method
func (s *PixelDifferenceScorer) Score(frames []mjpeg.MjpegFrame) int {
	if len(frames) < 2 {
		return -1
	}
	img1, _, _ := image.Decode(bytes.NewReader(frames[0].Body))
	img2, _, _ := image.Decode(bytes.NewReader(frames[1].Body))

	return kernelDiff(img1, img2) / (img1.Bounds().Dx() * img1.Bounds().Dy())
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

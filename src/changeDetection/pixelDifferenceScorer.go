package changeDetection

import (
	"github.com/anthonynsimon/bild/blend"
	"github.com/anthonynsimon/bild/segment"
	"image"
	"image/color"
	"math"
	"mjpeg_multiplexer/src/imageUtils"
	"mjpeg_multiplexer/src/mjpeg"
)

// PixelDifferenceScorer simple pixel scorer struct
type PixelDifferenceScorer struct{}

// Score implements Scorer.Score Method
func (s *PixelDifferenceScorer) Diff(frames []mjpeg.MjpegFrame) mjpeg.MjpegFrame {
	if len(frames) < 2 {
		return mjpeg.MjpegFrame{Body: mjpeg.Init()}
	}
	img1 := imageUtils.Decode(frames[0])
	img2 := imageUtils.Decode(frames[1])

	img := kernelPixelChangedThresholdImg(img1, img2)

	return imageUtils.Encode(img)
}

// Score implements Scorer.Score Method
func (s *PixelDifferenceScorer) Score(frames []mjpeg.MjpegFrame) float64 {
	if len(frames) < 2 {
		return -1
	}
	img1 := imageUtils.Decode(frames[0])
	img2 := imageUtils.Decode(frames[1])

	return kernelPixelChangedThreshold(img1, img2) / float64(img1.Bounds().Dx()*img1.Bounds().Dy())
}

func kernelPixelChangedThresholdImg(img1 image.Image, img2 image.Image) image.Image {
	//img := image.NewGray(img1.Bounds())

	//img1 = blur.Box(img1, 1)
	//img2 = blur.Box(img2, 1)

	imgOut := blend.Difference(img1, img2)
	return segment.Threshold(imgOut, 20)

	/*	for y := 0; y < img1.Bounds().Dy(); y++ {
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
	*/
}

//kernelPixelChangedThresholdImgEigenblur pixelates the image before comparing the images
func kernelPixelChangedThresholdImgEigenblur(img1 image.Image, img2 image.Image) image.Image {
	img := image.NewGray(img1.Bounds())

	radius := 2

	for y_ := radius; y_ < img1.Bounds().Dy()-radius; y_ += 1 + 2*radius {
		for x_ := radius; x_ < img1.Bounds().Dx()-radius; x_ += 1 + 2*radius {

			v1 := 0
			v2 := 0

			for dy := -radius; dy <= radius; dy++ {
				for dx := -radius; dx <= radius; dx++ {
					x := x_ + dx
					y := y_ + dy

					r1, g1, b1, _ := img1.At(x, y).RGBA()
					v1 += (int(r1/255) + int(g1/255) + int(b1/255)) / 3
					r2, g2, b2, _ := img2.At(x, y).RGBA()
					v2 += (int(r2/255) + int(g2/255) + int(b2/255)) / 3
				}
			}
			v1 = v1 / ((1 + 2*radius) * (1 + 2*radius))
			v2 = v2 / ((1 + 2*radius) * (1 + 2*radius))

			for dy := -radius; dy <= radius; dy++ {
				for dx := -radius; dx <= radius; dx++ {
					x := x_ + dx
					y := y_ + dy

					//		img.Set(x, y, color.Gray{uint8(v1)})

					if math.Abs(float64(v1-v2)) > 20 {
						img.Set(x, y, color.White)
					} else {
						img.Set(x, y, color.Black)
					}

				}
			}

		}
	}

	return img
}

//kernelPixelChangedThresholdImgEigenblur pixelates the image before comparing the images
func kernelPixelChangedThresholdEigenblur(img1 image.Image, img2 image.Image) int {
	score := 0

	radius := 2

	for y_ := radius; y_ < img1.Bounds().Dy()-radius; y_ += 1 + 2*radius {
		for x_ := radius; x_ < img1.Bounds().Dx()-radius; x_ += 1 + 2*radius {

			v1 := 0
			v2 := 0

			for dy := -radius; dy <= radius; dy++ {
				for dx := -radius; dx <= radius; dx++ {
					x := x_ + dx
					y := y_ + dy

					r1, g1, b1, _ := img1.At(x, y).RGBA()
					v1 += (int(r1/255) + int(g1/255) + int(b1/255)) / 3
					r2, g2, b2, _ := img2.At(x, y).RGBA()
					v2 += (int(r2/255) + int(g2/255) + int(b2/255)) / 3
				}
			}
			v1 = v1 / ((1 + 2*radius) * (1 + 2*radius))
			v2 = v2 / ((1 + 2*radius) * (1 + 2*radius))

			if math.Abs(float64(v1-v2)) > 20 {
				score += 1
			}

		}
	}

	return score
}

func kernelPixelChangedThreshold(img1 image.Image, img2 image.Image) float64 {
	score := 0.0

	//img1 = blur.Box(img1, 1)

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

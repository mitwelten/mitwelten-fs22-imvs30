package main

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"io/ioutil"
	"math"
	"strconv"
	"time"
)

var threshold float64 = 20

func main() {
	b1, err := ioutil.ReadFile("misc/imagediff/img1.jpg")
	if err != nil {
		panic("yikes")
	}

	img1, _ := jpeg.Decode(bytes.NewReader(b1))

	b2, _ := ioutil.ReadFile("misc/imagediff/img2.jpg")
	img2, _ := jpeg.Decode(bytes.NewReader(b2))

	for i := 0; i <= 255; i++ {
		threshold = float64(i)
		FrameDifferenceImage(img1, img2)
	}

}

func FrameDifferenceImage(img1 image.Image, img2 image.Image) image.Image {
	/*	img1 = imageUtils.Resize(img1, 100, 100)
		img2 = imageUtils.Resize(img2, 100, 100)
	*/

	/*	buff := bytes.NewBuffer([]byte{})
		_ = jpeg.Encode(buff, img1, nil)
		_ = ioutil.WriteFile("temp/"+timestamp+"img1.jpg", buff.Bytes(), 0644)

		buff = bytes.NewBuffer([]byte{})
		_ = jpeg.Encode(buff, img2, nil)
		_ = ioutil.WriteFile("temp/"+timestamp+"img2.jpg", buff.Bytes(), 0644)
	*/
	img3 := kernelPixelChangedThresholdImage(img1, img2)

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	buff := bytes.NewBuffer([]byte{})
	_ = jpeg.Encode(buff, img3, nil)
	_ = ioutil.WriteFile("misc/imagediff/"+timestamp+"_"+strconv.Itoa(int(threshold))+".jpg", buff.Bytes(), 0644)

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

			if math.Abs(float64(v1-v2)) > threshold {
				img.Set(x, y, color.White)
			} else {
				img.Set(x, y, color.Black)
			}
		}
	}
	return img
}

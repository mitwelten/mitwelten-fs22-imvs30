package image

// Reference: https://go.dev/blog/image-draw

import (
	"bytes"
	"github.com/anthonynsimon/bild/blur"
	"github.com/anthonynsimon/bild/transform"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"math"
	"mjpeg_multiplexer/src/mjpeg"
)

func GetImgDiffScore(frame1 mjpeg.MjpegFrame, frame2 mjpeg.MjpegFrame) int {
	img1_, _, _ := image.Decode(bytes.NewReader(frame1.Body))
	img1 := blur.Gaussian(img1_, 0)
	img2_, _, _ := image.Decode(bytes.NewReader(frame2.Body))
	img2 := blur.Gaussian(img2_, 0)

	score := 0

	for y := 0; y < img1.Bounds().Dy(); y++ {
		for x := 0; x < img1.Bounds().Dx(); x++ {

			r1, g1, b1, _ := img1.At(x, y).RGBA()
			v1 := (int(r1) + int(g1) + int(b1))

			r2, g2, b2, _ := img2.At(x, y).RGBA()
			v2 := (int(r2) + int(g2) + int(b2))

			d := int(math.Abs(float64(v1 - v2)))

			score = score + d
		}
	}

	return score
}
func kernelSquareOfColors(img1 image.Image, img2 image.Image, out *image.Gray) int {
	score := 0
	for y := 0; y < img1.Bounds().Dy(); y++ {
		for x := 0; x < img1.Bounds().Dx(); x++ {
			r1, g1, b1, _ := img1.At(x, y).RGBA()
			r2, g2, b2, _ := img2.At(x, y).RGBA()
			v1 := math.Abs(float64(r1 - r2))
			v2 := math.Abs(float64(g1 - g2))
			v3 := math.Abs(float64(b1 - b2))
			d := v1*v1 + v2*v2 + v3*v3
			score += int(d)
			out.Set(x, y, color.Gray{Y: uint8(d)})
		}
	}
	return score
}

func kernelDiff(img1 image.Image, img2 image.Image, out *image.Gray) int {
	score := 0
	for y := 0; y < img1.Bounds().Dy(); y++ {
		for x := 0; x < img1.Bounds().Dx(); x++ {
			r1, g1, b1, _ := img1.At(x, y).RGBA()
			v1 := int(r1) + int(g1) + int(b1)
			r2, g2, b2, _ := img2.At(x, y).RGBA()
			v2 := int(r2) + int(g2) + int(b2)
			d := uint8(math.Abs(float64(v1 - v2)))
			score += int(d)
			out.Set(x, y, color.Gray{Y: d})
		}
	}
	return score
}

func GetImgDiffResize(frame1 mjpeg.MjpegFrame, frame2 mjpeg.MjpegFrame) mjpeg.MjpegFrame {
	img1_, _, _ := image.Decode(bytes.NewReader(frame1.Body))
	img2_, _, _ := image.Decode(bytes.NewReader(frame2.Body))

	img1 := transform.Resize(img1_, img1_.Bounds().Dx()/8, img1_.Bounds().Dy()/8, transform.Box)
	img2 := transform.Resize(img2_, img2_.Bounds().Dx()/8, img2_.Bounds().Dy()/8, transform.Box)

	/*	img1 := transform.Resize(img1_, 8, 8, transform.Box)
		img2 := transform.Resize(img2_, 8, 8, transform.Box)
	*/
	out := image.NewGray(image.Rect(0, 0, img1.Bounds().Dx()/8, img1.Bounds().Dy()/8))

	score := kernelSquareOfColors(img1, img2, out)
	println(score)

	inSmall := transform.Resize(img1, img1_.Bounds().Dx()/8, img1_.Bounds().Dy()/8, transform.Box)
	inBig := transform.Resize(inSmall, img1_.Bounds().Dx(), img1_.Bounds().Dy(), transform.Box)
	_ = transform.Resize(out, img1_.Bounds().Dx(), img1_.Bounds().Dy(), transform.Box)

	buff := bytes.NewBuffer([]byte{})
	options := jpeg.Options{Quality: 100}
	err := jpeg.Encode(buff, inBig, &options)

	if err != nil {
		panic("can't encode jpeg")
	}
	return mjpeg.MjpegFrame{Body: buff.Bytes()}
}

func GetImgDiff(frame1 mjpeg.MjpegFrame, frame2 mjpeg.MjpegFrame) mjpeg.MjpegFrame {
	img1, _, _ := image.Decode(bytes.NewReader(frame1.Body))
	img2, _, _ := image.Decode(bytes.NewReader(frame2.Body))
	out := image.NewGray(image.Rect(0, 0, img1.Bounds().Dx(), img1.Bounds().Dy()))

	score := kernelDiff(img1, img2, out)
	println(score)

	buff := bytes.NewBuffer([]byte{})
	options := jpeg.Options{Quality: 100}
	err := jpeg.Encode(buff, out, &options)
	if err != nil {
		panic("can't encode jpeg")
	}
	return mjpeg.MjpegFrame{Body: buff.Bytes()}
}

func decode(frames []mjpeg.MjpegFrame) []image.Image {
	var images []image.Image
	for _, frame := range frames {
		var img, _, _ = image.Decode(bytes.NewReader(frame.Body))
		images = append(images, img)
	}
	return images
}

func Grid(row int, col int, frames ...mjpeg.MjpegFrame) mjpeg.MjpegFrame {
	var nCells = row * col
	var nFrames = len(frames)

	if nFrames > nCells {
		panic("Too many frames")
	}

	if nFrames == 0 {
		panic("At least one frame needed")
	}
	var images = decode(frames)

	// rectangle
	var i0 = images[0]
	var pointMax = image.Point{X: i0.Bounds().Dx() * col, Y: i0.Bounds().Dy() * row}
	var rectangle = image.Rectangle{Min: image.Point{}, Max: pointMax}

	// image
	var imageOut = image.NewRGBA(rectangle)

	for i := 0; i < nCells; i++ {
		var row_ = i / col
		var col_ = i % col

		if i >= nFrames {
			break
		}

		var sp = image.Point{X: i0.Bounds().Dx() * col_, Y: i0.Bounds().Dy() * row_}
		var r = image.Rectangle{Min: sp, Max: sp.Add(images[i].Bounds().Size())}
		draw.Draw(imageOut, r, images[i], image.Point{}, draw.Src)
	}

	buff := bytes.NewBuffer([]byte{})
	options := jpeg.Options{Quality: 100}
	err := jpeg.Encode(buff, imageOut, &options)
	if err != nil {
		panic("can't encode jpeg")
	}
	return mjpeg.MjpegFrame{Body: buff.Bytes()}
}

func MergeImages(frame1 mjpeg.MjpegFrame, frame2 mjpeg.MjpegFrame) (frame mjpeg.MjpegFrame) {
	var img1, _, _ = image.Decode(bytes.NewReader(frame1.Body))
	var img2, _, _ = image.Decode(bytes.NewReader(frame2.Body))

	//starting position of the second image (bottom left)
	sp2 := image.Point{X: img1.Bounds().Dx()}

	//new rectangle for the second image
	r2 := image.Rectangle{Min: sp2, Max: sp2.Add(img2.Bounds().Size())}

	//rectangle for the big image
	r := image.Rectangle{Min: image.Point{}, Max: r2.Max}

	rgba := image.NewRGBA(r)

	draw.Draw(rgba, img1.Bounds(), img1, image.Point{}, draw.Src)
	draw.Draw(rgba, r2, img2, image.Point{}, draw.Src)

	//out, customErrors := os.Create("./output.jpg")
	//if customErrors != nil {
	//	fmt.Println(customErrors)
	//}

	//var opt jpeg.Options
	//opt.Quality = 100
	buff := bytes.NewBuffer([]byte{})
	//jpeg.Encode(out, rgba, &opt)
	err := jpeg.Encode(buff, rgba, nil)
	if err != nil {
		panic("can't encode jpeg")
	}

	return mjpeg.MjpegFrame{Body: buff.Bytes()}
}

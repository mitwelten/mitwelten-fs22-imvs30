package image

// Reference: https://go.dev/blog/image-draw

import (
	"bytes"
	"image"
	"image/draw"
	"image/jpeg"
	"mjpeg_multiplexer/src/mjpeg"
)

func decode(frames []mjpeg.Frame) []image.Image {
	var images []image.Image
	for _, frame := range frames {
		var img, _, _ = image.Decode(bytes.NewReader(frame.Body))
		images = append(images, img)
	}
	return images
}

func Grid(row int, col int, frames ...mjpeg.Frame) mjpeg.Frame {
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
	err := jpeg.Encode(buff, imageOut, nil)
	if err != nil {
		panic("can't encode jpeg")
	}
	return mjpeg.Frame{Body: buff.Bytes()}
}

func MergeImages(frame1 mjpeg.Frame, frame2 mjpeg.Frame) (frame mjpeg.Frame) {
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

	return mjpeg.Frame{Body: buff.Bytes()}
}

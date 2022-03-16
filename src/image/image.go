package image

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
		if i >= nFrames {
			break
		}

		var sp = image.Point{X: i0.Bounds().Dx() * col, Y: i0.Bounds().Dy() * row}
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

func Grid4(frames ...mjpeg.Frame) mjpeg.Frame {
	if len(frames) != 4 {
		panic("must be at least 4 frames")
	}

	var images = decode(frames)
	var i0 = images[0]
	var sp1 = image.Point{X: i0.Bounds().Dx()}
	var sp2 = image.Point{Y: i0.Bounds().Dy()}
	var sp3 = image.Point{X: i0.Bounds().Dx(), Y: i0.Bounds().Dy()}

	//rectangle for the 4 grid
	var pointMax = image.Point{X: i0.Bounds().Dx() * 2, Y: i0.Bounds().Dy() * 2}
	var r = image.Rectangle{Min: image.Point{}, Max: pointMax}

	//create the image
	var imageOut = image.NewRGBA(r)

	//fill the grid
	draw.Draw(imageOut, images[0].Bounds(), images[0], image.Point{}, draw.Src)

	var r1 = image.Rectangle{Min: sp1, Max: sp1.Add(images[1].Bounds().Size())}
	draw.Draw(imageOut, r1, images[1], image.Point{}, draw.Src)

	var r2 = image.Rectangle{Min: sp2, Max: sp2.Add(images[2].Bounds().Size())}
	draw.Draw(imageOut, r2, images[2], image.Point{}, draw.Src)

	var r3 = image.Rectangle{Min: sp3, Max: sp3.Add(images[3].Bounds().Size())}
	draw.Draw(imageOut, r3, images[3], image.Point{}, draw.Src)

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

	//out, err := os.Create("./output.jpg")
	//if err != nil {
	//	fmt.Println(err)
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

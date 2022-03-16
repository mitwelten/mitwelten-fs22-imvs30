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

func Grid4(frames []mjpeg.Frame) mjpeg.Frame {
	if len(frames) != 4 {
		panic("must be 4 frames")
	}

	var images = decode(frames)
	var i0 = images[0]
	//rectangle for the 4 grid
	var pointMax = image.Point{X: i0.Bounds().Dx() * 2, Y: i0.Bounds().Dy() * 2}
	var r = image.Rectangle{Min: image.Point{}, Max: pointMax}

	//create the image
	var imageOut = image.NewRGBA(r)

	//fill the grid
	draw.Draw(imageOut, images[0].Bounds(), images[0], image.Point{}, draw.Src)
	draw.Draw(imageOut, images[1].Bounds(), images[1], image.Point{X: i0.Bounds().Dx()}, draw.Src)
	draw.Draw(imageOut, images[2].Bounds(), images[2], image.Point{Y: i0.Bounds().Dy()}, draw.Src)
	draw.Draw(imageOut, images[3].Bounds(), images[3], image.Point{X: i0.Bounds().Dx(), Y: i0.Bounds().Dy()}, draw.Src)

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

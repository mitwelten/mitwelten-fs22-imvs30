package image

import (
	"bytes"
	"image"
	"image/draw"
	"image/jpeg"
	"mjpeg_multiplexer/src/mjpeg"
)

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

	draw.Draw(rgba, img1.Bounds(), img1, image.Point{0, 0}, draw.Src)
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

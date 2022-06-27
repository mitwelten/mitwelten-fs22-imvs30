package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"github.com/pixiv/go-libjpeg/jpeg"
	"golang.org/x/image/draw"
	"image"
	"io/ioutil"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"
)

//go:embed image.jpg
var imageData []byte
var DecodeOptions = jpeg.DecoderOptions{ScaleTarget: image.Rectangle{}, DCTMethod: jpeg.DCTIFast, DisableFancyUpsampling: false, DisableBlockSmoothing: false}
var EncodingOptions = jpeg.EncoderOptions{Quality: 100, OptimizeCoding: true, ProgressiveMode: false, DCTMethod: jpeg.DCTISlow}

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6061", nil))
	}()

	img, _ := jpeg.Decode(bytes.NewReader(imageData), &DecodeOptions)

	images := []image.Image{img, img, img, img}
	nCols := 2
	nRows := 2
	totalWidth := img.Bounds().Dx() * nCols
	totalHeight := img.Bounds().Dy() * nRows
	var cellWidth int
	var cellHeight int

	cellWidth = totalWidth / nCols
	cellHeight = totalHeight / nRows

	pointMax := image.Point{X: totalWidth, Y: totalHeight}
	rectangle := image.Rectangle{Min: image.Point{}, Max: pointMax}

	nIterations := 20

	imageOut := image.NewRGBA(rectangle)
	fmt.Printf("Classic: ")
	start := time.Now()
	for i := 0; i < nIterations; i++ {
		classic(nRows, nCols, cellWidth, cellHeight, imageOut, images)
	}
	end := time.Since(start).Milliseconds()
	save(imageOut, "classic.jpg")
	fmt.Printf("Total %vms (%vms per iteration)\n", end, end/int64(nIterations))

	img, _ = jpeg.DecodeIntoRGBA(bytes.NewReader(imageData), &DecodeOptions)
	images = []image.Image{img, img, img, img}
	imageOut = image.NewRGBA(rectangle)
	fmt.Printf("Classic RGBA: ")
	start = time.Now()
	for i := 0; i < nIterations; i++ {
		classic(nRows, nCols, cellWidth, cellHeight, imageOut, images)
	}
	end = time.Since(start).Milliseconds()
	save(imageOut, "classic_RGBA.jpg")
	fmt.Printf("Total %vms (%vms per iteration)\n", end, end/int64(nIterations))

	imageOut = image.NewRGBA(rectangle)
	fmt.Printf("eigenbau RGBA: ")
	start = time.Now()
	for i := 0; i < nIterations; i++ {
		eigenbau(nRows, nCols, cellWidth, cellHeight, imageOut, images)
	}
	end = time.Since(start).Milliseconds()
	save(imageOut, "eigenbau_RGBA.jpg")
	fmt.Printf("Total %vms (%vms per iteration)\n", end, end/int64(nIterations))

}

func classic(nRows int, nCols int, cellWidth int, cellHeight int, imageOut *image.RGBA, images []image.Image) {
	for i := 0; i < nCols*nRows; i++ {
		var row_ = i / nCols
		var col_ = i % nCols

		var sp = image.Point{X: cellWidth * col_, Y: cellHeight * row_}
		//grid lines:

		var r = image.Rectangle{Min: sp, Max: sp.Add(images[i].Bounds().Size())}
		draw.Draw(imageOut, r, images[i], image.Point{}, draw.Src)
	}
}
func eigenbau(nRows int, nCols int, cellWidth int, cellHeight int, imageOut *image.RGBA, images []image.Image) {
	outStride := imageOut.Bounds().Dx() * 4

	for i := 0; i < len(images); i++ {
		var row_ = i / nCols
		var col_ = i % nCols
		img := images[i]
		cellStride := img.Bounds().Dx() * 4

		offsetX := cellWidth * col_
		offsetY := cellHeight * row_

		imgStride := img.Bounds().Dx() * 4
		img_ := img.(*image.RGBA)
		/*		imgStride := img.Bounds().Dx() * 4
				img_ := image.NewRGBA(img.Bounds())
				draw.Draw(img_, img.Bounds(), img, image.Point{}, draw.Src)
		*/
		for cellY := 0; cellY < img.Bounds().Dy(); cellY++ {
			cellIndex := cellStride * cellY
			outIndex := (outStride * (cellY + offsetY)) + (4 * (offsetX))
			copy(imageOut.Pix[outIndex:], img_.Pix[cellIndex:cellIndex+imgStride])
		}
	}
}

func save(imageOut *image.RGBA, name string) {
	//encode
	buff := bytes.NewBuffer([]byte{})
	err := jpeg.Encode(buff, imageOut, &EncodingOptions)

	if err != nil {
		panic("can't encode jpg")
	}
	err = ioutil.WriteFile("jpg/grid/"+name, buff.Bytes(), 0644)
}

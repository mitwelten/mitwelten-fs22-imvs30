package main

import (
	"bytes"
	_ "embed"
	"fmt"
	govips "github.com/davidbyttow/govips/v2/vips"
	"github.com/pixiv/go-libjpeg/jpeg"
	"golang.org/x/image/draw"
	"image"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"
)

var DecodeOptions = jpeg.DecoderOptions{ScaleTarget: image.Rectangle{}, DCTMethod: jpeg.DCTIFast, DisableFancyUpsampling: true, DisableBlockSmoothing: true}
var EncodingOptions = jpeg.EncoderOptions{Quality: 100, OptimizeCoding: false, ProgressiveMode: false, DCTMethod: jpeg.DCTISlow}

//go:embed image_real.jpg
var imageData []byte

var outWidth = 1600.0
var outHeight = 600.0

func vips_merge(iterations int, encode bool) {
	govips.Startup(nil)
	defer govips.Shutdown()

	fmt.Printf("    Number of runs: %v\n", iterations)
	start := time.Now()

	for i := 0; i < iterations; i++ {

		imgout, _ := govips.NewImageFromBuffer(imageData)
		w := float64(imgout.Width())
		h := float64(imgout.Height())
		imgout.ResizeWithVScale(outWidth/w, outHeight/h, govips.KernelNearest)
		img1, _ := govips.NewThumbnailWithSizeFromBuffer(imageData, 800, 600, govips.InterestingNone, govips.SizeForce)
		img2, _ := govips.NewThumbnailWithSizeFromBuffer(imageData, 800, 600, govips.InterestingNone, govips.SizeForce)

		imgout.Insert(img1, 0, 0, false, nil)
		imgout.Insert(img2, img1.Width(), 0, false, nil)

		if encode {

			ep := govips.NewJpegExportParams()
			ep.StripMetadata = true
			ep.Quality = 75
			ep.Interlace = true
			ep.OptimizeCoding = true
			ep.SubsampleMode = govips.VipsForeignSubsampleAuto
			ep.TrellisQuant = true
			ep.OvershootDeringing = true
			ep.OptimizeScans = true
			ep.QuantTable = 3

			//_, _, _ = imgout.ExportJpeg(ep)
			//_, _ = imgout.ToBytes()

			/*	p := govips.NewDefaultExportParams()
				p.Quality = 100
				_, _ = imgout.ToImage(p)
				//ioutil.WriteFile("go_vips_merge.jpg", goImg, 0644)

				//buff := bytes.NewBuffer([]byte{})
				//_ = jpeg.Encode(buff, imageOut, &EncodingOptions)*/
		}
	}

	end := time.Since(start).Milliseconds()
	fmt.Printf("    Total: %v ms\n", end)
	fmt.Printf("    Per iteration: %v ms\n", float64(end)/float64(iterations))
}

func baumarkt_merge(iterations int, encode bool) {
	fmt.Printf("    Number of runs: %v\n", iterations)
	start := time.Now()

	for i := 0; i < iterations; i++ {
		img1, _ := jpeg.Decode(bytes.NewReader(imageData), &DecodeOptions)
		img2, _ := jpeg.Decode(bytes.NewReader(imageData), &DecodeOptions)

		image1 := image.NewRGBA(image.Rect(0, 0, 800, 600))
		image2 := image.NewRGBA(image.Rect(0, 0, 800, 600))

		draw.NearestNeighbor.Scale(image1, image1.Rect, img1, img1.Bounds(), draw.Over, nil)
		draw.NearestNeighbor.Scale(image2, image1.Rect, img2, img2.Bounds(), draw.Over, nil)
		images := []image.Image{image1, image2}

		var i0 = image1
		var pointMax = image.Point{X: i0.Bounds().Dx() * 2, Y: i0.Bounds().Dy() * 1}
		var rectangle = image.Rectangle{Min: image.Point{}, Max: pointMax}

		// image
		var imageOut = image.NewRGBA(rectangle)

		nCells := 2
		col := 2
		nFrames := 2
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

		if encode {
			buff := bytes.NewBuffer([]byte{})
			//options := jpeg.EncoderOptions{Quality: 100}
			_ = jpeg.Encode(buff, imageOut, &EncodingOptions)
			//_ = ioutil.WriteFile("baumarkt_merge.jpg", buff.Bytes(), 0644)
		}

	}
	end := time.Since(start).Milliseconds()
	fmt.Printf("    Total: %v ms\n", end)
	fmt.Printf("    Per iteration: %v ms\n", float64(end)/float64(iterations))

}
func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	iterations := 50
	govips.LoggingSettings(nil, govips.LogLevelCritical)

	fmt.Printf("vips\n")
	vips_merge(iterations, true)

	fmt.Printf("baumarkt\n")
	baumarkt_merge(iterations, true)

	// =========================

}

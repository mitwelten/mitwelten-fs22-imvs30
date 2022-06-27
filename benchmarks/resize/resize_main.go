package main

import (
	"bytes"
	_ "embed"
	"fmt"
	govips "github.com/davidbyttow/govips/v2/vips"
	"github.com/discord/lilliput"
	"github.com/pixiv/go-libjpeg/jpeg"
	"github.com/vipsimage/vips"
	"golang.org/x/image/draw"
	"io/ioutil"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	//"golang.org/x/image/draw"
	"image"
	"time"
)

//go:embed image_real.jpg
var imageData []byte
var DecodeOptions = jpeg.DecoderOptions{ScaleTarget: image.Rectangle{}, DCTMethod: jpeg.DCTIFast, DisableFancyUpsampling: true, DisableBlockSmoothing: true}
var EncodingOptions = jpeg.EncoderOptions{Quality: 90, OptimizeCoding: false, ProgressiveMode: false, DCTMethod: jpeg.DCTISlow}

func resizeLilliPut(origName, newName string) {
	decoder, _ := lilliput.NewDecoder(imageData)

	ops := lilliput.NewImageOps(8192)
	defer ops.Close()

	// create a buffer to store the output image, 50MB in this case
	outputImg := make([]byte, 20*1024*1024)

	opts := &lilliput.ImageOptions{
		FileType:             filepath.Ext(newName),
		Width:                410,
		Height:               308,
		ResizeMethod:         lilliput.ImageOpsFit,
		NormalizeOrientation: true,
		EncodeOptions:        map[int]int{lilliput.JpegQuality: 100},
	}

	// resize and transcode image
	_, _ = ops.Transform(decoder, opts, outputImg)

	// image has been resized, now write file out

	//ioutil.WriteFile("lilliput.jpg", b, 0644)
}

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	/*	img := vips.NewFromBuffer(imageData, "")
		_ = img.ThumbnailImage(410)

		img.MagickSaveBuffer()
		img.Wrap()
		img.Save2file("out.png")
		iterations := 100
	*/
	iterations := 100

	fmt.Printf("baumarkt resize\n")
	fmt.Printf("    Number of runs: %v\n", iterations)
	opt := DecodeOptions
	src, _ := jpeg.Decode(bytes.NewReader(imageData), &opt)
	start := time.Now()

	for i := 0; i < iterations; i++ {
		dst := image.NewRGBA(image.Rect(0, 0, 100, 100))
		draw.NearestNeighbor.Scale(dst, dst.Rect, src, src.Bounds(), draw.Over, nil)
		if i == -1 {
			buff := bytes.NewBuffer([]byte{})
			_ = jpeg.Encode(buff, dst, &EncodingOptions)
			_ = ioutil.WriteFile("baumarkt_resize.jpg", buff.Bytes(), 0644)
		}
	}
	end := time.Since(start).Milliseconds()
	fmt.Printf("    Total: %v ms\n", end)
	fmt.Printf("    Per iteration: %v ms\n", float64(end)/float64(iterations))

	os.Exit(0)

	fmt.Printf("govips resize\n")
	fmt.Printf("    Number of runs: %v\n", iterations)
	start = time.Now()
	govips.LoggingSettings(nil, govips.LogLevelError)

	for i := 0; i < iterations; i++ {

		img, _ := govips.NewThumbnailWithSizeFromBuffer(imageData, 1920, 1080, govips.InterestingNone, govips.SizeForce)

		/*		p := govips.NewJpegExportParams()
				p.OptimizeCoding = false
				p.OptimizeScans = false
				p.StripMetadata = true
				p.Interlace = false
				_, _, _ = img.ExportJpeg(p)
		*/
		img.ExportJpeg(nil)
		//ioutil.WriteFile("govips.jpg", buff, 0644)
	}

	end = time.Since(start).Milliseconds()
	fmt.Printf("    Total: %v ms\n", end)
	fmt.Printf("    Per iteration: %v ms\n", float64(end)/float64(iterations))

	//---------------------

	fmt.Printf("bimg resize\n")
	fmt.Printf("    Number of runs: %v\n", iterations)
	start = time.Now()

	for i := 0; i < iterations; i++ {
		//_, _ = bimg.NewImage(imageData).ForceResize(410, 308)
		/*
			newImage, _ := bimg.NewImageFromBuffer(imageData).Resize(800, 600)
			bimg.Write("new.jpg", newImage)
		*/
	}
	end = time.Since(start).Milliseconds()
	fmt.Printf("    Total: %v ms\n", end)
	fmt.Printf("    Per iteration: %v ms\n", float64(end)/float64(iterations))

	fmt.Printf("libjpg resize\n")
	fmt.Printf("    Number of runs: %v\n", iterations)
	start = time.Now()

	for i := 0; i < iterations; i++ {
		opt := DecodeOptions
		opt.ScaleTarget = image.Rectangle{Max: image.Point{1, 1}}
		jpeg.Decode(bytes.NewReader(imageData), &opt)
	}

	end = time.Since(start).Milliseconds()
	fmt.Printf("    Total: %v ms\n", end)
	fmt.Printf("    Per iteration: %v ms\n", float64(end)/float64(iterations))

	fmt.Printf("baumarkt resize\n")
	fmt.Printf("    Number of runs: %v\n", iterations)
	opt = DecodeOptions
	src, _ = jpeg.Decode(bytes.NewReader(imageData), &opt)
	start = time.Now()

	for i := 0; i < iterations; i++ {
		dst := image.NewRGBA(image.Rect(0, 0, 100, 100))
		draw.NearestNeighbor.Scale(dst, dst.Rect, src, src.Bounds(), draw.Over, nil)
	}

	end = time.Since(start).Milliseconds()
	fmt.Printf("    Total: %v ms\n", end)
	fmt.Printf("    Per iteration: %v ms\n", float64(end)/float64(iterations))

	fmt.Printf("vips resize\n")
	fmt.Printf("    Number of runs: %v\n", iterations)
	start = time.Now()

	for i := 0; i < iterations; i++ {
		img := vips.NewFromBuffer(imageData, "")
		img.ThumbnailImage(410)
		img.JPEGSave("vips1_out.jpg")

	}

	end = time.Since(start).Milliseconds()
	fmt.Printf("    Total: %v ms\n", end)
	fmt.Printf("    Per iteration: %v ms\n", float64(end)/float64(iterations))

	fmt.Printf("lilliput resize\n")
	fmt.Printf("    Number of runs: %v\n", iterations)
	start = time.Now()

	for i := 0; i < iterations; i++ {
		resizeLilliPut("image_real.jpg", "lilli.jpg")
	}

	end = time.Since(start).Milliseconds()
	fmt.Printf("    Total: %v ms\n", end)
	fmt.Printf("    Per iteration: %v ms\n", float64(end)/float64(iterations))

}

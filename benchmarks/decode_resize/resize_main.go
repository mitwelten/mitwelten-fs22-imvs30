package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"github.com/h2non/bimg/v2"
	"github.com/pixiv/go-libjpeg/jpeg"
	"golang.org/x/image/draw"
	"image"
	"io/ioutil"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"
)

//go:embed image_real.jpg
var imageData []byte
var DecodeOptions = jpeg.DecoderOptions{ScaleTarget: image.Rectangle{}, DCTMethod: jpeg.DCTIFast, DisableFancyUpsampling: true, DisableBlockSmoothing: true}
var EncodingOptions = jpeg.EncoderOptions{Quality: 70, OptimizeCoding: false, ProgressiveMode: false, DCTMethod: jpeg.DCTISlow}

var targetWidth = 800
var targetHeight = 600

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	encode := true
	iterations := 10
	inputs := 1

	fmt.Printf("Number of runs: %v, inputs %v\n", iterations, inputs)

	fmt.Printf("libjpg\n")
	start := time.Now()

	for i := 0; i < iterations; i++ {
		opt := DecodeOptions
		opt.ScaleTarget = image.Rectangle{Max: image.Point{targetWidth, targetHeight}}
		//initial inaccruate resize

		var src image.Image
		var dst *image.RGBA
		for j := 0; j < inputs; j++ {
			src, _ = jpeg.Decode(bytes.NewReader(imageData), &opt)
			// and exact resize
			dst = image.NewRGBA(image.Rect(0, 0, targetWidth, targetWidth))
			draw.NearestNeighbor.Scale(dst, dst.Rect, src, src.Bounds(), draw.Over, nil)
		}

		if encode {
			buff := bytes.NewBuffer([]byte{})
			_ = jpeg.Encode(buff, dst, &EncodingOptions)
			if i == 0 {
				_ = ioutil.WriteFile("benchmark_decoderesize_libjpeg.jpg", buff.Bytes(), 0644)
			}
		}
	}

	end := time.Since(start).Milliseconds()
	fmt.Printf("    Total: %v ms\n", end)
	fmt.Printf("    Per iteration: %v ms\n", float64(end)/float64(iterations))

	//bimg
	fmt.Printf("bimg resize\n")
	fmt.Printf("    Number of runs: %v\n", iterations)
	start = time.Now()

	for i := 0; i < iterations; i++ {
		//_, _ = bimg.NewImage(imageData).ForceResize(410, 308)
		var img *bimg.Image
		var err error
		for j := 0; j < inputs; j++ {
			img, err = bimg.NewImageFromBuffer(imageData)
			if err = img.Resize(bimg.ResizeOptions{Width: targetWidth, Height: targetWidth}); err != nil {
				panic(err)
			}
		}
		if encode {
			b, _ := img.Save(bimg.SaveOptions{Type: bimg.JPEG, Compression: 1, Quality: 70, Speed: 1, Interlace: false, StripMetadata: true, Lossless: false})
			if i == 0 {
				_ = ioutil.WriteFile("benchmark_decoderesize_bimg.jpg", b, 0644)
			}
		}
	}
	end = time.Since(start).Milliseconds()
	fmt.Printf("    Total: %v ms\n", end)
	fmt.Printf("    Per iteration: %v ms\n", float64(end)/float64(iterations))

	// Eigenbau
	fmt.Printf("Marke Eigenbau\n")
	start = time.Now()

	for i := 0; i < iterations; i++ {
		opt := DecodeOptions
		var src image.Image
		var dst *image.RGBA

		for j := 0; j < inputs; j++ {
			src, _ = jpeg.Decode(bytes.NewReader(imageData), &opt)
			dst = image.NewRGBA(image.Rect(0, 0, targetWidth, targetWidth))
			draw.NearestNeighbor.Scale(dst, dst.Rect, src, src.Bounds(), draw.Over, nil)
		}

		if encode {
			buff := bytes.NewBuffer([]byte{})
			_ = jpeg.Encode(buff, dst, &EncodingOptions)
			if i == 0 {
				_ = ioutil.WriteFile("benchmark_decoderesize_eigenbau.jpg", buff.Bytes(), 0644)
			}

		}
	}

	end = time.Since(start).Milliseconds()
	fmt.Printf("    Total: %v ms\n", end)
	fmt.Printf("    Per iteration: %v ms\n", float64(end)/float64(iterations))
}

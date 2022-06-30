package JPEGBenchmark_libjpg

import (
	"github.com/pixiv/go-libjpeg/jpeg"
	"io/ioutil"
	"strconv"

	"bytes"
	_ "embed"
	"fmt"
	"image"
	"time"
)

//go:embed image_orig.jpg
var imageData []byte

var DecodeOptions = jpeg.DecoderOptions{ScaleTarget: image.Rectangle{}, DCTMethod: jpeg.DCTIFast, DisableFancyUpsampling: true, DisableBlockSmoothing: true}
var EncodingOptions = jpeg.EncoderOptions{Quality: 80, OptimizeCoding: false, ProgressiveMode: false, DCTMethod: jpeg.DCTIFast}

func Encode(iterations int) {
	img, _ := jpeg.Decode(bytes.NewReader(imageData), &DecodeOptions)
	fmt.Printf("  Encode\n")
	fmt.Printf("    Number of runs: %v\n", iterations)
	start := time.Now()
	for i := 0; i < iterations; i++ {
		//encode
		buff := bytes.NewBuffer([]byte{})
		err := jpeg.Encode(buff, img, &EncodingOptions)

		if err != nil {
			panic("can't encode jpg")
		}
	}
	end := time.Since(start).Milliseconds()
	fmt.Printf("    Total: %v ms\n", end)
	fmt.Printf("    Per iteration: %v ms\n", float64(end)/float64(iterations))
}
func EncodeRGBA(iterations int) {
	img, _ := jpeg.DecodeIntoRGBA(bytes.NewReader(imageData), &DecodeOptions)
	fmt.Printf("  Encode RGBA\n")
	fmt.Printf("    Number of runs: %v\n", iterations)
	start := time.Now()
	for i := 0; i < iterations; i++ {
		//encode
		buff := bytes.NewBuffer([]byte{})
		err := jpeg.Encode(buff, img, &EncodingOptions)

		if err != nil {
			panic("can't encode jpg")
		}
	}
	end := time.Since(start).Milliseconds()
	fmt.Printf("    Total: %v ms\n", end)
	fmt.Printf("    Per iteration: %v ms\n", float64(end)/float64(iterations))
}

func Decode(iterations int) {
	fmt.Printf("  Decode\n")
	fmt.Printf("    Number of runs: %v\n", iterations)
	start := time.Now()
	for i := 0; i < iterations; i++ {
		_, _ = jpeg.Decode(bytes.NewReader(imageData), &DecodeOptions)
	}
	end := time.Since(start).Milliseconds()
	fmt.Printf("    Total: %v ms\n", end)
	fmt.Printf("    Per iteration: %v ms\n", float64(end)/float64(iterations))
}

func DecodeRGBA(iterations int) {
	fmt.Printf("  Decode RGBA\n")
	fmt.Printf("    Number of runs: %v\n", iterations)
	start := time.Now()
	for i := 0; i < iterations; i++ {
		_, _ = jpeg.DecodeIntoRGBA(bytes.NewReader(imageData), &DecodeOptions)
	}
	end := time.Since(start).Milliseconds()
	fmt.Printf("    Total: %v ms\n", end)
	fmt.Printf("    Per iteration: %v ms\n", float64(end)/float64(iterations))
}
func DecodeEncode(iterations int) {
	fmt.Printf("  DecodeEncode\n")
	fmt.Printf("    Number of runs: %v\n", iterations)
	start := time.Now()
	for i := 0; i < iterations; i++ {
		//decode
		img, _ := jpeg.Decode(bytes.NewReader(imageData), &DecodeOptions)

		//encode
		buff := bytes.NewBuffer([]byte{})
		//options := jpeg.EncoderOptions{Quality: 100}
		err := jpeg.Encode(buff, img, &EncodingOptions)

		if err != nil {
			panic("can't encode jpg")
		}
	}
	end := time.Since(start).Milliseconds()
	fmt.Printf("    Total: %v ms\n", end)
	fmt.Printf("    Per iteration: %v ms\n", float64(end)/float64(iterations))
}

func DecodeEncodeScaled() {
	opt := jpeg.DecoderOptions{ScaleTarget: image.Rectangle{}, DCTMethod: jpeg.DCTIFast, DisableFancyUpsampling: true, DisableBlockSmoothing: true}
	opt.ScaleTarget = image.Rectangle{Min: image.Point{10, 0}, Max: image.Point{X: 1, Y: 2}}
	//decode
	img, _ := jpeg.Decode(bytes.NewReader(imageData), &opt)

	//encode
	buff := bytes.NewBuffer([]byte{})
	//options := jpeg.EncoderOptions{Quality: 100}
	_ = jpeg.Encode(buff, img, &EncodingOptions)
	_ = ioutil.WriteFile("jpg/jpgBenchmark_libjpg/image_out_libjpg.jpg", buff.Bytes(), 0644)

}

func DecodeEncodeSave() []byte {
	img, _ := jpeg.Decode(bytes.NewReader(imageData), &DecodeOptions)

	//encode
	buff := bytes.NewBuffer([]byte{})
	err := jpeg.Encode(buff, img, &EncodingOptions)

	if err != nil {
		panic("can't encode jpg")
	}
	err = ioutil.WriteFile("jpg/jpgBenchmark_libjpg/quality_"+strconv.Itoa(EncodingOptions.Quality)+".jpg", buff.Bytes(), 0644)
	if err != nil {
		panic("can't save jpg")
	}
	return buff.Bytes()
}

func BenchmarkDecode(iterations int, imageData []byte) float64 {
	fmt.Printf("  Decode libjpg\n")
	fmt.Printf("    Number of runs: %v\n", iterations)

	for i := 0; i < iterations/10; i++ {
		reader := bytes.NewReader(imageData)
		_, _ = jpeg.Decode(reader, &DecodeOptions)
	}

	start := time.Now()
	for i := 0; i < iterations; i++ {
		reader := bytes.NewReader(imageData)
		_, _ = jpeg.Decode(reader, &DecodeOptions)
	}
	end := time.Since(start).Milliseconds()
	fmt.Printf("    Total: %v ms\n", end)
	fmt.Printf("    Per iteration: %v ms\n", float64(end)/float64(iterations))
	return float64(end) / float64(iterations)
}
func BenchmarkEncode(iterations int, imageData []byte) float64 {
	fmt.Printf("  Encode libjpg\n")
	fmt.Printf("    Number of runs: %v\n", iterations)
	reader := bytes.NewReader(imageData)
	img, _ := jpeg.Decode(reader, &DecodeOptions)

	for i := 0; i < iterations/10; i++ {
		buff := bytes.NewBuffer([]byte{})
		err := jpeg.Encode(buff, img, &EncodingOptions)

		if err != nil {
			panic("can't encode jpg")
		}
	}

	start := time.Now()
	for i := 0; i < iterations; i++ {
		buff := bytes.NewBuffer([]byte{})
		err := jpeg.Encode(buff, img, &EncodingOptions)

		if err != nil {
			panic("can't encode jpg")
		}
	}
	end := time.Since(start).Milliseconds()
	fmt.Printf("    Total: %v ms\n", end)
	fmt.Printf("    Per iteration: %v ms\n", float64(end)/float64(iterations))
	return float64(end) / float64(iterations)
}

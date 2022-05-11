package JPEGBenchmark_libjpg

import (
	"github.com/pixiv/go-libjpeg/jpeg"
	"io/ioutil"

	"bytes"
	_ "embed"
	"fmt"
	"image"
	"time"
)

//go:embed image_big.jpg

var imageData []byte

var DecodeOptions = jpeg.DecoderOptions{ScaleTarget: image.Rectangle{}, DCTMethod: jpeg.DCTIFast, DisableFancyUpsampling: true, DisableBlockSmoothing: true}

// OptimizeCoding: Slightly more efficient compreesion, but way slower
// ProgressiveMode not needed in our case
// DCTMethod JDCT_ISLOW is the fastest on my (Tobi) system
// See https://github.com/libjpeg-turbo/libjpeg-turbo/blob/main/libjpeg.txt
var EncodingOptions = jpeg.EncoderOptions{Quality: 100, OptimizeCoding: false, ProgressiveMode: false, DCTMethod: jpeg.DCTISlow}

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
	err = ioutil.WriteFile("jpg/jpgBenchmark_libjpg/image_out_libjpg.jpg", buff.Bytes(), 0644)
	if err != nil {
		panic("can't save jpg")
	}
	return buff.Bytes()
}

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

//go:embed image.jpg

var imageData []byte

// OptimizeCoding: Slightly more efficient compreesion, but way slower
// ProgressiveMode not needed in our case
// DCTMethod JDCT_ISLOW is the fastest on my (Tobi) system
// See https://github.com/libjpeg-turbo/libjpeg-turbo/blob/main/libjpeg.txt
var EncodingOptions = jpeg.EncoderOptions{Quality: 100, OptimizeCoding: false, ProgressiveMode: false, DCTMethod: jpeg.DCTISlow}

func Encode(iterations int) {
	img, _, _ := image.Decode(bytes.NewReader(imageData))
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
		_, _, _ = image.Decode(bytes.NewReader(imageData))
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
		img, _, _ := image.Decode(bytes.NewReader(imageData))

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

func DecodeEncodeSave() []byte {
	img, _, _ := image.Decode(bytes.NewReader(imageData))

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

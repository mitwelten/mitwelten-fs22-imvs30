package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"image"
	"image/jpeg"
	"time"
)

//go:embed image.jpg
var imageData []byte

func main() {
	iterations := 250
	Decode(iterations)
	Encode(iterations)
	DecodeEncode(iterations)
}

func Encode(iterations int) {
	img, _, _ := image.Decode(bytes.NewReader(imageData))
	fmt.Printf("Encode\n")
	fmt.Printf("    Number of runs: %v\n", iterations)
	start := time.Now()
	for i := 0; i < iterations; i++ {
		//encode
		buff := bytes.NewBuffer([]byte{})
		options := jpeg.Options{Quality: 100}
		err := jpeg.Encode(buff, img, &options)

		if err != nil {
			panic("can't encode jpeg")
		}
	}
	end := time.Since(start).Milliseconds()
	fmt.Printf("    Total: %v ms\n", end)
	fmt.Printf("    Per iteration: %v ms\n", float64(end)/float64(iterations))
}

func Decode(iterations int) {
	fmt.Printf("Decode\n")
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
	fmt.Printf("DecodeEncode\n")
	fmt.Printf("    Number of runs: %v\n", iterations)
	start := time.Now()
	for i := 0; i < iterations; i++ {
		//decode
		img, _, _ := image.Decode(bytes.NewReader(imageData))

		//encode
		buff := bytes.NewBuffer([]byte{})
		options := jpeg.Options{Quality: 100}
		err := jpeg.Encode(buff, img, &options)

		if err != nil {
			panic("can't encode jpeg")
		}
	}
	end := time.Since(start).Milliseconds()
	fmt.Printf("    Total: %v ms\n", end)
	fmt.Printf("    Per iteration: %v ms\n", float64(end)/float64(iterations))

}

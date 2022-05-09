package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"github.com/anthonynsimon/bild/blur"
	"github.com/esimov/stackblur-go"
	"github.com/matsuyoshi30/song2"
	"image"
	_ "image/jpeg"
	"time"
)

//go:embed image.jpg
var imageData []byte

func main() {
	iterations := 200
	radius := 1.0

	img, _, err := image.Decode(bytes.NewReader(imageData))

	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("Warmup...\n")
	for i := 0; i < 20; i++ {
		song2.GaussianBlur(img, radius)
		stackblur.Process(img, uint32(radius))
		blur.Gaussian(img, radius)
	}

	fmt.Printf("Song2: \n")
	start := time.Now()
	for i := 0; i < iterations; i++ {
		song2.GaussianBlur(img, radius)
	}
	end := time.Since(start).Milliseconds()
	fmt.Printf("    Total: %v ms\n", end)
	fmt.Printf("    Per iteration: %v ms\n", float64(end)/float64(iterations))

	fmt.Printf("Stackblur: \n")
	start = time.Now()
	for i := 0; i < iterations; i++ {
		stackblur.Process(img, uint32(radius))
	}
	end = time.Since(start).Milliseconds()
	fmt.Printf("    Total: %v ms\n", end)
	fmt.Printf("    Per iteration: %v ms\n", float64(end)/float64(iterations))

	fmt.Printf("Bild: \n")
	start = time.Now()
	for i := 0; i < iterations; i++ {
		blur.Gaussian(img, radius)
	}
	end = time.Since(start).Milliseconds()
	fmt.Printf("    Total: %v ms\n", end)
	fmt.Printf("    Per iteration: %v ms\n", float64(end)/float64(iterations))
}

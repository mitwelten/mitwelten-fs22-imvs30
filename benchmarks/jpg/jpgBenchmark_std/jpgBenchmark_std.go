package JPEGBenchmark_std

import (
	"bytes"
	_ "embed"
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
	"time"
)

//go:embed image_real.jpg
var imageData []byte

func Encode(iterations int) {
	img, _, _ := image.Decode(bytes.NewReader(imageData))
	fmt.Printf("  Encode\n")
	fmt.Printf("    Number of runs: %v\n", iterations)
	start := time.Now()
	for i := 0; i < iterations; i++ {
		//encode
		buff := bytes.NewBuffer([]byte{})
		options := jpeg.Options{Quality: 100}
		err := jpeg.Encode(buff, img, &options)

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
		options := jpeg.Options{Quality: 100}
		err := jpeg.Encode(buff, img, &options)

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
	options := jpeg.Options{Quality: 100}
	err := jpeg.Encode(buff, img, &options)

	if err != nil {
		panic("can't encode jpg")
	}
	err = ioutil.WriteFile("jpg/jpgBenchmark_std/image_out_std.jpg", buff.Bytes(), 0644)

	if err != nil {
		panic("can't save jpg")
	}
	return buff.Bytes()
}

func BenchmarkDecode(iterations int, imageData []byte) float64 {
	fmt.Printf("  Decode \n")
	fmt.Printf("    Number of runs: %v\n", iterations)

	//warmup
	for i := 0; i < iterations/10; i++ {
		reader := bytes.NewReader(imageData)
		_, _ = jpeg.Decode(reader)
	}

	start := time.Now()
	for i := 0; i < iterations; i++ {
		reader := bytes.NewReader(imageData)
		_, _ = jpeg.Decode(reader)
	}
	end := time.Since(start).Milliseconds()
	fmt.Printf("    Total: %v ms\n", end)
	fmt.Printf("    Per iteration: %v ms\n", float64(end)/float64(iterations))
	return float64(end) / float64(iterations)
}
func BenchmarkEncode(iterations int, imageData []byte) float64 {
	fmt.Printf("  Encode \n")
	fmt.Printf("    Number of runs: %v\n", iterations)
	op := jpeg.Options{Quality: 80}

	reader := bytes.NewReader(imageData)
	img, _ := jpeg.Decode(reader)

	for i := 0; i < iterations/10; i++ {
		buff := bytes.NewBuffer([]byte{})
		err := jpeg.Encode(buff, img, &op)
		if err != nil {
			panic("can't encode jpg")
		}
	}

	start := time.Now()
	for i := 0; i < iterations; i++ {
		buff := bytes.NewBuffer([]byte{})
		err := jpeg.Encode(buff, img, &op)
		if err != nil {
			panic("can't encode jpg")
		}
	}
	end := time.Since(start).Milliseconds()
	fmt.Printf("    Total: %v ms\n", end)
	fmt.Printf("    Per iteration: %v ms\n", float64(end)/float64(iterations))
	return float64(end) / float64(iterations)
}

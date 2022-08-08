package JPEGBenchmark_std

import (
	"bytes"
	_ "embed"
	"fmt"
	"image/jpeg"
	"time"
)

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
			panic("can't encode benchmark_jpeg_libraries")
		}
	}

	start := time.Now()
	for i := 0; i < iterations; i++ {
		buff := bytes.NewBuffer([]byte{})
		err := jpeg.Encode(buff, img, &op)
		if err != nil {
			panic("can't encode benchmark_jpeg_libraries")
		}
	}
	end := time.Since(start).Milliseconds()
	fmt.Printf("    Total: %v ms\n", end)
	fmt.Printf("    Per iteration: %v ms\n", float64(end)/float64(iterations))
	return float64(end) / float64(iterations)
}

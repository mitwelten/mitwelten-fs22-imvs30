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

	end := time.Duration(0)

	for i := 0; i < iterations; i++ {
		// calculate the time for decode + some operation
		start_ := time.Now()
		reader := bytes.NewReader(imageData)
		_, _ = jpeg.Decode(reader)
		end += time.Since(start_)

	}

	avg := float64(end.Milliseconds()) / float64(iterations)
	fmt.Printf("    Total: %v ms\n", end.Milliseconds())
	fmt.Printf("    Per iteration: %.2f ms\n", avg)
	return avg
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

	end := time.Duration(0)

	for i := 0; i < iterations; i++ {
		// calculate the time for decode + some operation
		start_ := time.Now()
		buff := bytes.NewBuffer([]byte{})
		err := jpeg.Encode(buff, img, &op)
		if err != nil {
			panic("can't encode benchmark_jpeg_libraries")
		}
		end += time.Since(start_)

	}

	avg := float64(end.Milliseconds()) / float64(iterations)
	fmt.Printf("    Total: %v ms\n", end.Milliseconds())
	fmt.Printf("    Per iteration: %.2f ms\n", avg)
	return avg
}

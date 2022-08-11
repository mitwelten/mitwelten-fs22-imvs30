package JPEGBenchmark_libjpg

import (
	"bytes"
	_ "embed"
	"fmt"
	"github.com/pixiv/go-libjpeg/jpeg"
	"image"
	"time"
)

var DecodeOptions = jpeg.DecoderOptions{ScaleTarget: image.Rectangle{}, DCTMethod: jpeg.DCTIFast, DisableFancyUpsampling: true, DisableBlockSmoothing: true}
var EncodingOptions = jpeg.EncoderOptions{Quality: 80, OptimizeCoding: false, ProgressiveMode: false, DCTMethod: jpeg.DCTIFast}

func BenchmarkDecode(iterations int, imageData []byte) float64 {
	fmt.Printf("  Decode libjpg\n")
	fmt.Printf("    Number of runs: %v\n", iterations)

	for i := 0; i < iterations/10; i++ {
		reader := bytes.NewReader(imageData)
		_, _ = jpeg.Decode(reader, &DecodeOptions)
	}

	end := time.Duration(0)

	for i := 0; i < iterations; i++ {
		// calculate the time for decode + some operation
		start_ := time.Now()
		reader := bytes.NewReader(imageData)
		_, _ = jpeg.Decode(reader, &DecodeOptions)
		end += time.Since(start_)

	}

	avg := float64(end.Milliseconds()) / float64(iterations)
	fmt.Printf("    Total: %v ms\n", end.Milliseconds())
	fmt.Printf("    Per iteration: %.2f ms\n", avg)
	return avg
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
			panic("can't encode benchmark_jpeg_libraries")
		}
	}

	end := time.Duration(0)
	for i := 0; i < iterations; i++ {
		// calculate the time for decode + some operation
		start_ := time.Now()
		buff := bytes.NewBuffer([]byte{})
		err := jpeg.Encode(buff, img, &EncodingOptions)

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

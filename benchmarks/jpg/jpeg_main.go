package main

import (
	jpgBenchmark_libjpg "benchmarks/jpg/jpgBenchmark_libjpg"
	jpgBenchmark_std "benchmarks/jpg/jpgBenchmark_std"
	"fmt"
	"github.com/pixiv/go-libjpeg/jpeg"
)

func main() {
	iterations := 500

	fmt.Printf("Standard libraray: \n")
	jpgBenchmark_std.Decode(iterations)
	jpgBenchmark_std.Encode(iterations)
	jpgBenchmark_std.DecodeEncode(iterations)
	//	_ = jpgBenchmark_std.DecodeEncodeSave()

	fmt.Printf("libjpg: \n")
	jpgBenchmark_libjpg.Decode(iterations)
	jpgBenchmark_libjpg.Encode(iterations)
	jpgBenchmark_libjpg.DecodeEncode(iterations)
	//	_ = jpgBenchmark_libjpg.DecodeEncodeSave()
}

func compareDCTMethods() {
	iterations := 500

	fmt.Printf("libjpg DCTIFast: \n")
	jpgBenchmark_libjpg.EncodingOptions.DCTMethod = jpeg.DCTIFast
	jpgBenchmark_libjpg.Encode(iterations)

	fmt.Printf("libjpg DCTFloat: \n")
	jpgBenchmark_libjpg.EncodingOptions.DCTMethod = jpeg.DCTFloat
	jpgBenchmark_libjpg.Encode(iterations)

	fmt.Printf("libjpg DCTISlow: \n")
	jpgBenchmark_libjpg.EncodingOptions.DCTMethod = jpeg.DCTISlow
	jpgBenchmark_libjpg.Encode(iterations)

}

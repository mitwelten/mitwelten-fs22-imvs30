package main

import (
	jpgBenchmark_libjpg "benchmarks/jpg/jpgBenchmark_libjpg"
	jpgBenchmark_std "benchmarks/jpg/jpgBenchmark_std"
	"fmt"
	"github.com/pixiv/go-libjpeg/jpeg"
)

func main() {

	iterations := 100
	fmt.Printf("libjpg: \n")
	//jpgBenchmark_libjpg.Decode(iterations)
	//jpgBenchmark_libjpg.Encode(iterations)

	jpgBenchmark_libjpg.DecodeRGBA(iterations)
	jpgBenchmark_libjpg.EncodeRGBA(iterations)

	fmt.Printf("----\n")

	fmt.Printf("libjpg: DisableFancyUpsampling \n")
	jpgBenchmark_libjpg.DecodeOptions.DisableFancyUpsampling = true
	jpgBenchmark_libjpg.DecodeRGBA(iterations)
	jpgBenchmark_libjpg.DecodeOptions.DisableFancyUpsampling = false

	fmt.Printf("libjpg: DisableBlockSmoothing \n")
	jpgBenchmark_libjpg.DecodeOptions.DisableBlockSmoothing = true
	jpgBenchmark_libjpg.DecodeRGBA(iterations)
	jpgBenchmark_libjpg.DecodeOptions.DisableBlockSmoothing = false

	fmt.Printf("libjpg: combined \n")
	jpgBenchmark_libjpg.DecodeOptions.DisableFancyUpsampling = true
	jpgBenchmark_libjpg.DecodeOptions.DisableBlockSmoothing = true
	jpgBenchmark_libjpg.DecodeRGBA(iterations)
	jpgBenchmark_libjpg.DecodeOptions.DisableFancyUpsampling = false
	jpgBenchmark_libjpg.DecodeOptions.DisableBlockSmoothing = false

	fmt.Printf("----\n")

	fmt.Printf("libjpg: OptimizeCoding \n")
	jpgBenchmark_libjpg.EncodingOptions.OptimizeCoding = true
	jpgBenchmark_libjpg.EncodeRGBA(iterations)
	jpgBenchmark_libjpg.EncodingOptions.OptimizeCoding = false

	fmt.Printf("libjpg: ProgressiveMode \n")
	jpgBenchmark_libjpg.EncodingOptions.ProgressiveMode = true
	jpgBenchmark_libjpg.EncodeRGBA(iterations)
	jpgBenchmark_libjpg.EncodingOptions.ProgressiveMode = false

	fmt.Printf("libjpg: combined \n")
	jpgBenchmark_libjpg.EncodingOptions.OptimizeCoding = true
	jpgBenchmark_libjpg.EncodingOptions.ProgressiveMode = true
	jpgBenchmark_libjpg.EncodeRGBA(iterations)
	jpgBenchmark_libjpg.EncodingOptions.OptimizeCoding = false
	jpgBenchmark_libjpg.EncodingOptions.ProgressiveMode = false

	//jpgBenchmark_libjpg.DecodeEncode(iterations)
	//	_ = jpgBenchmark_libjpg.DecodeEncodeSave()

	fmt.Printf("Standard libraray: \n")
	jpgBenchmark_std.Decode(iterations)
	jpgBenchmark_std.Encode(iterations)
	//jpgBenchmark_std.DecodeEncode(iterations)
	//	_ = jpgBenchmark_std.DecodeEncodeSave()

}

func compareDCTMethods() {
	iterations := 100
	jpgBenchmark_libjpg.Decode(20)

	fmt.Printf("libjpg DCTIFast: \n")
	jpgBenchmark_libjpg.EncodingOptions.DCTMethod = jpeg.DCTIFast
	jpgBenchmark_libjpg.DecodeOptions.DCTMethod = jpeg.DCTIFast
	//jpgBenchmark_libjpg.Encode(iterations)
	jpgBenchmark_libjpg.Decode(iterations)

	fmt.Printf("libjpg DCTFloat: \n")
	jpgBenchmark_libjpg.EncodingOptions.DCTMethod = jpeg.DCTFloat
	jpgBenchmark_libjpg.DecodeOptions.DCTMethod = jpeg.DCTFloat
	//jpgBenchmark_libjpg.Encode(iterations)
	jpgBenchmark_libjpg.Decode(iterations)

	fmt.Printf("libjpg DCTISlow: \n")
	jpgBenchmark_libjpg.EncodingOptions.DCTMethod = jpeg.DCTISlow
	jpgBenchmark_libjpg.DecodeOptions.DCTMethod = jpeg.DCTISlow
	//jpgBenchmark_libjpg.Encode(iterations)
	jpgBenchmark_libjpg.Decode(iterations)

}

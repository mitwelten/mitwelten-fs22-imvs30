package main

import (
	jpgBenchmark_libjpg "benchmarks/jpg/jpgBenchmark_libjpg"
	jpgBenchmark_std "benchmarks/jpg/jpgBenchmark_std"
	"benchmarks/jpg/jpgBenchmark_vips"
	"bytes"
	_ "embed"
	"fmt"
	govips "github.com/davidbyttow/govips/v2/vips"
	"github.com/pixiv/go-libjpeg/jpeg"
	"os"
)

//go:embed image1.jpg
var img1 []byte

//go:embed image2.jpg
var img2 []byte

//go:embed image3.jpg
var img3 []byte

func main() {
	benchmark()
	os.Exit(0)

	iterations := 1000
	jpgBenchmark_libjpg.DecodeOptions.DisableFancyUpsampling = false
	jpgBenchmark_libjpg.DecodeOptions.DisableBlockSmoothing = false
	fmt.Printf("warmup: \n")
	jpgBenchmark_libjpg.DecodeRGBA(100)

	fmt.Printf("---------------------: \n")
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

func quality() {
	jpgBenchmark_libjpg.EncodingOptions.Quality = 100
	jpgBenchmark_libjpg.DecodeEncodeSave()

	jpgBenchmark_libjpg.EncodingOptions.Quality = 80
	jpgBenchmark_libjpg.DecodeEncodeSave()

	jpgBenchmark_libjpg.EncodingOptions.Quality = 60
	jpgBenchmark_libjpg.DecodeEncodeSave()

	jpgBenchmark_libjpg.EncodingOptions.Quality = 40
	jpgBenchmark_libjpg.DecodeEncodeSave()

	jpgBenchmark_libjpg.EncodingOptions.Quality = 20
	jpgBenchmark_libjpg.DecodeEncodeSave()
}
func log(f os.File, library string, function string, width int, height int, time float64) {
	_, _ = f.WriteString(fmt.Sprintf("%v,%v,%v,%v,%.2f\n", library, function, width, height, time))

}

func benchmark() {
	f, err := os.Create("log.csv")
	if err != nil {
		panic(err.Error())
	}
	f.WriteString("library,function,width,height,time\n")

	govips.LoggingSettings(nil, govips.LogLevelError)
	govips.Startup(nil)
	defer govips.Shutdown()
	benchmark1(*f, img1)
	benchmark1(*f, img2)
	benchmark1(*f, img3)
}

func benchmark1(f os.File, b []byte) {
	iterations := 128

	img, _ := jpeg.Decode(bytes.NewReader(b), &jpgBenchmark_libjpg.DecodeOptions)
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()
	var time float64
	time = jpgBenchmark_std.BenchmarkDecode(iterations, b)
	log(f, "go std library", "decode", width, height, time)
	time = jpgBenchmark_libjpg.BenchmarkDecode(iterations, b)
	log(f, "libjpeg-turbo", "decode", width, height, time)
	time = jpgBenchmark_vips.BenchmarkDecode(iterations, b)
	log(f, "libvips", "decode", width, height, time)

	time = jpgBenchmark_std.BenchmarkEncode(iterations, b)
	log(f, "go std library", "encode", width, height, time)
	time = jpgBenchmark_libjpg.BenchmarkEncode(iterations, b)
	log(f, "libjpeg-turbo", "encode", width, height, time)
	time = jpgBenchmark_vips.BenchmarkEncode(iterations, b)
	log(f, "libvips", "encode", width, height, time)
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

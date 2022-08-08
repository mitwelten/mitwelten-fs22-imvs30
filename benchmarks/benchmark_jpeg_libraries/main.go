package main

import (
	jpgBenchmark_libjpg "benchmarks/benchmark_jpeg_libraries/jpgBenchmark_libjpg"
	jpgBenchmark_std "benchmarks/benchmark_jpeg_libraries/jpgBenchmark_std"
	"benchmarks/benchmark_jpeg_libraries/jpgBenchmark_vips"
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
	//create log file and write header
	f, err := os.Create("benchmark_jpeg_libraries.csv")
	if err != nil {
		panic(err.Error())
	}
	_, err = f.WriteString("library,function,width,height,time\n")
	if err != nil {
		panic(err.Error())

	}

	//setup libvips stuff
	govips.LoggingSettings(nil, govips.LogLevelError)
	govips.Startup(nil)
	defer govips.Shutdown()

	iterations := 128
	benchmark(*f, img1, iterations)
	benchmark(*f, img2, iterations)
	benchmark(*f, img3, iterations)
}

func log(f os.File, library string, function string, width int, height int, time float64) {
	_, _ = f.WriteString(fmt.Sprintf("%v,%v,%v,%v,%.2f\n", library, function, width, height, time))
}

func benchmark(f os.File, b []byte, iterations int) {
	img, _ := jpeg.Decode(bytes.NewReader(b), &jpgBenchmark_libjpg.DecodeOptions)
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()

	var time float64
	//decode
	//go std
	time = jpgBenchmark_std.BenchmarkDecode(iterations, b)
	log(f, "go std library", "decode", width, height, time)
	//libjpeg
	time = jpgBenchmark_libjpg.BenchmarkDecode(iterations, b)
	log(f, "libjpeg-turbo", "decode", width, height, time)
	//libvips
	time = jpgBenchmark_vips.BenchmarkDecode(iterations, b)
	log(f, "libvips", "decode", width, height, time)

	//encode
	//go std
	time = jpgBenchmark_std.BenchmarkEncode(iterations, b)
	log(f, "go std library", "encode", width, height, time)
	//libjpeg
	time = jpgBenchmark_libjpg.BenchmarkEncode(iterations, b)
	log(f, "libjpeg-turbo", "encode", width, height, time)
	//libvips
	time = jpgBenchmark_vips.BenchmarkEncode(iterations, b)
	log(f, "libvips", "encode", width, height, time)
}

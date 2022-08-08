package jpgBenchmark_vips

import (
	"fmt"
	govips "github.com/davidbyttow/govips/v2/vips"
	"time"
)

func BenchmarkDecode(iterations int, imageData []byte) float64 {
	fmt.Printf("  Decode vips\n")
	fmt.Printf("    Number of runs: %v\n", iterations)

	for i := 0; i < iterations/10; i++ {
		_, _ = govips.NewImageFromBuffer(imageData)
	}

	start := time.Now()
	for i := 0; i < iterations; i++ {
		img, _ := govips.NewImageFromBuffer(imageData)
		img.Rotate(1)
	}

	end := time.Since(start).Milliseconds()
	fmt.Printf("    Total: %v ms\n", end)
	fmt.Printf("    Per iteration: %v ms\n", float64(end)/float64(iterations))
	return float64(end) / float64(iterations)
}

func BenchmarkEncode(iterations int, imageData []byte) float64 {
	fmt.Printf("  Encode vips\n")
	fmt.Printf("    Number of runs: %v\n", iterations)
	ep := govips.NewJpegExportParams()
	ep.StripMetadata = false
	ep.Quality = 80
	ep.Interlace = false
	ep.OptimizeCoding = false
	ep.SubsampleMode = govips.VipsForeignSubsampleAuto
	ep.TrellisQuant = true
	ep.OvershootDeringing = true
	ep.OptimizeScans = true
	ep.QuantTable = 3

	for i := 0; i < iterations/10; i++ {
		imgout, _ := govips.NewImageFromBuffer(imageData)
		_, _, err := imgout.ExportJpeg(ep)
		if err != nil {
			panic("can't encode benchmark_jpeg_libraries")
		}
	}

	start := time.Now()
	for i := 0; i < iterations; i++ {
		imgout, _ := govips.NewImageFromBuffer(imageData)
		imgout.Rotate(1)
		_, _, err := imgout.ExportJpeg(ep)
		if err != nil {
			panic("can't encode benchmark_jpeg_libraries")
		}
	}

	end := time.Since(start).Milliseconds()
	fmt.Printf("    Total: %v ms\n", end)
	fmt.Printf("    Per iteration: %v ms\n", float64(end)/float64(iterations))
	return float64(end) / float64(iterations)
}

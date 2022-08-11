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
		img, _ := govips.NewImageFromBuffer(imageData)
		_, _ = img.Average() // causes decode
	}

	// libvips works demand driven and won't decode until needed, so we'll force it with the addition of an image operation
	end := time.Duration(0)

	for i := 0; i < iterations; i++ {
		// calculate the time for decode + some operation
		start_ := time.Now()
		img, _ := govips.NewImageFromBuffer(imageData)
		_, _ = img.Average() // causes decode
		end += time.Since(start_)

		// subtract the time needed for the operation
		start_ = time.Now()
		_, _ = img.Average() // no more decode needed here
		end -= time.Since(start_)

	}

	avg := float64(end.Milliseconds()) / float64(iterations)
	fmt.Printf("    Total: %v ms\n", end.Milliseconds())
	fmt.Printf("    Per iteration: %.2f ms\n", avg)
	return avg
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

	end := time.Duration(0)

	for i := 0; i < iterations; i++ {
		// calculate the time for decode + some operation
		imgout, _ := govips.NewImageFromBuffer(imageData)
		start_ := time.Now()
		_, _, err := imgout.ExportJpeg(ep)
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

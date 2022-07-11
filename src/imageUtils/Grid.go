package imageUtils

import (
	"golang.org/x/image/draw"
	"image"
	"log"
	"mjpeg_multiplexer/src/global"
	"mjpeg_multiplexer/src/mjpeg"
)

func Grid(nRows int, nCols int, storages ...*mjpeg.FrameStorage) *mjpeg.MjpegFrame {
	var nCells = nRows * nCols
	var nFrames = len(storages)

	if nFrames > nCells {
		log.Fatalf("Too many frames for this grid configuartion: nRows %v nCols %v, but %v frames to compute\n", nRows, nCols, len(storages))
	}

	if nFrames == 0 {
		log.Fatalf("At least one frame needed\n")
	}

	firstWidthInitial, firstHeightInitial := storages[0].GetImageSize()
	totalWidth := firstWidthInitial * nCols
	totalHeight := firstHeightInitial * nRows

	// check the max totalWidth and totalHeight
	totalWidth, totalHeight = GetFinalImageSize(totalWidth, totalHeight)

	cellWidth := totalWidth / nCols
	cellHeight := totalHeight / nRows

	imageOut := getImageOut(totalWidth, totalHeight)

	var marginStartPoints []image.Point

	for i := 0; i < nCells; i++ {
		var row_ = i / nCols
		var col_ = i % nCols

		if i >= nFrames {
			break
		}

		var sp = image.Point{X: cellWidth * col_, Y: cellHeight * row_}

		//grid lines:
		if row_ == 0 || col_ == 0 {
			marginStartPoints = append(marginStartPoints, sp)
		}

		frame := storages[i].GetFrame()

		// don't redraw already drawn images
		if frame.CachedImage != nil {
			continue
		}

		// Check for resizing
		img := Decode(storages[i])

		if img.Bounds().Dx() != cellWidth || img.Bounds().Dy() != cellHeight {
			img = ResizeOutputFrame(img, cellWidth, cellHeight)
			frame.CachedImage = img
		}
		r := image.Rectangle{Min: sp, Max: sp.Add(img.Bounds().Size())}
		draw.Draw(imageOut, r, img, image.Point{}, draw.Src)
	}

	//place labels
	if global.Config.ShowInputLabel {

		for i := 0; i < nCells; i++ {
			var row_ = i / nCols
			var col_ = i % nCols

			if i >= nFrames {
				break
			}

			var sp = image.Point{X: cellWidth * col_, Y: cellHeight * row_}
			// adjust the start point relative to the border size
			x := sp.X
			y := sp.Y
			if col_ != 0 {
				x += getBorderSize(totalWidth) / 2
			}
			if row_ != 0 {
				y += getBorderSize(totalWidth) / 2
			}

			addLabel(imageOut, x, y, global.Config.InputConfigs[i].Label)

		}
	}

	// draw border
	if global.Config.ShowBorder {
		border := getBorderSize(totalWidth)
		borderVertical := image.Rectangle{Min: image.Point{X: -border / 2}, Max: image.Point{X: border / 2, Y: imageOut.Bounds().Dy()}}
		borderHorizontal := image.Rectangle{Min: image.Point{Y: -border / 2}, Max: image.Point{X: imageOut.Bounds().Dx(), Y: border / 2}}
		for i, p := range marginStartPoints {
			//ignore first point to avoid border lines
			if i == 0 {
				continue
			}

			if p.Y == 0 {
				//vertical
				draw.Draw(imageOut, borderVertical.Add(p), image.Black, image.Point{}, draw.Src)
			} else {
				//horizontal
				draw.Draw(imageOut, borderHorizontal.Add(p), image.Black, image.Point{}, draw.Src)
			}
		}
	}

	return Encode(imageOut)
}

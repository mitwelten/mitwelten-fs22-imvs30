package imageUtils

import (
	"golang.org/x/image/draw"
	"image"
	"mjpeg_multiplexer/src/global"
	"mjpeg_multiplexer/src/mjpeg"
	"mjpeg_multiplexer/src/utils"
)

// PanelLayout describes a layout of a certain panel layout as described in PanelLayouts.go
type PanelLayout struct {
	FirstWidth  float64
	FirstHeight float64

	ChildrenWidth  float64
	ChildrenHeight float64

	ChildrenPositions []utils.FloatPoint

	VerticalBorderPoints   []utils.Tuple[utils.FloatPoint]
	HorizontalBorderPoints []utils.Tuple[utils.FloatPoint]
}

// Panel combines frames by showing all at once, one frame being focused top left
func Panel(layout PanelLayout, startIndex int, storages ...*mjpeg.FrameStorage) *mjpeg.MjpegFrame {
	if imageInContainers == nil {
		imageInContainers = make([]*image.RGBA, len(storages))
	}

	firstWidthInitial, firstHeightInitial := storages[0].GetImageSize()
	totalWidth := int(float64(firstWidthInitial) / layout.FirstWidth)
	totalHeight := int(float64(firstHeightInitial) / layout.FirstHeight)

	totalWidth, totalHeight = GetFinalImageSize(totalWidth, totalHeight)

	//output img
	imageOut := getImageOut(totalWidth, totalHeight)

	//focused image
	sp := image.Point{}
	r := image.Rectangle{Min: sp, Max: image.Point{X: int(float64(totalWidth) * layout.FirstWidth), Y: int(float64(totalHeight) * layout.FirstHeight)}}
	drawFrame(storages[startIndex], r, startIndex)
	if global.Config.ShowInputLabel {
		addLabel(imageOut, sp.X, sp.Y, global.Config.InputConfigs[startIndex].Label)
	}

	//and all other images not in focus
	for i, child := range layout.ChildrenPositions {
		if i+1 >= len(storages) {
			break
		}
		index := (startIndex + i + 1) % len(storages)

		//draw the image
		sp = image.Point{X: int(float64(totalWidth) * child.X), Y: int(float64(totalHeight) * child.Y)}
		delta := image.Point{X: int(float64(totalWidth) * layout.ChildrenWidth), Y: int(float64(totalHeight) * layout.ChildrenHeight)}
		r = image.Rectangle{Min: sp, Max: sp.Add(delta)}
		drawFrame(storages[index], r, index)

		//and add the label
		if global.Config.ShowInputLabel {
			// add offset to avoid overlap with the borders
			offsetW := 0
			offsetH := 0

			//labels to the left / on top don't need an offset
			if global.Config.ShowBorder && child.X != 0 {
				offsetW = getBorderSize(totalWidth) / 2
			}
			if global.Config.ShowBorder && child.Y != 0 {
				offsetH = getBorderSize(totalWidth) / 2
			}
			addLabel(imageOut, sp.X+offsetW, sp.Y+offsetH, global.Config.InputConfigs[index].Label)
		}
	}

	// draw border
	// this is more complicated than in the grid, because the lines can't just be draw over the whole image
	// for each line, the start point of the cell + it's size determines the borders lines length
	if global.Config.ShowBorder {
		border := getBorderSize(totalWidth)

		//vertical lines
		for _, el := range layout.VerticalBorderPoints {
			// calculate width and height
			rectangleWidth := border
			rectangleHeight := int((el.T2.Y*float64(totalHeight))-(el.T1.Y*float64(totalHeight))) + border/2
			// start point is t1 - borderWidth
			sp := image.Point{X: int(el.T1.X*float64(totalWidth)) - border/2, Y: int(el.T1.Y * float64(totalHeight))}
			borderVertical := image.Rectangle{
				Min: sp,
				Max: sp.Add(image.Point{X: rectangleWidth, Y: rectangleHeight}),
			}
			draw.Draw(imageOut, borderVertical, image.Black, image.Point{}, draw.Src)

		}

		//horizontal lines
		for _, el := range layout.HorizontalBorderPoints {
			// calculate width and height
			rectangleWidth := int((el.T2.X*float64(totalWidth))-(el.T1.X*float64(totalWidth))) + border/2
			rectangleHeight := border
			// start point is t1 - borderHeight
			sp := image.Point{X: int(el.T1.X*float64(totalWidth)) - border/2, Y: int(el.T1.Y * float64(totalHeight))}
			borderVertical := image.Rectangle{
				Min: sp,
				Max: sp.Add(image.Point{X: rectangleWidth, Y: rectangleHeight}),
			}
			draw.Draw(imageOut, borderVertical, image.Black, image.Point{}, draw.Src)
		}
	}

	return Encode(imageOut)
}

// drawFrame is a helper method to draw the frame at the given location
func drawFrame(storage *mjpeg.FrameStorage, r image.Rectangle, i int) {
	if imageInContainers[i] != nil {
		DecodeContainer(storage, imageInContainers[i])
	} else {
		image_ := Decode(storage)
		imageInContainers[i] = image_.(*image.RGBA)
	}
	var img image.Image
	img = imageInContainers[i]

	// Check for resizing
	if img.Bounds().Dx() != r.Dx() || img.Bounds().Dy() != r.Dy() {
		img = ResizeOutputFrame(img, r.Dx(), r.Dy())
	}

	draw.Draw(imageOut, r, img, image.Point{}, draw.Src)
}

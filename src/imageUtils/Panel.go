package imageUtils

import (
	"golang.org/x/image/draw"
	"image"
	"mjpeg_multiplexer/src/global"
	"mjpeg_multiplexer/src/mjpeg"
	"mjpeg_multiplexer/src/utils"
)

type PanelLayout struct {
	FirstWidth  float64
	FirstHeight float64

	ChildrenWidth  float64
	ChildrenHeight float64

	ChildrenPositions []utils.FloatPoint

	VerticalBorderPoints   []utils.Tuple[utils.FloatPoint]
	HorizontalBorderPoints []utils.Tuple[utils.FloatPoint]
}

var Slots3 = PanelLayout{
	FirstWidth:     float64(2) / 3.0,
	FirstHeight:    1,
	ChildrenWidth:  float64(1) / 3,
	ChildrenHeight: float64(1) / 2,
	ChildrenPositions: []utils.FloatPoint{
		{X: float64(2) / 3, Y: 0},
		{X: float64(2) / 3, Y: float64(1) / 2},
	},

	VerticalBorderPoints: []utils.Tuple[utils.FloatPoint]{
		{T1: utils.FloatPoint{X: float64(2) / 3, Y: 0}, T2: utils.FloatPoint{X: float64(2) / 3, Y: 1}},
	},
	HorizontalBorderPoints: []utils.Tuple[utils.FloatPoint]{
		{T1: utils.FloatPoint{X: float64(2) / 3, Y: 0.5}, T2: utils.FloatPoint{X: 1, Y: 0.5}},
	},
}

var Slots4 = PanelLayout{
	FirstWidth:     float64(3) / 4.0,
	FirstHeight:    1,
	ChildrenWidth:  float64(1) / 4,
	ChildrenHeight: float64(1) / 3,
	ChildrenPositions: []utils.FloatPoint{
		{X: float64(3) / 4, Y: 0},
		{X: float64(3) / 4, Y: float64(1) / 3},
		{X: float64(3) / 4, Y: float64(2) / 3},
	},
	VerticalBorderPoints: []utils.Tuple[utils.FloatPoint]{
		{T1: utils.FloatPoint{X: float64(3) / 4, Y: 0}, T2: utils.FloatPoint{X: float64(3) / 4, Y: 1}},
	},
	HorizontalBorderPoints: []utils.Tuple[utils.FloatPoint]{
		{T1: utils.FloatPoint{X: float64(3) / 4, Y: float64(1) / 3}, T2: utils.FloatPoint{X: 1, Y: float64(1) / 3}},
		{T1: utils.FloatPoint{X: float64(3) / 4, Y: float64(2) / 3}, T2: utils.FloatPoint{X: 1, Y: float64(2) / 3}},
	},
}

var Slots6 = PanelLayout{
	FirstWidth:     float64(2) / 3.0,
	FirstHeight:    float64(2) / 3.0,
	ChildrenWidth:  float64(1) / 3,
	ChildrenHeight: float64(1) / 3,
	ChildrenPositions: []utils.FloatPoint{
		{X: 0, Y: float64(2) / 3},
		{X: float64(1) / 3, Y: float64(2) / 3},
		{X: float64(2) / 3, Y: float64(2) / 3},
		{X: float64(2) / 3, Y: float64(1) / 3},
		{X: float64(2) / 3, Y: 0},
	},
	VerticalBorderPoints: []utils.Tuple[utils.FloatPoint]{
		{T1: utils.FloatPoint{X: float64(1) / 3, Y: float64(2) / 3}, T2: utils.FloatPoint{X: float64(1) / 3, Y: 1}},
		{T1: utils.FloatPoint{X: float64(2) / 3, Y: 0}, T2: utils.FloatPoint{X: float64(2) / 3, Y: 1}},
	},
	HorizontalBorderPoints: []utils.Tuple[utils.FloatPoint]{
		{T1: utils.FloatPoint{X: float64(2) / 3, Y: float64(1) / 3}, T2: utils.FloatPoint{X: 1, Y: float64(1) / 3}},
		{T1: utils.FloatPoint{X: 0, Y: float64(2) / 3}, T2: utils.FloatPoint{X: 1, Y: float64(2) / 3}},
	},
}

var Slots8 = PanelLayout{
	FirstWidth:     float64(3) / 4.0,
	FirstHeight:    float64(3) / 4.0,
	ChildrenWidth:  float64(1) / 4,
	ChildrenHeight: float64(1) / 4,
	ChildrenPositions: []utils.FloatPoint{
		{X: 0, Y: float64(3) / 4},
		{X: float64(1) / 4, Y: float64(3) / 4},
		{X: float64(2) / 4, Y: float64(3) / 4},
		{X: float64(3) / 4, Y: float64(3) / 4},
		{X: float64(3) / 4, Y: float64(2) / 4},
		{X: float64(3) / 4, Y: float64(1) / 4},
		{X: float64(3) / 4, Y: 0},
	},
	VerticalBorderPoints: []utils.Tuple[utils.FloatPoint]{
		{T1: utils.FloatPoint{X: float64(1) / 4, Y: float64(3) / 4}, T2: utils.FloatPoint{X: float64(1) / 4, Y: 1}},
		{T1: utils.FloatPoint{X: float64(2) / 4, Y: float64(3) / 4}, T2: utils.FloatPoint{X: float64(2) / 4, Y: 1}},
		{T1: utils.FloatPoint{X: float64(3) / 4, Y: 0}, T2: utils.FloatPoint{X: float64(3) / 4, Y: 1}},
	},
	HorizontalBorderPoints: []utils.Tuple[utils.FloatPoint]{
		{T1: utils.FloatPoint{X: float64(3) / 4, Y: float64(1) / 4}, T2: utils.FloatPoint{X: 1, Y: float64(1) / 4}},
		{T1: utils.FloatPoint{X: float64(3) / 4, Y: float64(2) / 4}, T2: utils.FloatPoint{X: 1, Y: float64(2) / 4}},
		{T1: utils.FloatPoint{X: 0, Y: float64(3) / 4}, T2: utils.FloatPoint{X: 1, Y: float64(3) / 4}},
	},
}

func PanelDrawFrame(storage *mjpeg.FrameStorage, r image.Rectangle, flush bool) {
	frame := storage.GetLatestPtr()
	if flush {
		frame.CachedImage = nil
	}

	img := Decode(storage)

	// Check for resizing
	if img.Bounds().Dx() != r.Dx() || img.Bounds().Dy() != r.Dy() {
		img = ResizeOutputFrame(img, r.Dx(), r.Dy())
		//frame.CachedImage = img
	}

	draw.Draw(imageOut, r, img, image.Point{}, draw.Src)
}

var previousIndex = -1

func Panel(layout PanelLayout, startIndex int, storages ...*mjpeg.FrameStorage) *mjpeg.MjpegFrame {
	flush := false
	if previousIndex == -1 || previousIndex != startIndex {
		previousIndex = startIndex
		flush = true
	}

	firstWidthInitial, firstHeightInitial := storages[0].GetImageSize()
	totalWidth := int(float64(firstWidthInitial) / layout.FirstWidth)
	totalHeight := int(float64(firstHeightInitial) / layout.FirstHeight)

	totalWidth, totalHeight = GetFinalImageSize(totalWidth, totalHeight)

	//output img
	imageOut := getImageOut(totalWidth, totalHeight)

	//main character image
	sp := image.Point{}
	r := image.Rectangle{Min: sp, Max: image.Point{X: int(float64(totalWidth) * layout.FirstWidth), Y: int(float64(totalHeight) * layout.FirstHeight)}}

	PanelDrawFrame(storages[startIndex], r, flush)

	if global.Config.ShowInputLabel {
		addLabel(imageOut, sp.X, sp.Y, global.Config.InputConfigs[startIndex].Label)
	}

	for i, child := range layout.ChildrenPositions {
		if i+1 >= len(storages) {
			break
		}
		index := (startIndex + i + 1) % len(storages)
		sp = image.Point{X: int(float64(totalWidth) * child.X), Y: int(float64(totalHeight) * child.Y)}
		delta := image.Point{X: int(float64(totalWidth) * layout.ChildrenWidth), Y: int(float64(totalHeight) * layout.ChildrenHeight)}
		r = image.Rectangle{Min: sp, Max: sp.Add(delta)}

		PanelDrawFrame(storages[index], r, flush)

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

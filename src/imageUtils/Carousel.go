package imageUtils

import (
	"image"
	"mjpeg_multiplexer/src/global"
	"mjpeg_multiplexer/src/mjpeg"
)

func Carousel(input *mjpeg.FrameStorage, index int) *mjpeg.MjpegFrame {
	// with default settings (no resize, no quality change) the image doesn't need to be decoded and encoded at all
	if !global.DecodingNecessary() {
		return input.GetLatestPtr()
	}

	img := Decode(input)
	width, height := input.GetImageSize()

	// resizing
	if (global.Config.Width != width) || (global.Config.Height != height) {
		outputWidth, outputHeight := GetFinalImageSize(width, height)
		img = ResizeOutputFrame(img, outputWidth, outputHeight)
	}

	// label
	if global.Config.ShowInputLabel {
		addLabel(img.(*image.RGBA), 0, 0, global.Config.InputConfigs[index].Label)
	}

	// quality gets adjusted here
	return Encode(img)

}

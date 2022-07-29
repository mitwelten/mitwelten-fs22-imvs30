package imageUtils

import (
	"image"
	"mjpeg_multiplexer/src/global"
	"mjpeg_multiplexer/src/mjpeg"
)

//Carousel combines frames by showing one after the other
func Carousel(input *mjpeg.FrameStorage, index int) *mjpeg.MjpegFrame {
	// with default settings (no resize, no quality change) the image doesn't need to be decoded and encoded at all
	if !global.DecodingNecessary() {
		return input.GetFrame()
	}

	img := Decode(input)

	currentWidth, currentHeight := input.GetImageSize()
	outputWidth, outputHeight := GetFinalImageSize(currentWidth, currentHeight)

	// check for resizing
	if (outputWidth != currentWidth) || (outputHeight != currentHeight) {
		img = ResizeOutputFrame(img, outputWidth, outputHeight)
	}

	// label
	if global.Config.ShowInputLabel {
		addLabel(img.(*image.RGBA), 0, 0, global.Config.InputConfigs[index].Label)
	}

	// quality gets adjusted here
	return Encode(img)

}

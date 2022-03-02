package connection

import "mjpeg_multiplexer/src/mjpeg"

type Source interface {
	ReceiveFrame() (mjpeg.Frame, error)
	GetChannel() chan mjpeg.Frame
}

func RunSource(source Source) {
	go func() {
		for {
			var frame, err = source.ReceiveFrame()
			if err != nil {
				continue
			}
			source.GetChannel() <- frame
			//select {
			// case source  <- frame:
			//default:
			//}
		}
	}()
}

package connection

import "mjpeg_multiplexer/src/mjpeg"

type Source interface {
	ReceiveFrame() (mjpeg.Frame, error)
}

func RunSource(source Source) chan mjpeg.Frame {
	var channel = make(chan mjpeg.Frame)
	go func() {
		for {
			var frame, err = source.ReceiveFrame()
			if err != nil {
        println("Warning, error!")
				continue
			}

			//channel <- frame
			select {
			  case channel  <- frame:
			 default:
		  }
		}
	}()
	return channel
}

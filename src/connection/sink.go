package connection

import "mjpeg_multiplexer/src/mjpeg"

type Sink interface {
	ProcessFrame(frame mjpeg.Frame)
}

func RunSink(sink Sink, sources []Source) {
	//Create an aggregator channel which combines all channels into one
	agg := make(chan mjpeg.Frame)
	for _, s := range sources {
		go func(source Source) {
			for msg := range source.GetChannel() {
				agg <- msg
			}
		}(s)
	}

	go func(agg chan mjpeg.Frame) {
		for {
			frame := <-agg
			sink.ProcessFrame(frame)
		}
	}(agg)
}

package aggregator

import "mjpeg_multiplexer/src/mjpeg"

type AggregatorCombine struct {
}

func (combine AggregatorCombine) Aggregate(channels ...chan mjpeg.Frame) chan mjpeg.Frame {
	agg := make(chan mjpeg.Frame)
	for _, channel := range channels {

		go func(channel_ chan mjpeg.Frame) {
			for msg := range channel_ {
				agg <- msg
			}
		}(channel)

	}

	return agg
}

package aggregator

import "mjpeg_multiplexer/src/mjpeg"

type AggregatorCombine struct {
}

func (combine AggregatorCombine) Aggregate(channels ...chan mjpeg.MjpegFrame) chan mjpeg.MjpegFrame {
	agg := make(chan mjpeg.MjpegFrame)
	for _, channel := range channels {

		go func(channel_ chan mjpeg.MjpegFrame) {
			for msg := range channel_ {
				agg <- msg
			}
		}(channel)

	}

	return agg
}

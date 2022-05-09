package changeDetection

import "mjpeg_multiplexer/src/mjpeg"

// ChangeDetectionScorer compares multiple images and emits difference as int score
type ChangeDetectionScorer interface {
	Score([]mjpeg.MjpegFrame) float64
}

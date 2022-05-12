package global

import "time"

type GlobalConfig struct {
	// pair of url strings and the base64 encoded 'username:password'
	Authentications map[string]string
	// Print the time used for image operations
	LogTime bool
	// Maximal resolution of the resulting image, bigger images will be resized
	MaxWidth  int
	MaxHeight int
	// Minimum amount of time to wait between 2 consecutive reads from the inputs
	MinimumInputDelay time.Duration
	// max framerate for the output
	OutputFramerate float64
}

var Config GlobalConfig

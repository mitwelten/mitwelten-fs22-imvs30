package global

type InputConfig struct {
	Url            string
	Authentication string
	Label          string
}

type GlobalConfig struct {
	// Print the time used for image operations
	LogTime bool
	// Maximal resolution of the resulting image, bigger images will be resized
	Width             int
	Height            int
	IgnoreAspectRatio bool
	// Minimum amount of time to wait between 2 consecutive reads from the inputs
	InputFramerate float64
	// max framerate for the output
	OutputFramerate float64
	// quality for jpeg encoding
	EncodeQuality int
	// use border between images
	Border bool

	UseAuth            bool
	ShowInputLabel     bool
	InputLabelFontSize int
	UseMotion          bool
	InputConfigs       []InputConfig
}

var Config GlobalConfig

func initialConfig() GlobalConfig {
	return GlobalConfig{
		LogTime:            false,
		Width:              -1,
		Height:             -1,
		IgnoreAspectRatio:  false,
		InputFramerate:     -1,
		OutputFramerate:    -1,
		EncodeQuality:      100,
		Border:             false,
		UseAuth:            false,
		ShowInputLabel:     false,
		InputLabelFontSize: 32,
		UseMotion:          false,
		InputConfigs:       []InputConfig{},
	}
}
func SetupInitialConfig() {
	Config = initialConfig()
}

func DecodingNecessary() bool {
	return Config.Height != -1 || Config.Width != -1 || Config.EncodeQuality != 100 || Config.ShowInputLabel
}

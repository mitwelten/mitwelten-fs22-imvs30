package global

type InputConfig struct {
	Url            string
	Authentication string
}

type GlobalConfig struct {
	// Print the time used for image operations
	LogTime bool
	// Maximal resolution of the resulting image, bigger images will be resized
	MaxWidth  int
	MaxHeight int
	// Minimum amount of time to wait between 2 consecutive reads from the inputs
	InputFramerate float64
	// max framerate for the output
	OutputFramerate float64
	// quality for jpeg encoding
	EncodeQuality int
	// use border between images
	Border bool

	UseAuth        bool
	ShowInputLabel bool
	UseMotion      bool
	InputConfigs   []InputConfig
}

var Config GlobalConfig

func initialConfig() GlobalConfig {
	return GlobalConfig{
		LogTime:         false,
		MaxWidth:        -1,
		MaxHeight:       -1,
		InputFramerate:  -1,
		OutputFramerate: -1,
		EncodeQuality:   100,
		Border:          false,
		UseAuth:         false,
		ShowInputLabel:  false,
		UseMotion:       false,
		InputConfigs:    []InputConfig{},
	}
}
func SetupInitialConfig() {
	Config = initialConfig()
}

func DecodingNecessary() bool {
	return Config.MaxHeight != -1 || Config.MaxWidth != -1 || Config.EncodeQuality != 100
}

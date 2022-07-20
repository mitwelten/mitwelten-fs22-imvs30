package global

type InputConfig struct {
	Url            string
	Authentication string
	Label          string
}

type GlobalConfig struct {
	// Print the time used for image operations
	LogFPS bool
	// Maximal resolution of the resulting image, bigger images will be resized
	Width             int
	Height            int
	IgnoreAspectRatio bool
	// max framerate for the output
	OutputFramerate float64
	// quality for jpeg encoding
	EncodeQuality int
	// use border between images
	ShowBorder bool

	UseAuth            bool
	ShowInputLabel     bool
	InputLabelFontSize int
	UseMotion          bool
	InputConfigs       []InputConfig
	//hidden
	DisablePassthrough bool
	AlwaysActive       bool
	Debug              bool
}

var Config GlobalConfig

//initialConfig returns the default values to the global fields
func initialConfig() GlobalConfig {
	return GlobalConfig{
		LogFPS:             false,
		Width:              -1,
		Height:             -1,
		IgnoreAspectRatio:  false,
		OutputFramerate:    -1,
		EncodeQuality:      -1,
		ShowBorder:         false,
		UseAuth:            false,
		ShowInputLabel:     false,
		InputLabelFontSize: 32,
		UseMotion:          false,
		InputConfigs:       []InputConfig{},
		//hidden
		DisablePassthrough: false,
		AlwaysActive:       false,
		Debug:              false,
	}
}

//SetupInitialConfig Assigns the default values to the global config
func SetupInitialConfig() {
	Config = initialConfig()
}

//DecodingNecessary checks whether an jpeg image needs to be decoded based on the configs or if a passthrough is possible
func DecodingNecessary() bool {
	return Config.Height != -1 || Config.Width != -1 || (Config.EncodeQuality != 100 && Config.EncodeQuality != -1) || Config.ShowInputLabel || Config.DisablePassthrough
}

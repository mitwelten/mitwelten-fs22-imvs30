package global

type GlobalConfig struct {
	// pair of url strings and the base64 encoded 'username:password'
	Authentications map[string]string
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
}

var Config GlobalConfig

func initialConfig() GlobalConfig {
	return GlobalConfig{
		Authentications: map[string]string{},
		LogTime:         false,
		MaxWidth:        -1,
		MaxHeight:       -1,
		EncodeQuality:   100,
	}
}
func SetupInitialConfig() {
	Config = initialConfig()
}

func DecodingNecessary() bool {
	return (Config.MaxHeight != -1) || (Config.MaxWidth != -1) || Config.EncodeQuality != 100
}

package global

type GlobalConfig struct {
	Authentications map[string]string
	LogTime         bool
	MaxWidth        int
	MaxHeight       int
}

var Config GlobalConfig

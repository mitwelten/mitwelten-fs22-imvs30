package config

type GlobalConfig struct {
	Authentications map[string]string
}

var Config GlobalConfig

package config

import (
	"mjpeg_multiplexer/src/utils"
)

type GlobalConfig struct {
	Authentications []utils.Authentication
}

var Config GlobalConfig

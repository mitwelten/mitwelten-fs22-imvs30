package utils

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"mjpeg_multiplexer/src/global"
	"os"
)

type JSONAuthentication struct {
	Url      string
	Username string
	Password string
}

const authenticationFileLocation = "authentication.json"

func findInputConfigIndex(url string) int {
	for i, el := range global.Config.InputConfigs {
		if el.Url == url {
			return i
		}
	}
	return -1
}
func ParseAuthenticationFile() {
	bytes, err := os.ReadFile(authenticationFileLocation)
	if err != nil {
		log.Fatalf("Can't open authentication file: %v\n", err.Error())
	}

	var data []JSONAuthentication

	if err := json.Unmarshal(bytes, &data); err != nil {
		log.Fatalf("Can't parse authentication file json: %v\n", err.Error())
	}

	for _, entry := range data {
		auth := entry
		payload := base64.StdEncoding.EncodeToString([]byte(auth.Username + ":" + auth.Password))

		index := findInputConfigIndex(auth.Url)
		if index == -1 {
			continue
		}
		log.Printf("Entry for URL %v found\n", entry.Url)
		if global.Config.Debug {
			log.Printf("   => %v\n", auth.Username+":"+auth.Password)
		}
		global.Config.InputConfigs[index].Authentication = payload
	}
}

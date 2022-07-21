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

//findIndex searches for an entry in the json for input by comparing the URL
func findIndex(config *global.InputConfig, data []JSONAuthentication) int {
	for i, entry := range data {
		auth := entry
		if auth.Url == config.Url {
			return i
		}
	}

	return -1
}

//ParseAuthenticationFile parses the authentication file and writes the results into the global config
//A missing or malformed file will terminate the program
func ParseAuthenticationFile() {
	bytes, err := os.ReadFile(authenticationFileLocation)
	if err != nil {
		log.Fatalf("Can't open authentication file: %v\n", err.Error())
	}

	var data []JSONAuthentication

	if err := json.Unmarshal(bytes, &data); err != nil {
		log.Fatalf("Can't parse authentication file json: %v\n", err.Error())
	}

	// for each input, check if a json entry can be found
	for i, el := range global.Config.InputConfigs {
		jsonIndex := findIndex(&el, data)
		if jsonIndex == -1 {
			//no entry found
			if global.Config.Debug {
				log.Printf("No authentication entry for URL %v found\n", el.Url)
			}
			continue
		}
		//entry found
		auth := data[jsonIndex]
		log.Printf("Authentication entry for URL %v found\n", auth.Url)
		if global.Config.Debug {
			log.Printf("   => %v\n", auth.Username+":"+auth.Password)
		}
		payload := base64.StdEncoding.EncodeToString([]byte(auth.Username + ":" + auth.Password))
		global.Config.InputConfigs[i].Authentication = payload
	}
}

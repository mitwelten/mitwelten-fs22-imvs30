package utils

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"mjpeg_multiplexer/src/global"
	"os"
)

type JSONAuthentication struct {
	Url      string
	Username string
	Password string
}

//findIndex searches for an entry in the json for input by comparing the URL
func findIndex(url string, data []JSONAuthentication) int {
	for i, entry := range data {
		if entry.Url == url {
			return i
		}
	}

	return -1
}

//ParseAuthenticationFile parses the authentication file and for each passed url the authentication string will be created.
//If no entry is found for an url that index will be an empty string
func ParseAuthenticationFile(urls []string, authenticationFileLocation string) ([]string, error) {
	authentications := make([]string, len(urls))

	bytes, err := os.ReadFile(authenticationFileLocation)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Can't open authentication file: %v\n", err.Error()))
	}

	var data []JSONAuthentication

	if err := json.Unmarshal(bytes, &data); err != nil {
		return nil, errors.New(fmt.Sprintf("Can't parse authentication file json: %v\n", err.Error()))
	}

	// for each input, check if a json entry can be found
	for i, url := range urls {
		jsonIndex := findIndex(url, data)
		if jsonIndex == -1 {
			//no entry found
			if global.Config.Debug {
				log.Printf("No authentication entry for URL %v found\n", url)
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
		authentications[i] = payload
	}
	return authentications, nil
}

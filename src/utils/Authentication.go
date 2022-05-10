package utils

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"os"
)

type JSONAuthentication struct {
	Url      string
	Username string
	Password string
}

const authenticationFileLocation = "authentication.json"

func ParseAuthenticationFile() map[string]string {
	bytes, err := os.ReadFile(authenticationFileLocation)
	if err != nil {
		log.Fatalf("Can't open authentication file: %v\n", err.Error())
	}

	var data []JSONAuthentication

	if err := json.Unmarshal(bytes, &data); err != nil {
		log.Fatalf("Can't parse authentication file json: %v\n", err.Error())
	}

	authentications := make(map[string]string)
	keys := make([]string, 1)

	for _, entry := range data {
		auth := entry
		payload := base64.StdEncoding.EncodeToString([]byte(auth.Username + ":" + auth.Password))
		authentications[auth.Url] = payload
		keys = append(keys, auth.Url)
	}

	log.Printf("Found %v authenticaion configs: %v\n", len(authentications), keys)
	return authentications
}

package utils

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"os"
)

type Authentication struct {
	Url     string
	Payload string
}
type JSONAuthentication struct {
	Url      string
	Username string
	Password string
}

const authenticationFileLocation = "authentication.json"

func ParseAuthenticationFile() []Authentication {
	bytes, err := os.ReadFile(authenticationFileLocation)
	if err != nil {
		log.Fatalf("Can't open authentication file: %v\n", err.Error())
	}

	var data []JSONAuthentication

	if err := json.Unmarshal(bytes, &data); err != nil {
		log.Fatalf("Can't parse authentication file json: %v\n", err.Error())
	}

	authentications := make([]Authentication, 0)

	for _, entry := range data {
		auth := entry
		payload := base64.StdEncoding.EncodeToString([]byte(auth.Username + ":" + auth.Password))
		authentications = append(authentications, Authentication{Url: auth.Url, Payload: payload})
	}
	log.Printf("Found %v authenticaion configs.\n", len(authentications))
	return authentications
}

package utils

import (
	"testing"
)

const prefixMissing = "Can't open authentication file"
const prefixInvalidFile = "Can't parse authentication file json"
const fileOK = "test_authentication.json"
const fileMissing = "invalid/path/to/file"
const fileInvalid = "test_authentication_invalid.json"

var urls = []string{"192.168.0.42:8080", "localhost:8081", "hostname:8082"}

func TestParseAuthenticationFileMissing(t *testing.T) {
	_, err := ParseAuthenticationFile(urls, fileMissing)
	ExpectErrorMessage(t, prefixMissing, err)
}

func TestParseAuthenticationFileInvalidFile(t *testing.T) {
	_, err := ParseAuthenticationFile(urls, fileInvalid)
	ExpectErrorMessage(t, prefixInvalidFile, err)
}

func TestParseAuthenticationFileAllUrls(t *testing.T) {
	authentication, err := ParseAuthenticationFile(urls, fileOK)
	ExpectNoError(t, err)

	if len(authentication) != len(urls) {
		t.Fail()
	}

	for i := 0; i < len(authentication); i++ {
		if authentication[i] == "" {
			t.Fail()
		}
	}
}

func TestParseAuthenticationFileSomeUrls(t *testing.T) {
	authentication, err := ParseAuthenticationFile([]string{"test", urls[1], urls[2]}, fileOK)
	ExpectNoError(t, err)

	if len(authentication) != len(urls) {
		t.Fail()
	}

	if authentication[0] != "" {
		t.Fail()
	}

	for i := 1; i < len(authentication); i++ {
		if authentication[i] == "" {
			t.Fail()
		}
	}
}

func TestParseAuthenticationFileDuplicateUrls(t *testing.T) {
	urls[1] = urls[0]
	authentication, err := ParseAuthenticationFile([]string{urls[0], urls[0], urls[2]}, fileOK)
	ExpectNoError(t, err)

	if len(authentication) != len(urls) {
		t.Fail()
	}

	for i := 0; i < len(authentication); i++ {
		if authentication[i] == "" {
			t.Fail()
		}
	}
}
func TestParseAuthenticationFileNoUrls(t *testing.T) {
	// all urls in file
	authentication, err := ParseAuthenticationFile([]string{"1", "2", "3"}, fileOK)
	ExpectNoError(t, err)
	ExpectNoError(t, err)

	for _, el := range authentication {
		if el != "" {
			t.Fail()
		}
	}
}

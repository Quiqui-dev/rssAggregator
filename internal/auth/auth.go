package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	/*Extracts the API key from the headers of a http request
	Example:
		Authorization: ApiKey {insert apikey}
	*/

	val := headers.Get("Authorization")

	if val == "" {
		return "", errors.New("no authentication key provided")
	}

	vals := strings.Split(val, " ")
	if len(vals) != 2 {
		return "", errors.New("malformed auth header")
	}

	if vals[0] != "ApiKey" {
		return "", errors.New("malformed first part of auth header")
	}

	return vals[1], nil
}

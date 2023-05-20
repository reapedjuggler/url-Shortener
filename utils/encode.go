package utils

import (
	"encoding/base64"
	"fmt"
)

func ConvertToBase64(id string) (string, error) {

	// also add a check for checking if the url is valid or not

	if len(id) == 0 {
		return "", fmt.Errorf("Given String is not a valid string")
	}
	return base64.StdEncoding.EncodeToString([]byte(id)), nil
}

package utils

import (
	"encoding/base64"
	"math/big"
)

func ConvertToBase64(id int64) (string, error) {
	
	eb := big.NewInt(id)

	return base64.RawURLEncoding.EncodeToString(eb.Bytes()), nil
}

package utils

import (
	"encoding/base64"
	"math/big"
	"math/rand"
)

const randomChar string = "_*"
const URL_ALREADY_EXIST = "URL already exists"

func ConvertToBase64(id int64) (string, error) {
	eb := big.NewInt(id)
	return base64.RawURLEncoding.EncodeToString(eb.Bytes()), nil
}

func CompleteShortUrl(shortUrl string) string {
	for len(shortUrl) < 7 {
		shortUrl = shortUrl + string(randomChar[rand.Intn(2)])
	}
	return shortUrl
}

package utils

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var alphaNumSymCharset = "_ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789/"

var hexCharset = "abcdef0123456789"

func GenerateRandomStringFromCharset(length int, charset string) string {
	resp := make([]byte, length)
	for i := 0; i < length; i++ {
		resp[i] = charset[rand.Intn(len(charset))]
	}
	return string(resp)
}

func GenerateRandomHex(length int) string {
	return GenerateRandomStringFromCharset(length, hexCharset)
}

func GenerateRandomString(length int) string {
	return GenerateRandomStringFromCharset(length, alphaNumSymCharset)
}

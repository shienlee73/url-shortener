package shortener

import (
	"crypto/sha256"
	"fmt"
	"math/big"

	"github.com/itchyny/base58-go"
)

func sha256Of(input string) []byte {
	hash := sha256.Sum256([]byte(input))
	return hash[:]
}

func base58Encoded(bytes []byte) (string, error) {
	encoding := base58.BitcoinEncoding
	encoded, err := encoding.Encode(bytes)
	if err != nil {
		return "", fmt.Errorf("failed to base58 encode: %v", err)
	}
	return string(encoded), nil
}

func GenerateShortUrl(originalUrl string, userId string) (shortUrl string, err error) {
	if originalUrl == "" || userId == "" {
		return "", fmt.Errorf("originalUrl and userId cannot be empty")
	}

	urlHashBytes := sha256Of(originalUrl + userId)
	generatedNumber := new(big.Int).SetBytes(urlHashBytes).Uint64()
	if shortUrl, err = base58Encoded([]byte(fmt.Sprintf("%d", generatedNumber))) ; err != nil {
		return "", err
	}

	return shortUrl[:6], nil
}
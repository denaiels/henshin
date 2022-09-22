package shortener

import (
	"crypto/sha256"
	"fmt"
	"math/big"

	"github.com/itchyny/base58-go"
)

type ShortenerI interface {
	GenerateShortLink(initialUrl string, userId string) (string, error)
}

type shortener struct {
}

func NewShortener() ShortenerI {
	return &shortener{}
}

func hashSHA256(input string) []byte {
	algorithm := sha256.New()
	algorithm.Write([]byte(input))
	return algorithm.Sum(nil)
}

func base58Encoded(bytes []byte) (string, error) {
	encoding := base58.BitcoinEncoding
	encoded, err := encoding.Encode(bytes)
	if err != nil {
		return "", err
	}
	return string(encoded), nil
}

func (s *shortener) GenerateShortLink(initialUrl string, userId string) (string, error) {
	urlHashBytes := hashSHA256(initialUrl + userId)
	generatedNumber := new(big.Int).SetBytes(urlHashBytes).Uint64()
	finalString, err := base58Encoded([]byte(fmt.Sprintf("%d", generatedNumber)))
	if err != nil {
		return "", err
	}
	return finalString[:8], nil
}

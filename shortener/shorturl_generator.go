package shortener

import (
	"crypto/sha256"
	"fmt"
	"math/big"
	"strings"

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

func sha2560f(input string) []byte {
	algorithm := sha256.New()
	algorithm.Write([]byte(input))
	return algorithm.Sum(nil)
}

func base58Encoded(bytes []byte) string {
	encoding := base58.BitcoinEncoding
	encoded, err := encoding.Encode(bytes)
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to encode!!!")
	}
	return string(encoded)
}

func (s *shortener) GenerateShortLink(initialUrl string, userId string) (string, error) {
	if !strings.HasPrefix(initialUrl, "https://") {
		return "", fmt.Errorf("invalid url! please input a valid url")
	}
	urlHashBytes := sha2560f(initialUrl + userId)
	generatedNumber := new(big.Int).SetBytes(urlHashBytes).Uint64()
	finalString := base58Encoded([]byte(fmt.Sprintf("%d", generatedNumber)))
	return finalString[:8], nil
}

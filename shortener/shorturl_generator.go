package shortener

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/itchyny/base58-go"
	"github.com/rs/zerolog/log"
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

func base58Encoded(bytes []byte) (string, error) {
	encoding := base58.BitcoinEncoding
	encoded, err := encoding.Encode(bytes)
	if err != nil {
		return "", err
	}
	return string(encoded), nil
}

func (s *shortener) GenerateShortLink(initialUrl string, userId string) (string, error) {
	if !strings.HasPrefix(initialUrl, "https://") {
		err := errors.New("invalid url")
		log.Err(err).Msg("Please input a valid url!")
		return "", err
	}
	urlHashBytes := sha2560f(initialUrl + userId)
	generatedNumber := new(big.Int).SetBytes(urlHashBytes).Uint64()
	finalString, err := base58Encoded([]byte(fmt.Sprintf("%d", generatedNumber)))
	if err != nil {
		log.Err(err).Msg("Error while encoding with base58")
		return "", err
	}
	return finalString[:8], nil
}

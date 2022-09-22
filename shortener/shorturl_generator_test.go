package shortener_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"source.golabs.io/daniel.santoso/url-blaster/shortener"
)

const UserId = "e0dba740-fc4b-4977-872c-d360239e6b1a"

func TestShortLinkGeneratorWithYouTubeLink(t *testing.T) {
	s := shortener.NewShortener()

	initialUrl := "https://www.youtube.com/watch?v=dQw4w9WgXcQ"
	shortUrl, err := s.GenerateShortLink(initialUrl, UserId)

	assert.Equal(t, "ASzHLChJ", shortUrl)
	assert.NoError(t, err)
}

func TestShortLinkGeneratorWithGojekLink(t *testing.T) {
	s := shortener.NewShortener()

	initialUrl := "https://www.gojek.com/en-id/"
	shortUrl, err := s.GenerateShortLink(initialUrl, UserId)

	assert.Equal(t, "aSLo122q", shortUrl)
	assert.NoError(t, err)
}

func TestShortLinkGeneratorWithWikiLink(t *testing.T) {
	s := shortener.NewShortener()

	initialUrl := "https://ultra.fandom.com/wiki/Ultraman_(character)"
	shortUrl, err := s.GenerateShortLink(initialUrl, UserId)

	assert.Equal(t, "Y6edurWL", shortUrl)
	assert.NoError(t, err)
}

func TestEncodingFail(t *testing.T) {
	s := shortener.NewShortener()

	generatedNumber := "waokawokaowoakwoakw"
	finalString, err := s.Base58Encoded([]byte(generatedNumber))

	assert.Equal(t, "", finalString)
	assert.Error(t, err)
}

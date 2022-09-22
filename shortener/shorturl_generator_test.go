package shortener_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"source.golabs.io/daniel.santoso/url-blaster/shortener"
)

const UserId = "e0dba740-fc4b-4977-872c-d360239e6b1a"

func TestShortLinkGeneratorWithYouTubeLink(t *testing.T) {
	s := shortener.NewShortener()

	initialUrl1 := "https://www.youtube.com/watch?v=dQw4w9WgXcQ"
	shortUrl1, err := s.GenerateShortLink(initialUrl1, UserId)

	assert.Equal(t, "ASzHLChJ", shortUrl1)
	assert.NoError(t, err)
}

func TestShortLinkGeneratorWithGojekLink(t *testing.T) {
	s := shortener.NewShortener()

	initialUrl2 := "https://www.gojek.com/en-id/"
	shortUrl2, err := s.GenerateShortLink(initialUrl2, UserId)

	assert.Equal(t, "aSLo122q", shortUrl2)
	assert.NoError(t, err)
}

func TestShortLinkGeneratorWithWikiLink(t *testing.T) {
	s := shortener.NewShortener()

	initialUrl3 := "https://ultra.fandom.com/wiki/Ultraman_(character)"
	shortUrl3, err := s.GenerateShortLink(initialUrl3, UserId)

	assert.Equal(t, "Y6edurWL", shortUrl3)
	assert.NoError(t, err)
}

func TestShortLinkGeneratorWithInvalidLink(t *testing.T) {
	s := shortener.NewShortener()

	initialUrl := "hahaha"
	shortUrl, err := s.GenerateShortLink(initialUrl, UserId)

	assert.Equal(t, shortUrl, "")
	assert.NotNil(t, err)
}

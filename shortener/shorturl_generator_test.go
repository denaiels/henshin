package shortener

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const UserId = "e0dba740-fc4b-4977-872c-d360239e6b1a"

func TestShortLinkGenerator(t *testing.T) {
	initialUrl1 := "https://www.youtube.com/watch?v=dQw4w9WgXcQ"
	shortUrl1 := GenerateShortLink(initialUrl1, UserId)

	initialUrl2 := "https://www.gojek.com/en-id/"
	shortUrl2 := GenerateShortLink(initialUrl2, UserId)

	initialUrl3 := "https://ultra.fandom.com/wiki/Ultraman_(character)"
	shortUrl3 := GenerateShortLink(initialUrl3, UserId)

	assert.Equal(t, "ASzHLChJ", shortUrl1)
	assert.Equal(t, "aSLo122q", shortUrl2)
	assert.Equal(t, "Y6edurWL", shortUrl3)
}

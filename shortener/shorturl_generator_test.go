package shortener

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateShortUrl(t *testing.T) {
	// Test case: Empty original URL and user ID
	_, err := GenerateShortUrl("", "")
	assert.Error(t, err)

	// Test case: Non-empty original URL and user ID
	//          : Short URL length is less than 6 characters
	shortURL, err := GenerateShortUrl("https://example.com", "user123")
	assert.NoError(t, err)
	assert.NotEmpty(t, shortURL)
	assert.True(t, len(shortURL) >= 6)
}
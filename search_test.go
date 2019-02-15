package googleit

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBing(t *testing.T) {
	urls, err := Bing("cat animal wiki")
	assert.Nil(t, err)
	assert.True(t, len(urls) > 100)
	assert.Equal(t, "https://en.wikipedia.org/wiki/Cat", urls[0])
}

func TestDuckDuckGo(t *testing.T) {
	urls, err := DuckDuckGo("cat animal wiki")
	assert.Nil(t, err)
	assert.True(t, len(urls) > 100 && len(urls) < 300)
	assert.Equal(t, "https://en.wikipedia.org/wiki/Cat", urls[0])
}

func TestBoth(t *testing.T) {
	urls, err := Search("cats")
	assert.Nil(t, err)
	assert.True(t, len(urls) > 400)
	assert.Equal(t, "https://en.wikipedia.org/wiki/Cat", urls[0])
}

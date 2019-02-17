package googleit

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearchBing(t *testing.T) {
	urls, err := Bing("cat animal wiki")
	assert.Nil(t, err)
	assert.True(t, len(urls) > 100)
	assert.Equal(t, "https://en.wikipedia.org/wiki/Cat", urls[0])
}

func TestSearchDuckDuckGo(t *testing.T) {
	urls, err := DuckDuckGo("cat animal wiki")
	assert.Nil(t, err)
	assert.True(t, len(urls) > 100 && len(urls) < 300)
	assert.Equal(t, "https://en.wikipedia.org/wiki/Cat", urls[0])
}

func TestSearchBoth(t *testing.T) {
	urls, err := Search("cat animal wiki")
	assert.Nil(t, err)
	assert.True(t, len(urls) > 400)
	assert.Equal(t, "https://en.wikipedia.org/wiki/Cat", urls[0])
}

func TestSearchBingWithTor(t *testing.T) {
	urls, err := Bing("cat animal wiki", Options{NumPages: 1, UseTor: true})
	assert.Nil(t, err)
	assert.True(t, len(urls) >= 9)
	assert.Equal(t, "https://en.wikipedia.org/wiki/Cat", urls[0])
	fmt.Println(urls)
}

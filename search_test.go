package googleit

import (
	"testing"

	log "github.com/schollz/logger"
	"github.com/stretchr/testify/assert"
)

func TestSearchStartPage(t *testing.T) {
	log.SetLevel("trace")
	urls, err := StartPage("cat", Options{NumPages: 30})
	assert.Nil(t, err)
	assert.True(t, len(urls) > 100)
}

func TestSearchBing(t *testing.T) {
	log.SetLevel("trace")
	urls, err := Bing("cat animal wiki")
	assert.Nil(t, err)
	assert.True(t, len(urls) > 100)
}

func TestSearchDuckDuckGo(t *testing.T) {
	log.SetLevel("trace")
	urls, err := DuckDuckGo("cat animal wiki")
	assert.Nil(t, err)
	assert.True(t, len(urls) > 100 && len(urls) < 300)
}

func TestSearchBoth(t *testing.T) {
	log.SetLevel("trace")
	urls, err := Search("cat animal wiki")
	assert.Nil(t, err)
	assert.True(t, len(urls) > 10)
}

func TestSearchBingWithTor(t *testing.T) {
	log.SetLevel("trace")
	urls, err := Bing("cat animal wiki", Options{NumPages: 1, UseTor: true})
	assert.Nil(t, err)
	assert.True(t, len(urls) >= 9)
}

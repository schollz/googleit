package googleit

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBing(t *testing.T) {
	urls, err := Bing("cats")
	assert.Nil(t, err)
	assert.True(t, len(urls) > 100)
}

func TestDuckDuckGo(t *testing.T) {
	urls, err := DuckDuckGo("cats")
	assert.Nil(t, err)
	assert.True(t, len(urls) > 100 && len(urls) < 300)
}

func TestBoth(t *testing.T) {
	urls, err := Search("cats")
	assert.Nil(t, err)
	assert.True(t, len(urls) > 400)
}

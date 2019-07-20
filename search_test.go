package googleit

import (
	"fmt"
	"strings"
	"testing"

	log "github.com/schollz/logger"
	"github.com/stretchr/testify/assert"
)

func TestSearchStartPage(t *testing.T) {
	log.SetLevel("trace")
	urls, err := StartPage("banana chocolate chip cookie recipe", Options{NumPages: 30, MustInclude: []string{"banana", "chocolate", "chip", "cookie"}})
	assert.Nil(t, err)
	assert.True(t, len(urls) > 100)
	fmt.Println(strings.Join(urls, "\n"))
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

// func TestGo(t *testing.T) {
// 	log.SetLevel("trace")
// 	urls, err := Search("ingredients chocolate chip cookie recipe", Options{NumPages: 60})
// 	assert.Nil(t, err)
// 	ioutil.WriteFile("urls6.txt", []byte(strings.Join(urls, "\n")), 0644)
// }

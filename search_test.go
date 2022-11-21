package googleit

import (
	"fmt"
	"strings"
	"testing"

	log "github.com/schollz/logger"
	"github.com/stretchr/testify/assert"
)

func TestDuck2(t *testing.T) {
	client, err := GetClient(false)
	assert.Nil(t, err)

	next, err := duckDuckGoSearch(client, DuckSearch{Q: "chocolate chip cookies"})
	assert.Nil(t, err)
	assert.NotEmpty(t, next.Results)
	next, err = duckDuckGoSearch(client, DuckSearch{Q: "chocolate chip cookies"})
	assert.Nil(t, err)
	assert.NotEmpty(t, next.Results)
}

func TestSearchStartPage(t *testing.T) {
	log.SetLevel("trace")
	urls, err := StartPage("banana chocolate chip cookie recipe", Options{NumPages: 30, MustInclude: []string{"banana", "chocolate", "chip", "cookie"}})
	assert.Nil(t, err)
	assert.True(t, len(urls) > 10)
	fmt.Println(strings.Join(urls, "\n"))
}

func TestSearchBing(t *testing.T) {
	log.SetLevel("trace")
	urls, err := Bing("banana chocolate chip cookie recipe", Options{NumPages: 30, MustInclude: []string{"banana", "chocolate", "chip", "cookie"}})
	assert.Nil(t, err)
	assert.True(t, len(urls) > 100)
	fmt.Println(strings.Join(urls, "\n"))
}

func TestSearchDuckDuckGo(t *testing.T) {
	log.SetLevel("trace")
	urls, err := DuckDuckGo("banana chocolate chip cookie recipe", Options{NumPages: 10, MustInclude: []string{"chocolate", "chip", "cookie"}})
	assert.Nil(t, err)
	assert.True(t, len(urls) > 30)
	fmt.Println(strings.Join(urls, "\n"))
}

func TestSearchBoth(t *testing.T) {
	log.SetLevel("trace")
	urls, err := Search("banana chocolate chip cookie recipe", Options{NumPages: 30, MustInclude: []string{"banana", "chocolate", "chip", "cookie"}})
	assert.Nil(t, err)
	assert.True(t, len(urls) > 10)
	fmt.Println(urls)
	fmt.Printf("found %d urls.", len(urls))
}

func TestSearchRanking(t *testing.T) {
	log.SetLevel("trace")
	urls, err := Search("site:tabs.ultimate-guitar.com mr.magic amy winehouse", Options{NumPages: 1, Site: "tabs.ultimate-guitar.com"})
	assert.Nil(t, err)
	assert.True(t, len(urls) > 10)
	assert.Equal(t, "https://tabs.ultimate-guitar.com/tab/amy-winehouse/mr-magic-chords-640059", urls[0])
}

// func TestSearchBingWithTor(t *testing.T) {
// 	log.SetLevel("trace")
// 	urls, err := Bing("cat animal wiki", Options{NumPages: 1, UseTor: true})
// 	assert.Nil(t, err)
// 	assert.True(t, len(urls) >= 9)
// }

// func TestGo(t *testing.T) {
// 	log.SetLevel("trace")
// 	urls, err := Search("banana chocolate chip cookie recipe", Options{Site: "www.allrecipes.com", NumPages: 30, MustInclude: []string{"banana", "chocolate", "chip", "cookie"}})
// 	assert.Nil(t, err)
// 	ioutil.WriteFile("urls.txt", []byte(strings.Join(urls, "\n")), 0644)
// }

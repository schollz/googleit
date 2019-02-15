package googleit

import (
	"bytes"
	"io/ioutil"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHttpClient(t *testing.T) {
	httpClient, err := GetClient(false)
	assert.Nil(t, err)

	resp, err := httpClient.Get("http://icanhazip.com/")
	if err != nil {
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return
	}
	log.Printf("%s", bytes.TrimSpace(body))
	assert.Nil(t, httpClient.Close())
}
func TestHttpClientTor(t *testing.T) {
	httpClient, err := GetClient(true)
	assert.Nil(t, err)

	resp, err := httpClient.Get("http://icanhazip.com/")
	if err != nil {
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return
	}
	log.Printf("%s", bytes.TrimSpace(body))
	assert.Nil(t, httpClient.Close())
}

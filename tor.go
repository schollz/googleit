package googleit

import (
	"bytes"
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/cretz/bine/tor"
)

func useTor() (err error) {
	httpClient := &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 20,
		},
		Timeout: 10 * time.Second,
	}

	torconnection, err := tor.Start(nil, nil)
	if err != nil {
		log.Println(err)
		return
	}

	dialCtx, dialCancel := context.WithTimeout(context.Background(), 3000*time.Hour)
	defer dialCancel()
	// Make connection
	dialer, err := torconnection.Dialer(dialCtx, nil)
	if err != nil {
		log.Println(err)
		return
	}
	httpClient.Transport = &http.Transport{
		DialContext:         dialer.DialContext,
		MaxIdleConnsPerHost: 20,
	}

	// Get /
	resp, err := httpClient.Get("http://icanhazip.com/")
	if err != nil {
		log.Println(err)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("%s", bytes.TrimSpace(body))
	return
}

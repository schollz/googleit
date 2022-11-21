package googleit

import (
	"context"
	"net/http"
	"time"

	"github.com/cretz/bine/tor"
	log "github.com/schollz/logger"
)

// ListToSet convers a list to a set (removing duplicates)
// but preserving order
func ListToSet(s []string) (t []string) {
	m := make(map[string]struct{})
	t = make([]string, len(s))
	i := 0
	for _, v := range s {
		if _, ok := m[v]; ok {
			continue
		}
		m[v] = struct{}{}
		t[i] = v
		i++
	}
	if i == 0 {
		return []string{}
	}
	t = t[:i]
	return
}

type HTTPClient struct {
	Client        *http.Client
	torconnection *tor.Tor
}

func (h *HTTPClient) Get(url string) (resp *http.Response, err error) {
	return h.Client.Get(url)
}

func (h *HTTPClient) Do(req *http.Request) (resp *http.Response, err error) {
	return h.Client.Do(req)
}

func (h *HTTPClient) Close() (err error) {
	if h.torconnection != nil {
		err = h.torconnection.Close()
	}
	return
}

func GetClient(useTor bool) (httpClient *HTTPClient, err error) {
	httpClient = &HTTPClient{Client: &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 20,
		},
		Timeout: 10 * time.Second,
	},
	}
	if !useTor {
		return
	}

	log.Trace("starting tor")
	httpClient.torconnection, err = tor.Start(nil, nil)
	if err != nil {
		log.Error(err)
		return
	}

	dialCtx, dialCancel := context.WithTimeout(context.Background(), 3000*time.Hour)
	defer dialCancel()
	// Make connection
	log.Trace("dialing tor...")
	dialer, err := httpClient.torconnection.Dialer(dialCtx, nil)
	if err != nil {
		log.Error(err)
		return
	}
	log.Trace("making transport...")
	httpClient.Client.Transport = &http.Transport{
		DialContext:         dialer.DialContext,
		MaxIdleConnsPerHost: 20,
	}
	return
}

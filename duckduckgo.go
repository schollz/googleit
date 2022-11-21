package googleit

import (
	"bytes"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/andybalholm/brotli"
	log "github.com/schollz/logger"
)

func DuckDuckGo(query string, ops ...Options) (urls []string, err error) {
	if httpClient == nil {
		if len(ops) > 0 {
			httpClient, err = GetClient(ops[0].UseTor)
		} else {
			httpClient, err = GetClient(false)
		}
		if err != nil {
			return
		}
	}

	pageLimit := 10
	mustInclude := []string{}
	if len(ops) > 0 {
		pageLimit = ops[0].NumPages
		mustInclude = ops[0].MustInclude
		if ops[0].Site != "" {
			query = "site:" + ops[0].Site + " " + query
		}
	}

	currentCount := 1
	urls = []string{}
	nextParameters := DuckSearch{Q: query}
	for i := 0; i < pageLimit; i++ {
		log.Tracef("getting %s", nextParameters)

		nextParameters, err = duckDuckGoSearch(httpClient, nextParameters)
		if err != nil {
			log.Errorf("[duck] %s", err)
			return
		}
		for _, r := range nextParameters.Results {
			doesntHave := ""
			for _, word := range mustInclude {
				if !strings.Contains(r.Title, word) && !strings.Contains(r.URL, word) {
					doesntHave = word
					break
				}
			}
			if doesntHave != "" {
				log.Tracef("[duck] skipping '%s' as it doesn't have '%s'", r.Title, doesntHave)
				continue
			}
			urls = append(urls, r.URL)
			currentCount++
		}
		log.Tracef("[duck] finished page %d/%d", i+1, pageLimit)
	}

	urls = ListToSet(urls)
	return
}

type DuckSearch struct {
	Q          string
	S          string
	NextParams string
	V          string
	O          string
	DC         string
	API        string
	VQD        string
	Results    []Result
}

func duckDuckGoSearch(httpClient *HTTPClient, s DuckSearch) (next DuckSearch, err error) {
	params := url.Values{}
	params.Add("q", s.Q)
	params.Add("s", s.S)
	params.Add("nextParams", s.NextParams)
	params.Add("v", s.V)
	params.Add("o", s.O)
	params.Add("dc", s.DC)
	params.Add("api", s.API)
	params.Add("vqd", s.VQD)
	body := strings.NewReader(params.Encode())

	req, err := http.NewRequest("POST", "https://html.duckduckgo.com/html/", body)
	if err != nil {
		log.Error(err)
		return
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:107.0) Gecko/20100101 Firefox/107.0")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Referer", "https://duckduckgo.com/")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Origin", "https://duckduckgo.com")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Te", "trailers")

	log.Tracef("req: %+v", req)
	log.Tracef("params: %+v", params)
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Error(err)
		return
	}
	defer resp.Body.Close()

	brotliReader := brotli.NewReader(resp.Body)
	b := make([]byte, 100000)
	n, err := brotliReader.Read(b)
	if err != nil {
		log.Error(err)
		return
	}
	html := string(b[:n])

	next.Q = GetStringInBetween(html, `"hidden" name="q" value="`, `"`)
	next.S = GetStringInBetween(html, `"hidden" name="s" value="`, `"`)
	next.NextParams = GetStringInBetween(html, `"hidden" name="nextParams" value="`, `"`)
	next.V = GetStringInBetween(html, `"hidden" name="v" value="`, `"`)
	next.O = GetStringInBetween(html, `"hidden" name="o" value="`, `"`)
	next.DC = GetStringInBetween(html, `"hidden" name="dc" value="`, `"`)
	next.API = GetStringInBetween(html, `"hidden" name="api" value="`, `"`)
	next.VQD = GetStringInBetween(html, `"hidden" name="vqd" value="`, `"`)

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(b[:n]))
	if err != nil {
		log.Error(err)
		return
	}

	doc.Find("h2 > a").Each(func(i int, s *goquery.Selection) {
		href, ok := s.Attr("href")
		if !ok {
			return
		}
		href, _ = url.QueryUnescape(href)
		if strings.Contains(href, "=http") {
			href = "http" + strings.Split(href, "=http")[1]
		}
		if !strings.Contains(href, "http") || strings.Contains(href, "duckduckgo") || strings.Contains(href, "duck.co") {
			return
		}
		log.Tracef("[duck] %s", s.Text())
		log.Tracef("[duck] %s", href)
		next.Results = append(next.Results, Result{Title: s.Text(), URL: href})
	})

	return
}

// GetStringInBetween returns empty string if no start or end string found
func GetStringInBetween(str string, start string, end string) (result string) {
	s := strings.Index(str, start)
	if s == -1 {
		return
	}
	s += len(start)
	e := strings.Index(str[s:], end)
	if e == -1 {
		return
	}
	return str[s : s+e]
}

package googleit

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
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
	nextParameters := fmt.Sprintf(`q=%s&b=&kl=us-en`, url.QueryEscape(query))
	for i := 0; i < pageLimit; i++ {
		log.Tracef("getting %s", nextParameters)
		body := strings.NewReader(nextParameters)
		req, errReq := http.NewRequest("POST", "https://lite.duckduckgo.com/lite/", body)
		if errReq != nil {
			err = errReq
			log.Errorf("[duck] %s", err)
			return
		}

		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:79.0) Gecko/20100101 Firefox/79.0")
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
		req.Header.Set("Accept-Language", "en-US,en;q=0.5")
		req.Header.Set("Referer", "https://lite.duckduckgo.com/")
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Origin", "https://lite.duckduckgo.com")
		req.Header.Set("Connection", "keep-alive")
		req.Header.Set("Upgrade-Insecure-Requests", "1")
		req.Header.Set("Pragma", "no-cache")
		req.Header.Set("Cache-Control", "no-cache")
		req.Header.Set("TE", "Trailers")

		resp, err2 := httpClient.Client.Do(req)
		if err2 != nil {
			err = err2
			log.Errorf("[duck] %s", err)
			return
		}
		if resp.StatusCode != 200 {
			err = fmt.Errorf("[duck] status code error: %d %s", resp.StatusCode, resp.Status)
			log.Errorf("[duck] %s", err)
			return
		}

		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		resp.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		var newResults []Result
		newResults, nextParameters, err2 = captureDuckDuckGo(resp)
		if err2 != nil {
			err = err2
			log.Errorf("[duck] %s", err)
			return
		}
		if len(newResults) == 0 {
			log.Tracef("[duck] no new results: %s", bodyBytes)
			break
		}
		for _, r := range newResults {
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
		log.Tracef("[duck] finished page %d/%d", i, pageLimit)
	}

	urls = ListToSet(urls)
	return
}

func captureDuckDuckGo(res *http.Response) (results []Result, nextParameters string, err error) {
	defer res.Body.Close()
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return
	}

	// Find the review items
	parameters := []string{}
	doc.Find("input[type='hidden']").Each(func(i int, s *goquery.Selection) {
		name, ok := s.Attr("name")
		if !ok {
			return
		}
		value, ok := s.Attr("value")
		if !ok {
			return
		}
		parameters = append(parameters, name+"="+value)
		// fmt.Printf("%d) %s: %s\n", i, name, value)
	})
	nextParameters = strings.Join(parameters, "&")

	// Find the urls
	results = []Result{}
	doc.Find("td > a").Each(func(i int, s *goquery.Selection) {
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
		log.Tracef("[duck] %s", href)
		results = append(results, Result{
			URL:   href,
			Title: strings.ToLower(strings.TrimSpace(s.Text())),
		})
	})
	return
}

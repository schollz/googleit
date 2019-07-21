package googleit

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	log "github.com/schollz/logger"
)

func Bing(query string, ops ...Options) (urls []string, err error) {
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
	}
	if pageLimit < 1 {
		pageLimit = 10
	}

	currentCount := 1
	urls = []string{}
	for i := 0; i < pageLimit; i++ {
		req, err2 := http.NewRequest("GET", fmt.Sprintf("https://www.bing.com/search?q=%s&count=50&first=%d", url.QueryEscape(query), currentCount), nil)
		if err2 != nil {
			err = err2
			return
		}
		req.Header.Set("Dnt", "1")
		req.Header.Set("Accept-Language", "en-US,en;q=0.9")
		req.Header.Set("Upgrade-Insecure-Requests", "1")
		req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3325.181 Safari/537.36")
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
		req.Header.Set("Referer", "https://www.bing.com/")
		req.Header.Set("Authority", "www.bing.com")

		resp, err2 := httpClient.Client.Do(req)
		if err2 != nil {
			err = err2
			return
		}

		var newResults []Result
		newResults, err2 = captureBing(resp)
		if err2 != nil {
			err = err2
			return
		}
		if len(newResults) == 0 {
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
				log.Tracef("[bing] skipping '%s' as it doesn't have '%s'", r.Title, doesntHave)
				continue
			}
			urls = append(urls, r.URL)
			currentCount++
		}
		log.Tracef("[bing] finished page %d", i)
	}

	urls = ListToSet(urls)
	return
}

func captureBing(res *http.Response) (results []Result, err error) {
	defer res.Body.Close()
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return
	}

	// Find the urls
	results = []Result{}
	doc.Find("h2 > a").Each(func(i int, s *goquery.Selection) {
		href, ok := s.Attr("href")
		if !ok {
			return
		}
		href, _ = url.QueryUnescape(href)
		if strings.Contains(href, "=http") {
			href = "http" + strings.Split(href, "=http")[1]
		}
		if !strings.Contains(href, "http") || strings.Contains(href, "bing") || strings.Contains(href, "bing.co") || strings.Contains(href, "clickserve") {
			return
		}
		results = append(results, Result{URL: href})
		log.Tracef("[bing] %s", href)
		results = append(results, Result{
			URL:   href,
			Title: strings.ToLower(strings.TrimSpace(s.Text())),
		})
	})

	return
}

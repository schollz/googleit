package googleit

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func DuckDuckGo(query string, numPages ...int) (urls []string, err error) {
	pageLimit := 20
	if len(numPages) > 0 {
		pageLimit = numPages[0]
	}
	if pageLimit < 1 {
		pageLimit = 20
	}

	urls = []string{}
	nextParameters := fmt.Sprintf(`q=%s&b=&kl=us-en`, url.QueryEscape(query))
	for i := 0; i < pageLimit; i++ {
		body := strings.NewReader(nextParameters)
		req, errReq := http.NewRequest("POST", "https://duckduckgo.com/html/", body)
		if errReq != nil {
			err = errReq
			return
		}
		req.Header.Set("Origin", "https://duckduckgo.com")
		req.Header.Set("Accept-Language", "en-US,en;q=0.9")
		req.Header.Set("Upgrade-Insecure-Requests", "1")
		req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3325.181 Safari/537.36")
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
		req.Header.Set("Cache-Control", "max-age=0")
		req.Header.Set("Authority", "duckduckgo.com")
		req.Header.Set("Referer", "https://duckduckgo.com/")
		req.Header.Set("Dnt", "1")

		resp, err2 := http.DefaultClient.Do(req)
		if err2 != nil {
			err = err2
			return
		}
		if resp.StatusCode != 200 {
			err = fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
			return
		}
		var urls2 []string
		urls2, nextParameters, err2 = captureDuckDuckGo(resp)
		if err2 != nil {
			err = err2
			return
		}
		if len(urls2) == 0 {
			break
		}
		urls = append(urls, urls2...)
	}

	urls = ListToSet(urls)
	return
}

func captureDuckDuckGo(res *http.Response) (urls []string, nextParameters string, err error) {
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
	urls = []string{}
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
		urls = append(urls, href)
		// fmt.Printf("%d) %s\n", i, href)
	})

	return
}

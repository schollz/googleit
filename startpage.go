package googleit

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	log "github.com/schollz/logger"
)

func StartPage(query string, ops ...Options) (urls []string, err error) {
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
	escapedQuery, _ := url.QueryUnescape(query)
	var code string
	for i := 0; i < pageLimit; i++ {
		log.Tracef("[startpage] working on page %d with code %s", i, code)
		var req *http.Request
		if i == 0 {
			body := strings.NewReader(`query=` + escapedQuery + `&cat=web&cmd=process_search&language=english&engine0=v1all&abp=1`)
			req, err = http.NewRequest("POST", "https://www.startpage.com/do/search", body)
			if err != nil {
				return
			}
			req.Header.Set("Referer", "https://www.startpage.com/")
		} else {
			if code == "" {
				break
			}
			body := strings.NewReader(fmt.Sprintf("language=english&abp=1&lui=english&sc=%s&cat=web&query=%s&page=%d",
				code,
				escapedQuery,
				i+1,
			))
			req, err = http.NewRequest("POST", "https://www.startpage.com/sp/search", body)
			if err != nil {
				return
			}
			req.Header.Set("Referer", "https://www.startpage.com/sp/search")
		}
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:68.0) Gecko/20100101 Firefox/68.0")
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
		req.Header.Set("Accept-Language", "en-US,en;q=0.5")
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Connection", "keep-alive")
		req.Header.Set("Upgrade-Insecure-Requests", "1")

		resp, err2 := httpClient.Client.Do(req)
		if err2 != nil {
			err = err2
			return
		}
		var newResults []Result
		newResults, code, err2 = captureStartPage(resp)
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
				if !strings.Contains(r.Title, word) {
					doesntHave = word
					break
				}
			}
			if doesntHave != "" {
				log.Tracef("[startpage] skipping '%s' as it doesn't have '%s'", r.Title, doesntHave)
			}
			urls = append(urls, r.URL)
			currentCount++
		}
	}

	urls = ListToSet(urls)
	return
}

func captureStartPage(res *http.Response) (results []Result, code string, err error) {
	defer res.Body.Close()
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return
	}

	// Find the urls
	results = []Result{}
	doc.Find(".w-gl__result-title").Each(func(i int, s *goquery.Selection) {
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
		results = append(results, Result{strings.ToLower(s.Text()), href})
		log.Tracef("[startpage] '%s' %s", strings.ToLower(s.Text()), href)
	})

	doc.Find("input#sc").Each(func(i int, s *goquery.Selection) {
		id := s.AttrOr("id", "")
		if id == "sc" {
			code, _ = s.Attr("value")
		}
	})

	return
}

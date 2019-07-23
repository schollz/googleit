package googleit

import (
	log "github.com/schollz/logger"
)

func init() {
	log.SetLevel("error")
}

type Options struct {
	Site        string
	NumPages    int
	UseTor      bool
	MustInclude []string
}

type Result struct {
	Title string
	URL   string
}

var httpClient *HTTPClient

func Search(query string, ops ...Options) (urls []string, err error) {
	type Job struct {
		service string
		query   string
	}
	type Result struct {
		err  error
		urls []string
	}
	jobs := make(chan Job, 3)
	results := make(chan Result, 3)

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

	workers := 3
	for w := 1; w <= workers; w++ {
		go func(id int, jobs <-chan Job, results chan<- Result) {
			// generate sha256 filename
			for j := range jobs {
				var r Result
				if j.service == "duck" {
					r.urls, r.err = DuckDuckGo(j.query, ops...)
				} else if j.service == "startpage" {
					r.urls, r.err = StartPage(j.query, ops...)
				} else {
					r.urls, r.err = Bing(j.query, ops...)
				}
				results <- r
			}

		}(w, jobs, results)
	}

	jobs <- Job{
		service: "duck",
		query:   query,
	}
	jobs <- Job{
		service: "startpage",
		query:   query,
	}
	jobs <- Job{
		service: "bing",
		query:   query,
	}
	close(jobs)

	urls = []string{}
	for i := 0; i < 3; i++ {
		r := <-results
		if r.err != nil {
			err = r.err
		}
		urls = append(urls, r.urls...)
	}
	urls = ListToSet(urls)
	return
}

package googleit

type Options struct {
	NumPages int
	UseTor   bool
}

func Search(query string, ops ...Options) (urls []string, err error) {
	type Job struct {
		service string
		query   string
	}
	type Result struct {
		err  error
		urls []string
	}
	jobs := make(chan Job, 2)
	results := make(chan Result, 2)

	workers := 2
	for w := 1; w <= workers; w++ {
		go func(id int, jobs <-chan Job, results chan<- Result) {
			// generate sha256 filename
			for j := range jobs {
				var r Result
				if j.service == "duck" {
					r.urls, r.err = DuckDuckGo(j.query, ops...)
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
		service: "bing",
		query:   query,
	}
	close(jobs)

	urls = []string{}
	for i := 0; i < 2; i++ {
		r := <-results
		if r.err != nil {
			err = r.err
		}
		urls = append(urls, r.urls...)
	}
	urls = ListToSet(urls)
	return
}

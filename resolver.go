package main

import (
	"net/http"
	"sync"
	"time"
)

type Resolver struct {
	workers  int
	wait     sync.WaitGroup
	client   *http.Client
	jobs     chan string
	Resolved chan *Result
}

type Result struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Error string `json:"error"`
}

func NewResolver(workers int) *Resolver {
	r := &Resolver{workers: workers}

	timeout := time.Duration(5 * time.Second)
	r.client = &http.Client{
		Timeout: timeout,
	}

	r.jobs = make(chan string, 100)
	r.Resolved = make(chan *Result, 100)

	return r
}

func (r *Resolver) Start(urls []string) {
	for w := 1; w <= r.workers; w++ {
		r.wait.Add(1)
		go r.worker(r.jobs, r.Resolved)
	}

	for _, url := range urls {
		r.jobs <- url
	}
	close(r.jobs)
}

func (r *Resolver) Stop() {
	close(r.Resolved)
}

func (r *Resolver) Wait() {
	r.wait.Wait()
}

func (r *Resolver) worker(jobs <-chan string, results chan<- *Result) {
	defer r.wait.Done()
	for url := range jobs {
		results <- r.resolve(url)
	}
}

func (r *Resolver) resolve(url string) *Result {
	res := &Result{From: url}

	resp, err := r.client.Head(url)

	if err != nil {
		res.Error = err.Error()
		return res
	}

	res.To = resp.Request.URL.String()

	return res
}

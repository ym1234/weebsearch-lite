package crawl

import (
	"net/http"
	"errors"
	// "fmt"
	"io/ioutil"
	"context"
)

type ProducerFunc func(string) ([]string, []string)
type ConsumerFunc func(string)
type GetURL func(string) ([]byte, error)

type Crawler struct {
	results, next chan string
	Producer ProducerFunc
	Consumer ConsumerFunc
	URL GetURL
	threads int
}

func Default(threads int) *Crawler {
	return &Crawler{
		make(chan string, threads * 2),
		make(chan string, threads * 2),
		func(string)([]string, []string){return nil, nil},
		func(string){},
		DefaultGetter,
		threads}
}

// TODO(ym): clear channnels before starting, run doesn't have to be a one time thing (or does it?)
func (crawler *Crawler) Run(ctx context.Context, start string) {
	for i := 0; i < crawler.threads; i++ {
		go crawler.produce(ctx)
		go crawler.consume(ctx)
	}
	crawler.next <- start
}

func (crawler *Crawler) consume(ctx context.Context) {
	select {
	// case <-ctx.Done():
		// return
	case item := <-crawler.results:
		println()
		crawler.Consumer(item)
	default:
		println("whatever")
	}
}

func (crawler *Crawler) produce(ctx context.Context) {
	select {
	// case <-ctx.Done():
	// 	println("Exiting")
	// 	return
	case item := <-crawler.next:
		// TODO(ym): Handle errors
		body, _ := crawler.URL(item)
		results, _ := crawler.Producer(string(body))
		// results, nexts := crawler.Producer(string(body))
		for _, result := range results {
			crawler.results <- result
		}
		// if len(nexts) == 0 && len(crawler.next) == 0 {
		// 	return
		// }
		// for _, next := range nexts {
		// 	crawler.next <- next
		// }
	default:
		println("Whatever producer")
	}
}

func DefaultGetter(query string) ([]byte, error) {
	response, err := http.Get(query)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		return nil, err
	}

	if string(body) == "" {
		return nil, errors.New("Nothing was received.")
	}
	return body, nil
}


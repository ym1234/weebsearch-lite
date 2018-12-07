package main
import (
	// "time"
	// "bufio"
	// "os"
	"github.com/davecgh/go-spew/spew"
	"fmt"
	"regexp"
	"github.com/ym1234/weebsearch-lite/crawl"
	"context"
	"runtime"
)


func main() {
	testTrie()
	// testCrawler()
}

func testTrie() {
	trie := New()
	trie.Insert("Hello, world!", []string{"test2", "test3", "test10"})
	spew.Dump(trie)
	fmt.Println(trie.GetRecurse("Hello, world!"))
}

// Crawler test, fails miserable lmao
func testCrawler() {
	pattern := regexp.MustCompile(`https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{2,256}\.[a-z]{2,6}\b([-a-zA-Z0-9@:%_\+.~#?&//=]*)`)
	crawler := crawl.Default(runtime.NumCPU())
	out := make(chan string, 10000)
	crawler.Producer = func(body string) ([]string, []string) {
		result := pattern.FindAllString(body, -1)
		out <- fmt.Sprint(result)
		return result, nil
	}
	crawler.Consumer = func(result string) {
			out <- result
	}
	newContext, cancel := context.WithCancel(context.Background())
	crawler.Run(newContext, "https://github.com/coreos/go-systemd")
	for i := range out {
		println(i)
	}
	cancel()
}

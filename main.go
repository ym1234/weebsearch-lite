package main
import (
	// "time"
	// "bufio"
	// "os"
	"github.com/davecgh/go-spew/spew"
	"fmt"
	"regexp"
	"context"
	"runtime"
)


func main() {
	// testTrie()
	testCrawler()
}

// TODO(ym): make actual _test.go
func testTrie() {
	trie := New()
	trie.Insert("Hello, world!", []string{"test2", "test3", "test10"})
	trie.Add("Hello, world!", "test4")
	spew.Dump(trie)
	fmt.Println(trie.GetRecurse("Hello, world!"))
	trie.Clear("Hello, world!")
	fmt.Println(trie.GetRecurse("Hello, world!"))
}

// Crawler test, fails miserable lmao
func testCrawler() {
	pattern := regexp.MustCompile(`https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{2,256}\.[a-z]{2,6}\b([-a-zA-Z0-9@:%_\+.~#?&//=]*)`)
	crawler := Default(runtime.NumCPU())
	out := make(chan string, 10)

	crawler.Producer = func(body string) ([]string, []string) {
		result := pattern.FindAllString(body, -1)
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

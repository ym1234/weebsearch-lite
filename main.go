package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

func main() {
	whatever, _ := authenticate(os.Getenv("USERNAME"), os.Getenv("PASSWORD"))
	Crawl(whatever)
}

func authenticate(username, password string) ([]*http.Cookie, error) {
	client := &http.Client{
		Transport: &http.Transport{TLSClientConfig: &tls.Config{}},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}}

	resp, err := http.Get("http://subs.com.ru/login.php")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// TODO(ym): Error handling
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	x := regexp.MustCompile(`.*<input type="hidden" name="x" value="(.*?)" />`).FindStringSubmatch(string(respBody))[1]
	rcookiettl := regexp.MustCompile(`<input type="hidden" name="rcookiettl" value="(.*?)"/>`).FindStringSubmatch(string(respBody))[1]

	// Not sure if sessid is needed
	endpoint := "http://subs.com.ru/login.php?a=check"
	req, err := http.NewRequest("POST", endpoint, strings.NewReader(url.Values{"rcookiettl": {rcookiettl}, "rpassword": {password}, "rremember": {"on"}, "rusername": {username}, "x": {x}}.Encode()))
	req.Header = map[string][]string{
		"Content-Type": {"application/x-www-form-urlencoded"},
		"Origin":       {"http://subs.com.ru"},
		"Referer":      {"http://subs.com.ru/login.php"}}
	sessid := GetCookie(resp.Cookies(), "PHPSESSID")
	req.AddCookie(sessid)

	newResp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer newResp.Body.Close()

	// Not sure if sessid is needed here either
	return []*http.Cookie{GetCookie(newResp.Cookies(), "ctd25e6ac8e8ab0e48"), sessid}, nil
}

// TODO(ym): make actual _test.go
func testTrie() {
	trie := New()
	trie.Insert("Hello, world!", []string{"test2", "test3", "test10"})
	trie.Add("Hello, world!", "test4")
	spew.Dump(trie)
	trie.Clear("Hello, world!")
	fmt.Println(trie.GetRecurse("Hello, world!"))
}

// // Crawler test, fails miserably lmao
// func testCrawler() {
// 	pattern := regexp.MustCompile(`https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{2,256}\.[a-z]{2,6}\b([-a-zA-Z0-9@:%_\+.~#?&//=]*)`)
// 	crawler := Default(runtime.NumCPU())
// 	out := make(chan string, 10)

// 	crawler.Producer = func(body string) ([]string, []string) {
// 		result := pattern.FindAllString(body, -1)
// 		return result, result
// 	}
// 	crawler.Consumer = func(result string) {
// 		out <- result
// 	}

// 	newContext, cancel := context.WithCancel(context.Background())
// 	crawler.Run(newContext, "http://subs.com.ru/list.php?c=enganime")

// 	for i := range out {
// 		println(i)
// 	}

// 	cancel()
// }

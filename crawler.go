package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type archieve struct {
	ID   string
	Name string
}

var client = http.DefaultClient

const pageEndpoint = "http://subs.com.ru/list.php?c=enganime"
const archieveEndpPoint = "http://subs.com.ru/page.php"

// TODO(ym): context for safely ending

func Crawl(tokenCookie []*http.Cookie) {
	t := time.NewTimer(time.Hour)
	for {
		crawler(tokenCookie)
		t.Reset(time.Hour)
		<-t.C
	}
}

func crawler(tokenCookie []*http.Cookie) error {
	buf, err := ioutil.ReadFile("downloaded")
	downloaded := 0
	if err == nil {
		downloaded, _ = strconv.Atoi(string(buf)[:len(buf)-1])
	}
	println(downloaded)

	buf, err = GetPage(tokenCookie, 2)
	if err != nil {
		return err
	}

	regex := regexp.MustCompile(`<td colspan="8" style="text-align:right;" class="smalltext">Всего архивов в этой секции: (.*?)</td>`)
	results := regex.FindStringSubmatch(string(buf))
	if len(results) != 2 {
		return errors.New("Couldn't find the number of archieves")
	}

	totalArchieves, err := strconv.Atoi(results[1])
	if err != nil {
		return err
	}
	println(totalArchieves)

	numArchieves := totalArchieves-downloaded
	numPages := float64(numArchieves) / 30
	println(numPages)

	return nil
}

func GetPage(tokenCookies []*http.Cookie, num int) ([]byte, error) {
	req, err := http.NewRequest("GET", pageEndpoint+"&d"+strconv.Itoa(num), nil)
	if err != nil {
		return nil, err
	}
	AddCookies(req, tokenCookies)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func GetArchieves(num int, token []*http.Cookie, consumer chan<- archieve) error {
	buf, err := GetPage(token, num)
	if err != nil {
		return err
	}
	GetArchievesFromBuf(buf, consumer)
	return nil
}

func GetArchievesFromBuf(buf []byte, consumer chan<- archieve) {
	// wew this monsterous
	regex := regexp.MustCompile(`(?s)<div><a href="page\.php\?id=(\d*)" title=".*?">.*?</a>.*?<br /></div.*?<div class="smalltext">(.*?)</div>`)
	results := regex.FindAllStringSubmatch(string(buf), -1)

	for _, v := range results {
		// Not sure if this can happen lol
		if len(v) != 3 {
			continue
		}
		// For now, may need to unescape using the html package later
		thisArchieve := archieve{v[1], strings.Replace(v[2], "&nbsp;", "", -1)}
		consumer <- thisArchieve
	}
	close(consumer)
}

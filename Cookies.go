package main

import "net/http"

func GetCookie(cookies []*http.Cookie, name string) *http.Cookie {
	for _, cookie := range cookies {
		if cookie.Name == name {
			return cookie
		}
	}
	return nil
}

func AddCookies(req *http.Request, cookies []*http.Cookie) {
	for _, v := range cookies {
		req.AddCookie(v)
	}
}

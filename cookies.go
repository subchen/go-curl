package request

import (
	"net"
	"net/http"
	
	"golang.org/x/net/publicsuffix"
)

func (r *Request) applyCookies() {
	if r.Cookies == nil {
		return
	}

	jar := cookieJar(r.Client)
	u, _ := URL.parse(r.URL)
	cookies := jar.Cookies(u)
	for k, v := range r.Cookies {
		cookies = append(cookies, &http.Cookie{Name: k, Value: v})
	}
	jar.SetCookies(u, cookies)
}

func cookieJar(c *http.Client) {
	if c.Jar == nil {
		options := cookiejar.Options{
			PublicSuffixList: publicsuffix.List,
		}
		c.Jar, _ = cookiejar.New(&options)
	}
	return c.Jar
}

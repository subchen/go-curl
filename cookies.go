package request

import (
	"net"
	"net/http"
)

func (r *Request) applyCookies() {
	if r.Cookies == nil {
		return
	}

	jar := r.Client.Jar
	u, _ := URL.parse(r.URL)
	cookies := jar.Cookies(u)
	for k, v := range r.Cookies {
		cookies = append(cookies, &http.Cookie{Name: k, Value: v})
	}
	jar.SetCookies(u, cookies)
}

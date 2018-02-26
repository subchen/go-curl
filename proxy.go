package curl

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/net/proxy"
)

func (r *Request) applyProxy(req *http.Request) (err error) {
	if r.ProxyURL == "" {
		return nil
	}

	u, err := url.Parse(r.ProxyURL)
	if err != nil {
		return err
	}

	switch u.Scheme {
	case "http", "https":
		r.Client.Transport.Proxy = http.ProxyURL(u)
	case "socks5":
		dialer, err := proxy.FromURL(u, proxy.Direct)
		if err != nil {
			return err
		}
		r.Client.Transport.Dial = dialer.Dial
		r.Client.Transport.Proxy = http.ProxyFromEnvironment(req)
	}

	return nil
}

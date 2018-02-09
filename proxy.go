package request

import (
	"net"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/net/proxy"
)

func (r *Request) applyProxy() (err error) {
	if r.Proxy == "" {
		return nil
	}

	u, err := url.Parse(r.Proxy)
	if err != nil {
		return err
	}

	switch u.Scheme {
	case "http", "https":
		dialer := &net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}
		r.Client.Transport = &http.Transport{
			Proxy: 	             http.ProxyURL(u),
			Dial:                dialer.Dial,
			TLSHandshakeTimeout: 10 * time.Second,
		}
	case "socks5":
		dialer, err := proxy.FromURL(u, proxy.Direct)
		if err != nil {
			return err
		}
		r.Client.Transport = &http.Transport{
			Proxy:               http.ProxyFromEnvironment,
			Dial:                dialer.Dial,
			TLSHandshakeTimeout: 10 * time.Second,
		}
	}

	return nil
}

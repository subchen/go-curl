package curl

import (
	"net/http"
	"net/url"

	"golang.org/x/net/proxy"
)

func (r *Request) applyProxy() (err error) {
	if r.ProxyURL == "" {
		return nil
	}

	u, err := url.Parse(r.ProxyURL)
	if err != nil {
		return err
	}

	switch u.Scheme {
	case "http", "https":
		r.Client.Transport = &http.Transport{
			Proxy:               http.ProxyURL(u),
			Dial:                r.Transport.Dial,
			TLSHandshakeTimeout: r.Transport.TLSHandshakeTimeout,
			TLSClientConfig:     r.Transport.TLSClientConfig,
		}
	case "socks5":
		dialer, err := proxy.FromURL(u, proxy.Direct)
		if err != nil {
			return err
		}
		r.Client.Transport = &http.Transport{
			Proxy:               http.ProxyFromEnvironment,
			Dial:                dialer.Dial,
			TLSHandshakeTimeout: r.Transport.TLSHandshakeTimeout,
			TLSClientConfig:     r.Transport.TLSClientConfig,
		}
	}

	return nil
}

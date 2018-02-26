package curl

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/url"
	"time"

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
		dialer := &net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}
		r.Client.Transport = &http.Transport{
			Proxy: 	             http.ProxyURL(u),
			Dial:                dialer.Dial,
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: r.InsecureSkipVerify},
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
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: r.InsecureSkipVerify},
			TLSHandshakeTimeout: 10 * time.Second,
		}
	}

	return nil
}

package curl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Request struct {
	Client           *http.Client
	Method           string
	URL              string
	QueryString      interface{}	// *url.Values, map[string]string, map[string][]string
	Headers          map[string]string
	Cookies          map[string]string
	Body             interface{}	// io.Reader, string
	Json             interface{} // any
	Form             interface{} // *url.Values, map[string]string, map[string][]string
	Files            []File
	Auth             interface{} // authorization(BasicAuth, TokenAuth, ...), string
	Proxy            string // http(s)://..., sock5://...
	RedirectDisabled bool
}

var Version = "1.0.0"

func NewRequest() *Request {
	return &Request{
		Client:  new(http.Client),
		Method:  "GET",
		Headers: map[string]string{},
	}
}

func (r *Request) Do() (*Response, error) {
	if r.Client == nil {
		r.Client = new(http.Client)
	}
	if r.Method == "" {
		r.Method = "GET"
	}

	body, contentType, err := r.newBody()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(r.Method, r.newURL(), body)
	if err != nil {
		return nil, err
	}

	r.applyAuth()
	r.applyCookies(req)
	r.applyHeaders(req, contentType)

	resp, err = r.Client.Do(req)
	if err != nil {
		return nil, err
	}
	return &Response{resp}, nil
}

func (r *Request) newURL() string {
	if r.QueryString == nil {
		return r.URL
	}

	qs := newURLValues(r.QueryString)
	if strings.Contains(u, "?") {
		return r.URL + "&" + qs.Encode()
	}
	return r.URL + "?" + qs.Encode()
}

func (r *Request) Reset() {
	r.Method = "GET"
	r.URL = nil
	r.QueryString = nil
	r.Headers = map[string]string{}
	r.Cookies = nil
	r.Body = nil
	r.Json = nil
	r.Form = nil
	r.Files = nil
	r.Auth = nil
	r.Proxy = ""
	r.RedirectDisabled = false
}

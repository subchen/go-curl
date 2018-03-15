package curl

import (
	"net/http"
	"strings"
)

type Request struct {
	Client        *http.Client
	GlobalHeaders map[string]string
	Headers       map[string]string
	Cookies       map[string]string
	Auth          interface{}
}

func NewRequest(client *http.Client) *Request {
	return &Request{
		Client: client,
	}
}

func (r *Request) Call(method string, url string, body interface{}) (*Response, error) {
	payload, err := newPayload(body)
	if err != nil {
		return nil, err
	}

	defer r.reset(payload)

	req, err := http.NewRequest(method, url, payload.reader)
	if err != nil {
		return nil, err
	}

	if r.Client == nil {
		r.Client = new(http.Client)
	}

	applyAuth(r)
	applyHeaders(req, r, payload.contentType, payload.contentLength)
	applyCookies(req, r)

	resp, err := r.Client.Do(req)
	if err != nil {
		return nil, err
	}
	return &Response{resp, nil}, nil
}

func (r *Request) Get(url string) (*Response, error) {
	return r.Call("GET", url, nil)
}

func (r *Request) Post(url string, body interface{}) (*Response, error) {
	return r.Call("POST", url, body)
}

func (r *Request) Put(url string, body interface{}) (*Response, error) {
	return r.Call("PUT", url, body)
}

func (r *Request) Patch(url string, body interface{}) (*Response, error) {
	return r.Call("PATCH", url, body)
}

func (r *Request) Delete(url string) (*Response, error) {
	return r.Call("DELETE", url, nil)
}

func (r *Request) Head(url string) (*Response, error) {
	return r.Call("HEAD", url, nil)
}

func (r *Request) Options(url string) (*Response, error) {
	return r.Call("OPTIONS", url, nil)
}

func (r *Request) WithGlobalHeader(name, value string) *Request {
	if r.GlobalHeaders == nil {
		r.GlobalHeaders = make(map[string]string)
	}
	r.GlobalHeaders[name] = value
	return r
}

func (r *Request) WithHeader(name, value string) *Request {
	if r.Headers == nil {
		r.Headers = make(map[string]string)
	}
	r.Headers[name] = value
	return r
}

func (r *Request) WithCookie(name, value string) *Request {
	if r.Cookies == nil {
		r.Cookies = make(map[string]string)
	}
	r.Cookies[name] = value
	return r
}

func (r *Request) WithBasicAuth(name, passwd string) *Request {
	r.Auth = &BasicAuth{name, passwd}
	return r
}

func (r *Request) WithTokenAuth(token string) *Request {
	r.Auth = &TokenAuth{token}
	return r
}

func (r *Request) reset(payload *Payload) {
	r.Headers = nil
	r.Cookies = nil

	if payload.closer != nil {
		payload.closer.Close()
	}
}

func NewURL(u string, query interface{}) string {
	if query == nil {
		return u
	}

	qs := newValues(query)
	if strings.Contains(u, "?") {
		return u + "&" + qs.Encode()
	}
	return u + "?" + qs.Encode()
}

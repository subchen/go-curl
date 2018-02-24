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
	Client      *http.Client
	Method      string
	URL         string
	QueryString interface{}	// *url.Values, map[string]string, map[string][]string
	Headers     map[string]string
	Cookies     map[string]string
	Body        interface{}	// io.Reader, string
	Json        interface{} // any
	Form        interface{} // *url.Values, map[string]string, map[string][]string
	Files       []File
	Auth        interface{} // authorization(BasicAuth, TokenAuth, ...), string
	Proxy       string // http(s)://..., sock5://...
	Redirect    bool
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

	body, err := r.newBody()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(r.Method, r.newURL(), body)
	if err != nil {
		return nil, err
	}

	r.applyAuth()
	r.applyCookies(req)
	r.applyHeaders(req)

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

func (r *Request) newBody() (io.Reader, error) {
	// html5 payload
	if r.Body != nil {
		switch v := r.Body.(type) {
		case io.Reader:
			r.applyContentType(DefaultPayloadContentType)
			return nil
		case string:
			r.Body = strings.NewReader(v)
			r.applyContentType(DefaultPayloadContentType)
			return nil
		default:
			panic(fmt.Errorf("unsupport request.Body type: %T", v))
		}
	}

	// json
	if r.Json != nil {
		body, err := json.Marshal(r.Json)
		if err != nil {
			return err
		}
		r.Body = bytes.NewReader(b)
		r.applyContentType(DefaultJsonContentType)
		return nil
	}

	// multipart body
	if r.Files != nil {
		return r.newMultipartBody()
	}

	// form data
	if r.Form != nil {
		form := newURLValues(r.Form)
		r.Body = strings.NewReader(form.Encode())
		r.applyContentType(DefaultFormContentType)
		return nil
	}
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
	r.Redirect = false
}

func newURLValues(value interface{}) *url.Values {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case *url.Values:
		return v
	case map[string]string:
		vals := new(url.Values)
		for k, v := range v {
			vals.Set(k, v)
		}
		return vals
	case map[string][]string:
		vals = new(url.Values)
		for k, vs := range v {
			for _, v := range vs {
				vals.Add(k, v)
			}
		}
		return vals
	}
	panic(fmt.Errorf("unable to convert type %T to *url.Values", value))
}

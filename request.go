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
	QueryString interface{}
	Headers     map[string]string
	Cookies     map[string]string
	Body        interface{}
	Json        interface{}
	Form        interface{}
	Files       []File
	Auth        interface{}
	Proxy       string
}

const (
	DEFAULT_FORM_CONTENT_TYPE = "application/x-www-form-urlencoded; charset=utf-8"
	DEFAULT_JSON_CONTENT_TYPE = "application/json; charset=utf-8"
)

var Version = "1.0.0"

func NewRequest() *Request {
	return &Request{
		Client:  new(http.Client),
		Method:  "GET",
		Headers: map[string]string{},
	}
}

func (r *Request) Do() (*Response, error) {
	req, err := r.build()
	if err != nil {
		return nil, err
	}

	resp, err = r.Client.Do(req)

	return nil, &Response{resp, err}
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
}

func (r *Request) newHttpRequest() (*http.Request, error) {
	if r.Client == nil {
		r.Client = new(http.Client)
	}
	if r.Method == "" {
		r.Method = "GET"
	}

	req, err := http.NewRequest(r.Method, r.newURL(), r.newBody())
	if err != nil {
		return nil, err
	}

	r.applyAuth()
	r.applyCookies(req)
	r.applyHeaders(req)

	return req, nil
}

func (r *Request) newURL() {
	if r.QueryString == nil {
		return r.URL
	}

	qs := newURLValues(r.QueryString)
	if strings.Contains(u, "?") {
		return r.URL + "&" + qs.Encode()
	}
	return r.URL + "?" + qs.Encode()
}

func (r *Request) newBody() (body io.Reader, contentType string, err error) {
	// html5 payload
	if r.Body != nil {
		switch v := r.Body.(type) {
		case io.Reader:
			return r.Body, DEFAULT_CONTENT_TYPE, nil
		case string:
			return strings.NewReader(v), DEFAULT_CONTENT_TYPE, nil
		default:
			panic(fmt.Errorf("unsupport request.Body type: %T", v))
		}
	}

	// json
	if r.Json != nil {
		body, err := json.Marshal(r.Json)
		if err != nil {
			return nil, "", err
		}
		return bytes.NewReader(b), DEFAULT_JSON_CONTENT_TYPE, nil
	}

	// multipart body
	if r.Files != nil {
		bodyBuffer := new(bytes.Buffer)
		bodyWriter := multipart.NewWriter(bodyBuffer)
		defer bodyWriter.Close()

		for _, file := range r.Files {
			fileWriter, err := bodyWriter.CreateFormFile(file.Fieldname, file.Filename)
			if err != nil {
				return nil, "", err
			}

			f, err := os.Open(file.Filename)
			if err != nil {
				return nil, "", err
			}
			defer f.Close()

			_, err = io.Copy(fileWriter, f)
			if err != nil {
				return nil, "", err
			}
		}
		if r.Form != nil {
			form := newURLValues(r.Form)
			for k, v := range form {
				bodyWriter.WriteField(k, v)
			}
		}

		return bodyBuffer, bodyWriter.FormDataContentType(), nil
	}

	if r.Form != nil {
		form := newURLValues(r.Form)
		return strings.NewReader(form.Encode()), DEFAULT_FORM_CONTENT_TYPE, nil
	}
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

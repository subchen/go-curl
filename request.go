package curl

import (
	"crypto/tls"
	"net"
	"net/http"
	"strings"
	"time"
)

var (
	DefaultRequestTimeout      = 30 * time.Second
	DefaultDailTimeout         = 30 * time.Second
	DefaultTLSHandshakeTimeout = 30 * time.Second
	DefaultInsecureSkipVerify  = false
)

type Request struct {
	Client    *http.Client
	Transport *http.Transport

	Method           string
	URL              string
	QueryString      interface{} // url.Values, map[string]string, map[string][]string
	Headers          map[string]string
	Cookies          map[string]string
	Body             interface{} // io.Reader, string
	JSON             interface{} // any
	Form             interface{} // url.Values, map[string]string, map[string][]string
	Files            []UploadFile
	Auth             interface{} // authorization(BasicAuth, TokenAuth, ...), string
	ProxyURL         string      // http(s)://..., sock5://...
	RedirectDisabled bool
}

func NewRequest() *Request {
	return &Request{
		Method: "GET",
	}
}

func (r *Request) Do() (*Response, error) {
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

	r.newClient()
	r.applyProxy()

	r.applyAuth()
	r.applyHeaders(req, contentType)
	r.applyCookies(req)

	resp, err := r.Client.Do(req)
	if err != nil {
		return nil, err
	}
	return &Response{resp, nil}, nil
}

func (r *Request) newClient() {
	if r.Transport == nil {
		r.Transport = &http.Transport{
			Dial: (&net.Dialer{
				Timeout: DefaultDailTimeout,
			}).Dial,
			TLSHandshakeTimeout: DefaultTLSHandshakeTimeout,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: DefaultInsecureSkipVerify,
			},
		}
	}

	if r.Client == nil {
		r.Client = &http.Client{
			Timeout:   DefaultRequestTimeout,
			Transport: r.Transport,
		}
	}
}

func (r *Request) newURL() string {
	if r.QueryString == nil {
		return r.URL
	}

	qs := newURLValues(r.QueryString)
	if strings.Contains(r.URL, "?") {
		return r.URL + "&" + qs.Encode()
	}
	return r.URL + "?" + qs.Encode()
}

func (r *Request) SetQueryString(qs interface{}) *Request {
	r.QueryString = qs
	return r
}

func (r *Request) SetBody(body interface{}) *Request {
	r.Body = body
	return r
}

func (r *Request) SetForm(form interface{}) *Request {
	r.Form = form
	return r
}

func (r *Request) SetJSON(json interface{}) *Request {
	r.JSON = json
	return r
}

func (r *Request) AddFile(field, filename string) *Request {
	r.Files = append(r.Files, UploadFile{
		Fieldname: field,
		Filename:  filename,
	})
	return r
}

func (r *Request) AddHeader(name, value string) *Request {
	if r.Headers == nil {
		r.Headers = make(map[string]string)
	}
	r.Headers[name] = value
	return r
}

func (r *Request) AddCookie(name, value string) *Request {
	if r.Cookies == nil {
		r.Cookies = make(map[string]string)
	}
	r.Cookies[name] = value
	return r
}

func (r *Request) SetBasicAuth(name, passwd string) *Request {
	r.Auth = &BasicAuth{name, passwd}
	return r
}

func (r *Request) SetTokenAuth(token string) *Request {
	r.Auth = &TokenAuth{token}
	return r
}

func (r *Request) SetProxyURL(url string) *Request {
	r.ProxyURL = url
	return r
}

func (r *Request) Get(url string) (*Response, error) {
	r.Method = "GET"
	r.URL = url
	return r.Do()
}

func (r *Request) Post(url string) (*Response, error) {
	r.Method = "POST"
	r.URL = url
	return r.Do()
}

func (r *Request) Put(url string) (*Response, error) {
	r.Method = "PUT"
	r.URL = url
	return r.Do()
}

func (r *Request) Patch(url string) (*Response, error) {
	r.Method = "PATCH"
	r.URL = url
	return r.Do()
}

func (r *Request) Delete(url string) (*Response, error) {
	r.Method = "DELETE"
	r.URL = url
	return r.Do()
}

func (r *Request) Head(url string) (*Response, error) {
	r.Method = "HEAD"
	r.URL = url
	return r.Do()
}

func (r *Request) Options(url string) (*Response, error) {
	r.Method = "OPTIONS"
	r.URL = url
	return r.Do()
}

func (r *Request) Reset() *Request {
	r.Method = "GET"
	r.URL = ""
	r.QueryString = nil
	r.Headers = nil
	r.Cookies = nil
	r.Body = nil
	r.JSON = nil
	r.Form = nil
	r.Files = nil
	r.Auth = nil
	return r
}

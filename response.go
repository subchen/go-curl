package curl

import (
	"compress/gzip"
	"compress/zlib"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Response ...
type Response struct {
	*http.Response
	body []byte
}

// JSON return Response Body as data.Query
func (resp *Response) JSON() (*data.Query, error) {
	b, err := resp.Body()
	if err != nil {
		return nil, err
	}
	return json.NewQueryFromBytes(b)
}

// Content return Response Body as []byte
func (resp *Response) Body() ([]byte, error) {
	if resp.body != nil {
		return resp.body, nil
	}

	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		if reader, err = gzip.NewReader(resp.Body); err != nil {
			return nil, err
		}
	case "deflate":
		if reader, err = zlib.NewReader(resp.Body); err != nil {
			return nil, err
		}
	default:
		reader = resp.Body
	}

	defer reader.Close()
	if b, err = ioutil.ReadAll(reader); err != nil {
		return nil, err
	}

	resp.body = b
	return b, err
}

// Text return Response Body as string
func (resp *Response) Text() (string, error) {
	b, err := resp.Body()
	if err != nil {
		return "", nil
	}
	return string(b), nil
}

// OK check Response StatusCode < 400 ?
func (resp *Response) OK() bool {
	return resp.StatusCode < 400
}

// URL return finally request url
func (resp *Response) URL() (*url.URL, error) {
	u := resp.Request.URL
	switch resp.StatusCode {
	case http.StatusMovedPermanently, http.StatusFound,
		http.StatusSeeOther, http.StatusTemporaryRedirect:
		location, err := resp.Location()
		if err != nil {
			return nil, err
		}
		u = u.ResolveReference(location)
	}
	return u, nil
}

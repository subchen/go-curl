package curl

import (
	"compress/gzip"
	"compress/zlib"
	"encoding/json"
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

// JSON return Response Body as map[string]interface{}
func (resp *Response) JSON() (map[string]interface{}, error) {
	v := make(map[string]interface{})
	err := resp.JSONObject(&v)
	return v, err
}

// UnmarshalJSON unmarshal Response Body
func (resp *Response) UnmarshalJSON(data interface{}) error {
	b, err := resp.Body()
	if err != nil {
		return nil, err
	}
	return json.Unmarshal(b, data)
}

// RequestURL return finally request url
func (resp *Response) RequestURL() (*url.URL, error) {
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

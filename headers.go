package curl

import (
	"net/http"
)

var DefaultUserAgent = "subchen/go-curl"

var DefaultHeaders = map[string]string{
	"Connection":      "keep-alive",
	"Accept-Encoding": "gzip, deflate",
	"Accept":          "*/*",
	"User-Agent":      DefaultUserAgent,
}


const (	
	DefaultPayloadContentType = "application/octoc-streams"	
	DefaultJsonContentType    = "application/json; charset=utf-8"
	DefaultFormContentType    = "application/x-www-form-urlencoded; charset=utf-8"
)

func (r *Request) applyHeaders(req *http.Request, contentType string) {
	// apply custom Headers
	for k, v := range r.Headers {
		req.Header.Set(k, v)
	}

	// apply contentType
	if _, ok := r.Headers["Content-Type"]; !ok {
		req.Header.Set("Content-Type", contentType)
	}

	// apply default headers
	for k, v := range DefaultHeaders {
		if _, ok := r.Headers[k]; !ok {
			req.Header.Set(k, v)
		}
	}
}

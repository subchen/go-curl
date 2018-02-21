package requests

import (
	"net/http"
)

var DefaultUserAgent = "subchen/go-cli:v0.1.0"

var DefaultHeaders = map[string]string{
	"Connection":      "keep-alive",
	"Accept-Encoding": "gzip, deflate",
	"Accept":          "*/*",
	"User-Agent":      DefaultUserAgent,
}


func (r *Request) applyHeaders(req *http.Request) {
	// apply custom Headers
	for k, v := range r.Headers {
		req.Header.Set(k, v)
	}

	// apply default headers
	for k, v := range DefaultHeaders {
		if _, ok := r.Headers[k]; !ok {
			req.Header.Set(k, v)
		}
	}
}

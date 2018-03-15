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

func applyHeaders(req *http.Request, r *Request, contentType string) {
	// apply contentType
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	// apply custom Headers
	for k, v := range r.Headers {
		req.Header.Set(k, v)
	}

	// apply custom global Headers
	for k, v := range r.GlobalHeaders {
		if _, ok := req.Header[k]; !ok {
			req.Header.Set(k, v)
		}
	}

	// apply default headers
	for k, v := range DefaultHeaders {
		if _, ok := req.Header[k]; !ok {
			req.Header.Set(k, v)
		}
	}
}

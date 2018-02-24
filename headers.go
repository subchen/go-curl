package curl

import (
	"net/http"
)

var DefaultUserAgent = "subchen/go-curl:" + Version

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

func (r *Request) setContentType(contentType string) {
	if _, ok := r.Headers["Content-type"]; !ok {
		r.Headers["Content-type"] = contentType
	}	
}

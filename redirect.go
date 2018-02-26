package curl

import (
	"errors"
	"net/http"
)

// DefaultRedirectLimit define max redirect counts
var DefaultRedirectLimit = 10

// ErrMaxRedirect when redirect times great than DefaultRedirectLimit will return this error
var ErrMaxRedirect = errors.New("Exceeded max redirects")

func defaultCheckRedirect(req *http.Request, via []*http.Request) error {
	if len(via) > DefaultRedirectLimit {
		return ErrMaxRedirect
	}
	if len(via) == 0 {
		return nil
	}
	// Redirect requests with the first Header
	for key, val := range via[0].Header {
		// Don't copy Referer Header
		if key != "Referer" {
			req.Header[key] = val
		}
	}
	return nil
}

func disableCheckRedirect(req *http.Request, via []*http.Request) error {
	return http.ErrUseLastResponse
}

func (r *Request) applyRedirect() {
	if r.Client.CheckRedirect == nil {
		if r.RedirectDisabled {
			r.Client.CheckRedirect = disableCheckRedirect
		} else {
			r.Client.CheckRedirect = defaultCheckRedirect
		}
	}
}

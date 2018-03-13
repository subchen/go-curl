package curl

import (
	"encoding/base64"
	"fmt"
)

type authenticator interface {
	HeaderValue() string
}

type BasicAuth struct {
	Username string
	Password string
}

type TokenAuth struct {
	Token string
}

func (a *BasicAuth) HeaderValue() string {
	auth := a.Username + ":" + a.Password
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
}

func (a *TokenAuth) HeaderValue() string {
	return "token " + a.Token
}

func applyAuth(r *Request) {
	if r.Auth == nil {
		return
	}

	if r.Headers == nil {
		r.Headers = make(map[string]string)
	}

	switch v := r.Auth.(type) {
	case authenticator:
		r.Headers["Authorization"] = v.HeaderValue()
	case string:
		r.Headers["Authorization"] = v
	default:
		panic(fmt.Errorf("unsupported request.Auth type: %T", v))
	}
}

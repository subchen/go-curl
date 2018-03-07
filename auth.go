package curl

import (
	"encoding/base64"
	"fmt"
)

type authenticator interface {
	RequestHeaderValue() string
}

type BasicAuth struct {
	Username string
	Password string
}

type TokenAuth struct {
	Token string
}

func (a *BasicAuth) RequestHeaderValue() string {
	auth := a.Username + ":" + a.Password
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
}

func (a *TokenAuth) RequestHeaderValue() string {
	return "token " + a.Token
}

func (r *Request) applyAuth() {
	if r.Auth == nil {
		return
	}

	switch v := r.Auth.(type) {
	case authenticator:
		r.Headers["Authorization"] = v.RequestHeaderValue()
	case string:
		r.Headers["Authorization"] = v
	default:
		panic(fmt.Errorf("unsupported request.Auth type: %T", v))
	}
}

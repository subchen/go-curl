package request

import (
	"encoding/base64"
)

type authorize interface {
	Authorization() string
}

type BasicAuth struct {
	Username string
	Password string
}

type TokenAuth struct {
	Token string
}

func (a *BasicAuth) Authorization() string {
	auth := a.Username + ":" + b.Password
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
}

func (a *TokenAuth) Authorization() string {
	return "token " + a.Token
}

func (r *Request) applyAuth() {
	if r.Auth == nil {
		return
	}

	switch v := r.Auth.(type) {
	case authorize:
		r.Headers["Authorization"] = v.Authorization()
	case string:
		r.Headers["Authorization"] = v
	default:
		panic(fmt.Errorf("unsupport request.Auth type: %T", v))
	}
}

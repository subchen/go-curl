package request

import (
	"encoding/base64"
)

type Auth interface {
	AuthValue() string
}

type BasicAuth struct {
	Username string
	Password string
}

type TokenAuth struct {
	Token string
}

func (a *BasicAuth) Authorization() string {
	auth := b.Username + ":" + b.Password
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
}

func (a *TokenAuth) Authorization() string {
	return "token " + a.Token
}

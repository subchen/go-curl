# go-curl

[![Go Report Card](https://goreportcard.com/badge/github.com/subchen/go-curl)](https://goreportcard.com/report/github.com/subchen/go-curl)
[![GoDoc](https://godoc.org/github.com/subchen/go-curl?status.svg)](https://godoc.org/github.com/subchen/go-curl)

A Go HTTP client library for creating and sending API requests

## Examples

```go
import "github.com/subchen/go-curl"
```

### Basic request

```go
req := curl.NewRequest(nil)

// GET
resp, err := req.Get("http://example.com/api/users")
if err != nil {
	log.Fatalln("Unable to make request: ", err)
}
fmt.Println(resp.Text())

// POST
user := &User{...}
resp, err := req.Post("http://example.com/api/users", user)
if err != nil {
	log.Fatalln("Unable to make request: ", err)
}
fmt.Println(resp.Text())
```

### Chained request

```go
user := newUser()
req := curl.NewRequest()
resp, err := req.WithBasicAuth("admin", "passwd").WithHeader("x-trace-id", "123").Post("http://example.com/api/users")
if err != nil {
	log.Fatalln("Unable to make request: ", err)
}
fmt.Println(resp.Text())
```

# go-curl

[![Go Report Card](https://goreportcard.com/badge/github.com/subchen/go-curl)](https://goreportcard.com/report/github.com/subchen/go-curl)
[![GoDoc](https://godoc.org/github.com/subchen/go-curl?status.svg)](https://godoc.org/github.com/subchen/go-curl)

A Go HTTP client library for creating and sending API requests

## Examples

```
import "github.com/subchen/go-curl"
```

### Basic request

```go
req := curl.NewRequest()
req.Method = "GET"
req.URL = "http://example.com/api/users"
resp, err := req.Do()
if err != nil {
	log.Fatalln("Unable to make request: ", err)
}
fmt.Println(resp.Text())
```

### Chained request

```go
user := newUser()
req := curl.NewRequest()
resp, err := req.SetJSON(user).Post("http://example.com/api/users")
if err != nil {
	log.Fatalln("Unable to make request: ", err)
}
fmt.Println(resp.Text())
```

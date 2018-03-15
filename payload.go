package curl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

type Payload struct {
	reader      io.Reader
	closer      io.Closer
	contentType string
}

type UploadFile struct {
	Fieldname string
	Filename  string
}

var emptyPayload = new(Payload)

func newPayload(body interface{}) *Payload {
	if body == nil {
		return emptyPayload
	}

	switch v := body.(type) {
	case *Payload:
		return v
	case Payload:
		return &v
	case string:
		return NewStringPayload(v)
	case []byte:
		return NewBytesPayload(v)
	case map[string]string:
		return NewFormPayload(v)
	case map[string][]string:
		return NewFormPayload(v)
	case url.Values:
		return NewFormPayload(v)
	}

	// io.reader
	if v, ok := body.(io.Reader); ok {
		return NewReaderPayload(v)
	}

	// struct
	t := reflect.TypeOf(body)
	if t.Kind() == reflect.Struct {
		return NewJSONPayload(&body)
	}
	// point to struct
	if t.Kind() == reflect.Ptr && reflect.ValueOf(body).Elem().Kind() == reflect.Struct {
		return NewJSONPayload(body)
	}

	panic(fmt.Errorf("unsupported payload type: %T", body))
}

func NewStringPayload(body string) *Payload {
	return &Payload{
		reader: strings.NewReader(body),
	}
}

func NewBytesPayload(body []byte) *Payload {
	return &Payload{
		reader: bytes.NewReader(body),
	}
}

func NewReaderPayload(reader io.Reader) *Payload {
	return &Payload{
		reader: reader,
	}
}

func NewFilePayload(filename string) *Payload {
	f, err := os.Open(filename)
	if err != nil {
		return nil
	}

	ext := filepath.Ext(filename)

	return &Payload{
		reader:      f,
		closer:      f,
		contentType: mime.TypeByExtension(ext),
	}
}

func NewJSONPayload(obj interface{}) *Payload {
	body, err := json.Marshal(obj)
	if err != nil {
		return nil
	}
	return &Payload{
		reader:      bytes.NewReader(body),
		contentType: "application/json; charset=utf-8",
	}
}

func NewFormPayload(form interface{}) *Payload {
	body := newValues(form)
	return &Payload{
		reader:      strings.NewReader(body.Encode()),
		contentType: "application/x-www-form-urlencoded; charset=utf-8",
	}
}

func NewMultipartPayload(files []UploadFile, form interface{}) (*Payload, error) {
	bodyBuffer := new(bytes.Buffer)
	bodyWriter := multipart.NewWriter(bodyBuffer)
	defer bodyWriter.Close()

	for _, file := range files {
		fileWriter, err := bodyWriter.CreateFormFile(file.Fieldname, file.Filename)
		if err != nil {
			return nil, err
		}

		f, err := os.Open(file.Filename)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		_, err = io.Copy(fileWriter, f)
		if err != nil {
			return nil, err
		}
	}

	if form != nil {
		for k, vs := range newValues(form) {
			for _, v := range vs {
				bodyWriter.WriteField(k, v)
			}
		}
	}

	return &Payload{
		reader:      bodyBuffer,
		contentType: bodyWriter.FormDataContentType(),
	}, nil
}

func newValues(value interface{}) url.Values {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case url.Values:
		return v
	case map[string]string:
		vals := url.Values{}
		for k, v := range v {
			vals.Set(k, v)
		}
		return vals
	case map[string][]string:
		vals := url.Values{}
		for k, vs := range v {
			for _, v := range vs {
				vals.Add(k, v)
			}
		}
		return vals
	}
	panic(fmt.Errorf("unable to convert type %T to url.Values", value))
}

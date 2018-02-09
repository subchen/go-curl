package request

import (
	"bytes"
	"io"
	"mime/multipart"
	"os"
)

type File struct {
	Fieldname string
	Filename  string
}

func (r *Request) newMultipartBody() error {
	bodyBuffer := new(bytes.Buffer)
	bodyWriter := multipart.NewWriter(bodyBuffer)
	defer bodyWriter.Close()

	for _, file := range r.Files {
		fileWriter, err := bodyWriter.CreateFormFile(file.Fieldname, file.Filename)
		if err != nil {
			return err
		}

		f, err := os.Open(file.Filename)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = io.Copy(fileWriter, f)
		if err != nil {
			return err
		}
	}

	if r.Form != nil {
		form := newURLValues(r.Form)
		for k, v := range form {
			bodyWriter.WriteField(k, v)
		}
	}

	r.Headers["Content-Type"] = bodyWriter.FormDataContentType()
	r.Body = bodyBuffer

	return nil
}

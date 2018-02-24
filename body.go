package curl

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

func (r *Request) applyBody() (body io.Reader, err error) {
		
}

func (r *Request) newJSONBody() (body io.Reader, string, error) {
	body, err := json.Marshal(r.Json)
	if err != nil {
		return nli, "", err
	}
	return bytes.NewReader(b), DefaultJsonContentType, nil
}

func (r *Request) newFormBody() (body io.Reader, string, error) {
	form := newURLValues(r.Form)
	if r.Files != nil {
		return newMultipartBody(r.Files, form)
	}
	return strings.NewReader(form.Encode()), DefaultFormContentType, nil
}

func newMultipartBody(files []File, form *url.Values) (body io.Reader, string, error) {
	bodyBuffer := new(bytes.Buffer)
	bodyWriter := multipart.NewWriter(bodyBuffer)
	defer bodyWriter.Close()

	for _, file := range files {
		fileWriter, err := bodyWriter.CreateFormFile(file.Fieldname, file.Filename)
		if err != nil {
			return nli, "", err
		}

		f, err := os.Open(file.Filename)
		if err != nil {
			return nli, "", err
		}
		defer f.Close()

		_, err = io.Copy(fileWriter, f)
		if err != nil {
			return nli, "", err
		}
	}

	if form != nil {
		for k, vs := range form {
			for _, v := range vs {
				bodyWriter.WriteField(k, v)
			}
		}
	}

	return bodyBuffer, bodyWriter.FormDataContentType(), nil
}

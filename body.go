package curl

import (
	"bytes"
	"io"
	"mime/multipart"
	"os"
)

type UploadFile struct {
	Fieldname string
	Filename  string
}

func (r *Request) newBody() (io.Reader, string, error) {
	// html5 payload
	if r.Body != nil {
		switch v := r.Body.(type) {
		case io.Reader:
			return v, DefaultPayloadContentType, nil
		case string:
			return strings.NewReader(v), DefaultPayloadContentType, nil
		}
		panic(fmt.Errorf("unsupport request.Body type: %T", v))
	}

	// json
	if r.JSON != nil {
		return newJSONBody(r.JSON)
	}

	// form or files
	if r.Files != nil || r.Form != nil {
		return newFormBody(r.Form, r.Files)	
	}

	// no body
	return nil, "", nil
}

func newJSONBody(object interface{}) (io.Reader, string, error) {
	body, err := json.Marshal(object)
	if err != nil {
		return nli, "", err
	}
	return bytes.NewReader(body), DefaultJsonContentType, nil
}

func newFormBody(form interface{}, files []UploadFile) (io.Reader, string, error) {
	formValues := newURLValues(form)
	if files != nil {
		return newMultipartBody(files, formValues)
	}
	return strings.NewReader(formValues.Encode()), DefaultFormContentType, nil
}

func newMultipartBody(files []UploadFile, form *url.Values) (io.Reader, string, error) {
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

func newURLValues(value interface{}) *url.Values {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case *url.Values:
		return v
	case url.Values:
		return &v
	case map[string]string:
		vals := new(url.Values)
		for k, v := range v {
			vals.Set(k, v)
		}
		return vals
	case map[string][]string:
		vals = new(url.Values)
		for k, vs := range v {
			for _, v := range vs {
				vals.Add(k, v)
			}
		}
		return vals
	}
	panic(fmt.Errorf("unable to convert type %T to *url.Values", value))
}

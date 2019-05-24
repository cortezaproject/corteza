package store

import (
	"io"
	"mime/multipart"
	"net/url"

	"github.com/pkg/errors"

	"github.com/cortezaproject/corteza-server/internal/http"
)

func FromURL(fileURL string) (io.ReadCloser, error) {
	if u, err := url.ParseRequestURI(fileURL); err != nil {
		return nil, errors.WithStack(err)
	} else if u.Scheme != "https" {
		return nil, errors.New("Only HTTPS is supported for file uploads")
	}

	client, err := http.New(&http.Config{
		Timeout: 10,
	})

	if err != nil {
		return nil, errors.WithStack(err)
	}

	req, err := client.Get(fileURL)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func FromMultipartFile(file *multipart.FileHeader) (io.ReadCloser, error) {
	reader, err := file.Open()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return reader, nil
}

func FromAny(file *multipart.FileHeader, url string) (io.ReadCloser, error) {
	if file != nil {
		return FromMultipartFile(file)
	}
	return FromURL(url)
}

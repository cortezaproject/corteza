package types

import (
	"io"
)

type (
	HttpRequest struct {
		Method   string              `json:"method"`
		URL      *Url                `json:"url"`
		Header   map[string][]string `json:"header"`
		Body     io.Reader           `json:"body"`
		Form     map[string][]string `json:"form"`
		PostForm map[string][]string `json:"post_form"`
	}

	Url struct {
		Scheme      string
		Opaque      string
		Host        string
		Path        string
		RawPath     string
		ForceQuery  bool
		RawQuery    string
		Fragment    string
		RawFragment string
	}
)

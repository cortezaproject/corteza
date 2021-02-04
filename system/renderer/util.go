package renderer

import (
	"encoding/base64"
	htpl "html/template"
	ttpl "html/template"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/Masterminds/sprig"
)

func preprocHTMLTemplate(pl *driverPayload) (*htpl.Template, error) {
	bb, err := ioutil.ReadAll(pl.Template)
	if err != nil {
		return nil, err
	}

	gtpl := htpl.New("text/html_render").
		Funcs(sprig.FuncMap()).
		Funcs(htpl.FuncMap{
			// "attachDataURL": func(name string) htpl.URL {
			// 	// Find the attachment
			// 	att, has := pl.Attachments[name]
			// 	if !has {
			// 		return htpl.URL(fmt.Sprintf("error: attachment not found: %s", name))
			// 	}

			// 	// Process source
			// 	bb, err := ioutil.ReadAll(att.Source)
			// 	if err != nil {
			// 		return htpl.URL(fmt.Sprintf("error: %s", err.Error()))
			// 	}

			// 	return htpl.URL("data:" + att.Mime + ";base64," + base64.RawStdEncoding.EncodeToString(bb))
			// },

			"inlineRemote": func(url string) (htpl.URL, error) {
				rsp, err := http.Get(url)
				if err != nil {
					return "", err
				}

				defer rsp.Body.Close()

				bb, err := ioutil.ReadAll(rsp.Body)
				if err != nil {
					return "", err
				}

				raw := base64.RawStdEncoding.EncodeToString(bb)
				return htpl.URL("data:" + rsp.Header.Get("Content-Type") + ";base64," + raw), nil
			},
		})

	// Prep the original template
	t, err := gtpl.Parse(string(bb))
	if err != nil {
		return nil, err
	}

	// Prep partials
	for _, p := range pl.Partials {
		bb, err = ioutil.ReadAll(p)
		if err != nil {
			return nil, err
		}

		t, err = gtpl.Parse(string(bb))
		if err != nil {
			return nil, err
		}
	}

	return t, nil
}

func preprocPlainTemplate(tpl io.Reader, pp map[string]io.Reader) (*ttpl.Template, error) {
	bb, err := ioutil.ReadAll(tpl)
	if err != nil {
		return nil, err
	}

	gtpl := ttpl.New("text/plain_render")

	// Prep the original template
	t, err := gtpl.Parse(string(bb))
	if err != nil {
		return nil, err
	}

	// Prep partials
	for _, p := range pp {
		bb, err = ioutil.ReadAll(p)
		if err != nil {
			return nil, err
		}

		t, err = gtpl.Parse(string(bb))
		if err != nil {
			return nil, err
		}
	}

	return t, nil
}

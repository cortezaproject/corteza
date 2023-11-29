package renderer

import (
	"encoding/base64"
	"html/template"
	"io"
	"net/http"

	"github.com/Masterminds/sprig"
)

func preprocHTMLTemplate(pl *driverPayload) (*template.Template, error) {
	bb, err := io.ReadAll(pl.Template)
	if err != nil {
		return nil, err
	}

	gtpl := template.New("text/html_render").
		Funcs(sprig.FuncMap()).
		Funcs(template.FuncMap{
			// "attachDataURL": func(name string) template.URL {
			// 	// Find the attachment
			// 	att, has := pl.Attachments[name]
			// 	if !has {
			// 		return template.URL(fmt.Sprintf("error: attachment not found: %s", name))
			// 	}

			// 	// Process source
			// 	bb, err := io.ReadAll(att.Source)
			// 	if err != nil {
			// 		return template.URL(fmt.Sprintf("error: %s", err.Error()))
			// 	}

			// 	return template.URL("data:" + att.Mime + ";base64," + base64.RawStdEncoding.EncodeToString(bb))
			// },

			"inlineRemote": func(url string) (template.URL, error) {
				rsp, err := http.Get(url)
				if err != nil {
					return "", err
				}

				defer rsp.Body.Close()

				bb, err := io.ReadAll(rsp.Body)
				if err != nil {
					return "", err
				}

				raw := base64.RawStdEncoding.EncodeToString(bb)
				return template.URL("data:" + rsp.Header.Get("Content-Type") + ";base64," + raw), nil
			},
			"env": envGetter(),
		})

	// Prep the original template
	t, err := gtpl.Parse(string(bb))
	if err != nil {
		return nil, err
	}

	// Prep partials
	for _, p := range pl.Partials {
		bb, err = io.ReadAll(p)
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

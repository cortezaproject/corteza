package renderer

import (
	"encoding/base64"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"

	"github.com/Masterminds/sprig"
)

func preprocHTMLTemplate(pl *driverPayload) (*template.Template, error) {
	bb, err := ioutil.ReadAll(pl.Template)
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
			// 	bb, err := ioutil.ReadAll(att.Source)
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

				bb, err := ioutil.ReadAll(rsp.Body)
				if err != nil {
					return "", err
				}

				raw := base64.RawStdEncoding.EncodeToString(bb)
				return template.URL("data:" + rsp.Header.Get("Content-Type") + ";base64," + raw), nil
			},
			// emits the value as-is, without the contextual escaping
			// html/template applies to every other interpolated value;
			// unset values render empty, the same as plain interpolation
			"safeHTML": func(v any) template.HTML {
				if v == nil {
					return ""
				}

				return template.HTML(fmt.Sprintf("%v", v))
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

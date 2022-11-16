package renderer

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/cortezaproject/corteza/server/system/types"
)

type (
	gotenbergPDF struct {
		url string
		def DriverDefinition
	}
	gotenbergPDFDriver struct {
		url string
	}
)

// @todo healthcheck, different input data formats
func newGotenbergPDF(url string) driverFactory {
	return &gotenbergPDF{
		url: url,

		def: DriverDefinition{
			Name: "gotenbergPDF",
			InputTypes: []types.DocumentType{
				types.DocumentTypePlain,
				types.DocumentTypeHTML,
			},
			OutputTypes: []types.DocumentType{
				types.DocumentTypePDF,
			},
		},
	}
}

func (d *gotenbergPDF) Define() DriverDefinition {
	return d.def
}

func (d *gotenbergPDF) CanRender(t types.DocumentType) bool {
	for _, i := range d.def.InputTypes {
		if i == t {
			return true
		}
	}
	return false
}

func (d *gotenbergPDF) CanProduce(t types.DocumentType) bool {
	for _, o := range d.def.OutputTypes {
		if o == t {
			return true
		}
	}
	return false
}

func (d *gotenbergPDF) Driver() driver {
	return &gotenbergPDFDriver{
		url: d.url,
	}
}

func (d *gotenbergPDFDriver) Render(ctx context.Context, pl *driverPayload) (io.ReadSeeker, error) {
	// HTTP request body stuff
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	// index.html is required by the rendering container
	part, err := writer.CreateFormFile("file", "index.html")
	if err != nil {
		return nil, err
	}
	err = d.prepareContent(part, pl)
	if err != nil {
		return nil, err
	}

	// Document configurations
	err = d.applyOptions(writer, pl.Options)
	if err != nil {
		return nil, err
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	// HTTP request header stuff
	// @todo make sure to use the propper endpoint when you add support
	//       for different inputs.
	url := d.url
	if n, has := pl.Options["url"]; has {
		url = n
	}
	req, err := http.NewRequest("POST", url+"/convert/html", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	ss, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return bytes.NewReader(ss), err
}

func (d *gotenbergPDFDriver) prepareContent(w io.Writer, pl *driverPayload) error {
	t, err := preprocHTMLTemplate(pl)
	if err != nil {
		return err
	}

	return t.Execute(w, pl.Variables)
}

func (d *gotenbergPDFDriver) applyOptions(mw *multipart.Writer, opts map[string]string) (err error) {
	if opts == nil {
		return nil
	}

	for k, v := range opts {
		switch k {
		case "marginTop",
			"marginBottom",
			"marginLeft",
			"marginRight":
			err = d.addFormField(mw, k, v)

		case "marginY":
			err = d.addFormField(mw, "marginTop", v)
			if err != nil {
				return err
			}
			err = d.addFormField(mw, "marginBottom", v)

		case "marginX":
			err = d.addFormField(mw, "marginLeft", v)
			if err != nil {
				return err
			}
			err = d.addFormField(mw, "marginRight", v)

		case "margin":
			err = d.addFormField(mw, "marginTop", v)
			if err != nil {
				return err
			}
			err = d.addFormField(mw, "marginBottom", v)
			if err != nil {
				return err
			}
			err = d.addFormField(mw, "marginLeft", v)
			if err != nil {
				return err
			}
			err = d.addFormField(mw, "marginRight", v)

		case "documentSize":
			w, h := d.documentDimensions(v)
			if w+h != "" {
				err = d.addFormField(mw, "paperWidth", w)
				if err != nil {
					return err
				}
				err = d.addFormField(mw, "paperHeight", h)
			}

		case "documentWidth":
			err = d.addFormField(mw, "paperWidth", v)

		case "documentHeight":
			err = d.addFormField(mw, "paperHeight", v)

		case "contentScale":
			err = d.addFormField(mw, "scale", v)

		case "orientation":
			if v == "landscape" {
				err = d.addFormField(mw, "landscape", "true")
			} else {
				err = d.addFormField(mw, "landscape", "false")
			}
		}

		if err != nil {
			return err
		}
	}
	return nil
}

func (d *gotenbergPDFDriver) addFormField(mw *multipart.Writer, k, v string) error {
	w, err := mw.CreateFormField(k)
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(v))
	return err
}

// documentDimensions returns the ISO216 standard document dimensions in inches (Gotenberg uses inches)
func (d *gotenbergPDFDriver) documentDimensions(isoDoc string) (string, string) {
	switch strings.ToLower(isoDoc) {
	// A series
	case "a0":
		return "33.1", "46.8"
	case "a1":
		return "23.4", "33.1"
	case "a2":
		return "16.5", "23.4"
	case "a3":
		return "11.7", "16.5"
	case "a4":
		return "8.3", "11.7"
	case "a5":
		return "5.8", "8.3"
	case "a6":
		return "4.1", "5.8"
	case "a7":
		return "2.9", "4.1"
	case "a8":
		return "2.0", "2.9"
	case "a9":
		return "1.5", "2.0"
	case "a10":
		return "1.0", "1.5"

	// B series
	case "b0":
		return "39.4", "55.7"
	case "b1":
		return "27.8", "39.4"
	case "b2":
		return "19.7", "27.8"
	case "b3":
		return "13.9", "19.7"
	case "b4":
		return "9.8", "13.9"
	case "b5":
		return "6.9", "9.8"
	case "b6":
		return "4.9", "6.9"
	case "b7":
		return "3.5", "4.9"
	case "b8":
		return "2.4", "3.5"
	case "b9":
		return "1.7", "2.4"
	case "b10":
		return "1.2", "1.7"

	// C series
	case "c0":
		return "36.1", "51.1"
	case "c1":
		return "25.5", "36.1"
	case "c2":
		return "18.0", "25.5"
	case "c3":
		return "12.8", "18.0"
	case "c4":
		return "9.0", "12.8"
	case "c5":
		return "6.4", "9.0"
	case "c6":
		return "4.5", "6.4"
	case "c7":
		return "3.2", "4.5"
	case "c8":
		return "2.2", "3.2"
	case "c9":
		return "1.6", "2.2"
	case "c10":
		return "1.1", "1.6"

		// ANSI
	case "ansi a":
		return "8.5", "11"
	case "ansi b":
		return "11", "17"
	case "ansi c":
		return "17", "22"
	case "ansi d":
		return "22", "34"
	case "ansi e":
		return "34", "44"

	// Proprietary NA
	case "junior legal":
		return "8", "5"
	case "letter":
		return "8.5", "11"
	case "legal":
		return "8.5", "14"
	case "tabloid":
		return "11", "17"
	}

	// This will fallback to the default (A4)
	return "", ""
}

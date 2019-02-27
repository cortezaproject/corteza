package handlers

import (
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/crusttech/crust/crm/rest/request"
)

type Downloadable interface {
	Name() string
	Download() bool
	ModTime() time.Time
	Content() io.ReadSeeker
	Valid() bool
}

func NewAttachmentDownloadable(ctrl AttachmentAPI) *Attachment {
	h := NewAttachment(ctrl)
	h.Original = func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		params := request.NewAttachmentOriginal()
		params.Fill(r)

		f, err := ctrl.Original(r.Context(), params)
		serveFile(f, err, w, r)
	}

	h.Preview = func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		params := request.NewAttachmentPreview()
		params.Fill(r)

		f, err := ctrl.Preview(r.Context(), params)
		serveFile(f, err, w, r)
	}

	return h
}

func serveFile(f interface{}, err error, w http.ResponseWriter, r *http.Request) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else if dl, ok := f.(Downloadable); ok {
		if !dl.Valid() {
			w.WriteHeader(http.StatusNotFound)
		} else {
			if dl.Download() {
				w.Header().Add("Content-Disposition", "attachment; filename="+url.QueryEscape(dl.Name()))
			} else {
				w.Header().Add("Content-Disposition", "inline; filename="+url.QueryEscape(dl.Name()))
			}

			http.ServeContent(w, r, dl.Name(), dl.ModTime(), dl.Content())
		}
	} else {
		http.Error(w, "Got incompatible type from controller", http.StatusInternalServerError)
	}
}

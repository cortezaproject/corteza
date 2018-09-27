package handlers

import (
	"net/http"

	"github.com/crusttech/crust/sam/rest/request"
	"io"
	"net/url"
	"time"
)

// HTTP API interface
type AttachmentDownloadable struct {
	Original func(http.ResponseWriter, *http.Request)
	Preview  func(http.ResponseWriter, *http.Request)
}

type Downloadable interface {
	Name() string
	Download() bool
	ModTime() time.Time
	Content() io.ReadSeeker
}

func NewAttachmentDownloadable(ah AttachmentAPI) *Attachment {
	serve := func(f interface{}, err error, w http.ResponseWriter, r *http.Request) {
		if err != nil {
			switch true {
			case err.Error() == "crust.sam.repository.AttachmentNotFound":
				w.WriteHeader(http.StatusNotFound)
			default:
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else if dl, ok := f.(Downloadable); ok {
			if dl.Download() {
				w.Header().Add("Content-Disposition", "attachment; filename="+url.QueryEscape(dl.Name()))
			}

			http.ServeContent(w, r, dl.Name(), dl.ModTime(), dl.Content())
		} else {
			http.Error(w, "Got incompatible type from controller", http.StatusInternalServerError)
		}
	}

	return &Attachment{
		Original: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewAttachmentOriginal()
			params.Fill(r)

			f, err := ah.Original(r.Context(), params)
			serve(f, err, w, r)
		},

		Preview: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewAttachmentPreview()
			params.Fill(r)

			f, err := ah.Preview(r.Context(), params)
			serve(f, err, w, r)
		},
	}
}

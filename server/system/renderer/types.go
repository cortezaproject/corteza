package renderer

import (
	"context"
	"io"

	"github.com/cortezaproject/corteza/server/system/types"
)

type (
	RendererPayload struct {
		Template     io.Reader
		TemplateType types.DocumentType
		TargetType   types.DocumentType
		Variables    map[string]interface{}
		Options      map[string]string
		Partials     []*TemplatePartial
		Attachments  AttachmentIndex
	}

	TemplatePartial struct {
		Handle       string
		Template     io.Reader
		TemplateType types.DocumentType
	}

	driverPayload struct {
		Template    io.Reader
		Variables   map[string]interface{}
		Options     map[string]string
		Partials    map[string]io.Reader
		Attachments AttachmentIndex
	}

	AttachmentIndex map[string]*Attachment
	Attachment      struct {
		Source io.Reader
		Mime   string
		Name   string
	}

	DriverDefinition struct {
		Name        string               `json:"name"`
		InputTypes  []types.DocumentType `json:"inputTypes"`
		OutputTypes []types.DocumentType `json:"outputTypes"`
	}

	driverFactory interface {
		Define() DriverDefinition

		CanRender(t types.DocumentType) bool
		CanProduce(t types.DocumentType) bool
		Driver() driver
	}
	driver interface {
		Render(ctx context.Context, pl *driverPayload) (io.ReadSeeker, error)
	}
)

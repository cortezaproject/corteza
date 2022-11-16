package automation

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/cortezaproject/corteza/server/compose/types"
)

type (
	attachmentService interface {
		FindByID(ctx context.Context, namespaceID, attachmentID uint64) (*types.Attachment, error)
		CreateRecordAttachment(ctx context.Context, namespaceID uint64, name string, size int64, fh io.ReadSeeker, moduleID, recordID uint64, fieldName string) (att *types.Attachment, err error)
		DeleteByID(ctx context.Context, namespaceID uint64, attachmentID uint64) error
		OpenOriginal(att *types.Attachment) (io.ReadSeekCloser, error)
		OpenPreview(att *types.Attachment) (io.ReadSeekCloser, error)
	}

	attachmentHandler struct {
		reg modulesHandlerRegistry
		svc attachmentService
	}

	attachmentLookup interface {
		GetAttachment() (bool, uint64, *types.Attachment)
	}
)

func AttachmentHandler(reg modulesHandlerRegistry, svc attachmentService) *attachmentHandler {
	h := &attachmentHandler{
		reg: reg,
		svc: svc,
	}

	h.register()
	return h
}

func (h attachmentHandler) lookup(ctx context.Context, args *attachmentLookupArgs) (results *attachmentLookupResults, err error) {
	results = &attachmentLookupResults{}
	results.Attachment, err = h.svc.FindByID(ctx, 0, args.Attachment)
	return
}

func (h attachmentHandler) delete(ctx context.Context, args *attachmentDeleteArgs) error {
	return h.svc.DeleteByID(ctx, 0, args.Attachment)
}

func (h attachmentHandler) openOriginal(ctx context.Context, args *attachmentOpenOriginalArgs) (*attachmentOpenOriginalResults, error) {
	att, err := lookupAttachment(ctx, h.svc, args)
	if err != nil {
		return nil, err
	}

	r := &attachmentOpenOriginalResults{}
	r.Content, err = h.svc.OpenOriginal(att)
	if err != nil {
		return nil, err
	}

	// @todo we need to call Close() when file is read (or at the end of the workflow)
	//       some kind of workflow-cleanup facility is needed

	return r, nil
}

func (h attachmentHandler) openPreview(ctx context.Context, args *attachmentOpenPreviewArgs) (*attachmentOpenPreviewResults, error) {
	att, err := lookupAttachment(ctx, h.svc, args)
	if err != nil {
		return nil, err
	}

	r := &attachmentOpenPreviewResults{}
	r.Content, err = h.svc.OpenOriginal(att)
	if err != nil {
		return nil, err
	}

	// @todo we need to call Close() when file is read (or at the end of the workflow)
	//       some kind of workflow-cleanup facility is needed

	return r, nil
}

func (h attachmentHandler) create(ctx context.Context, args *attachmentCreateArgs) (*attachmentCreateResults, error) {
	var (
		err error
		att *types.Attachment

		fh   io.ReadSeeker
		size int64
	)

	switch {
	case len(args.contentBytes) > 0:
		size = int64(len(args.contentBytes))
		fh = bytes.NewReader(args.contentBytes)

	case args.contentStream != nil:
		if rs, is := args.contentStream.(io.ReadSeeker); is {
			_, err = rs.Seek(0, 0)
			if err != nil {
				return nil, err
			}

			size, err = getReaderSize(rs)
			if err != nil {
				return nil, err
			}

			_, err = rs.Seek(0, 0)
			if err != nil {
				return nil, err
			}

			fh = rs
		} else {
			// In case we only got a reader...
			//
			// For future proofing, for handling larger attachment, we create a temp.
			// file which we then use as a reader

			// Preparations
			tmpf, err := ioutil.TempFile("", "reader")
			if err != nil {
				return nil, err
			}
			defer tmpf.Close()
			defer os.Remove(tmpf.Name())

			// Writing content to file
			w := bufio.NewWriter(tmpf)
			r := args.contentStream
			buf := make([]byte, 1024)
			for {
				// read
				n, err := r.Read(buf)
				if err != nil && err != io.EOF {
					return nil, err
				}
				if n == 0 {
					break
				}

				// on-the-fly size calculation
				size += int64(n)

				// write
				if _, err := w.Write(buf[:n]); err != nil {
					return nil, err
				}
			}
			if err = w.Flush(); err != nil {
				return nil, err
			}

			_, err = tmpf.Seek(0, 0)
			if err != nil {
				return nil, err
			}
			fh = tmpf
		}

	default:
		fh = strings.NewReader(args.contentString)
		size = int64(len(args.contentString))
	}

	switch {
	case args.Resource != nil:
		att, err = h.svc.CreateRecordAttachment(
			ctx,
			args.Resource.NamespaceID,
			args.Name,
			size,
			fh,
			args.Resource.ModuleID,
			args.Resource.ID,
			args.FieldName,
		)
	default:
		return nil, fmt.Errorf("unknown resource")
	}

	if err != nil {
		return nil, err
	}

	return &attachmentCreateResults{att}, nil
}

func lookupAttachment(ctx context.Context, svc attachmentService, args attachmentLookup) (*types.Attachment, error) {
	_, ID, attachment := args.GetAttachment()

	switch {
	case attachment != nil:
		return attachment, nil
	case ID > 0:
		return svc.FindByID(ctx, 0, ID)
	}

	return nil, fmt.Errorf("empty attachment lookup params")
}

func getReaderSize(r io.Reader) (size int64, err error) {
	buf := make([]byte, 1024)
	var n int
	for {
		// read a chunk
		n, err = r.Read(buf)
		if err != nil && err != io.EOF {
			return 0, err
		}
		if n == 0 {
			break
		}
		size += int64(n)
	}

	return size, nil
}

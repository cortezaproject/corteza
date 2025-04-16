package automation

import (
	"context"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
)

type (
	apigwBodyHandler struct {
		reg apigwBodyHandlerRegistry
	}
)

func ApigwBodyHandler(reg queueHandlerRegistry) *apigwBodyHandler {
	h := &apigwBodyHandler{
		reg: reg,
	}

	h.register()
	return h
}

func (h apigwBodyHandler) read(ctx context.Context, args *apigwBodyReadArgs) (res *apigwBodyReadResults, err error) {
	res = &apigwBodyReadResults{}

	bb, err := io.ReadAll(args.Request.Body)

	if err != nil {
		return
	}

	res.Body = string(bb)

	return
}

func (h apigwBodyHandler) readFile(ctx context.Context, args *apigwBodyReadFileArgs) (res *apigwBodyReadFileResults, err error) {
	// @note the contents are already parsed; will probably need to rework behind the scenes bits
	res = &apigwBodyReadFileResults{
		Exists: true,
	}

	var fh *multipart.FileHeader
	res.File, fh, err = args.Request.FormFile(args.Name)
	if err != nil {
		// In case the file is missing or empty, do this
		// @todo do we want to just check if missing?
		if errors.Is(err, http.ErrMissingFile) || strings.Contains(err.Error(), "EOF") {
			res.Exists = false
			err = nil
		} else {
			return
		}
	}

	if !res.Exists {
		return
	}

	res.FileName = fh.Filename

	return
}

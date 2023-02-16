package envoyx

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func (nvyx *service) decodeUri(ctx context.Context, p DecodeParams) (nn NodeSet, err error) {
	aUri, ok := p.Params["uri"]
	if !ok {
		err = fmt.Errorf("cannot decode URI: no uri parameter provided")
		return
	}

	uri, ok := aUri.(string)
	if !ok {
		err = fmt.Errorf("cannot decode URI: uri should be string encoded: got %v", aUri)
		return
	}

	pp := strings.Split(uri, "://")

	proto := pp[0]
	rest := pp[1]

	switch proto {
	case "file":
		fileInfo, err := os.Stat(rest)
		if err != nil {
			return nil, err
		}

		if fileInfo.IsDir() {
			// decode whole directory
			return nvyx.decodeDirectory(ctx, p, rest)
		}

		// decode specific file
		return nvyx.decodeFile(ctx, p, rest)

	default:
		err = fmt.Errorf("unsupported URI protocol %s", proto)
		return
	}
}

func (nvyx *service) encodeIo(ctx context.Context, dg *depGraph, p EncodeParams) (err error) {
	for rt, nn := range NodesByResourceType(dg.Roots()...) {
		for _, se := range nvyx.encoders[EncodeTypeIo] {
			err = se.Encode(ctx, p, rt, OmitPlaceholderNodes(nn...), dg)
			if err != nil {
				return
			}
		}
	}

	return
}

func (nvyx *service) decodeDirectory(ctx context.Context, p DecodeParams, path string) (nn NodeSet, err error) {
	return nn, filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// @todo consider supporting nested directories
		if info.IsDir() {
			return nil
		}

		aux, err := nvyx.decodeFile(ctx, p, path)
		if err != nil {
			return err
		}

		nn = append(nn, aux...)
		return nil
	})
}

func (nvyx *service) decodeFile(ctx context.Context, p DecodeParams, path string) (nn NodeSet, err error) {
	var aux NodeSet
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	p.Params["stream"] = f

	for _, d := range nvyx.decoders[DecodeTypeURI] {
		aux, err = d.Decode(ctx, p)
		if err != nil {
			return
		}
		nn = append(nn, aux...)

		_, err = f.Seek(0, 0)
		if err != nil {
			return
		}
	}

	return
}

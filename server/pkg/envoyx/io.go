package envoyx

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/cortezaproject/corteza/server/pkg/envoyx/csv"
	"github.com/cortezaproject/corteza/server/pkg/envoyx/json"
)

func (svc *Service) decodeIo(ctx context.Context, p DecodeParams) (nodes NodeSet, providers []Provider, err error) {
	// Get and validate the reader param
	r, ok := p.Params["reader"]
	if !ok {
		err = fmt.Errorf("cannot decode IO: no reader parameter provided")
		return
	}
	reader, ok := r.(io.Reader)
	if !ok {
		err = fmt.Errorf("cannot decode IO: reader should be io.Reader encoded")
		return
	}

	// Get and validate the mimetype param
	_, ok = p.Params["mime"].(string)
	if !ok {
		// @todo guess mime type
		p.Params["mime"] = "text/yaml"
	}

	// @todo different ext based on mimetype -- currently yaml is the only supported one
	tmpFileName := "*.yaml"
	f, err := os.CreateTemp(os.TempDir(), tmpFileName)
	if err != nil {
		return
	}
	defer f.Close()
	defer os.Remove(f.Name())

	_, err = io.Copy(f, reader)
	if err != nil {
		return
	}

	return svc.decodeFile(ctx, p, os.TempDir(), f.Name())
}

func (svc *Service) decodeUri(ctx context.Context, p DecodeParams) (nodes NodeSet, providers []Provider, err error) {
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
			return nil, nil, err
		}

		if fileInfo.IsDir() {
			// decode whole directory
			return svc.decodeDirectory(ctx, p, rest)
		}

		// decode specific file
		var base string
		bits := strings.Split(rest, string(os.PathSeparator))
		base = strings.Join(bits[0:len(bits)-1], string(os.PathSeparator))
		return svc.decodeFile(ctx, p, base, rest)

	default:
		err = fmt.Errorf("unsupported URI protocol %s", proto)
		return
	}
}

func (nvyx *Service) encodeIo(ctx context.Context, dg *DepGraph, p EncodeParams) (err error) {
	for rt, nn := range NodesByResourceType(dg.Roots()...) {
		for _, se := range nvyx.encoders[EncodeTypeIo] {
			nn = OmitPlaceholderNodes(nn...)
			if len(nn) == 0 {
				continue
			}

			err = se.Encode(ctx, p, rt, OmitPlaceholderNodes(nn...), dg)
			if err != nil {
				return
			}
		}
	}

	return
}

func (nvyx *Service) decodeDirectory(ctx context.Context, p DecodeParams, path string) (nodes NodeSet, providers []Provider, err error) {
	basePath := path
	return nodes, providers, filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// @todo consider supporting nested directories
		if info.IsDir() {
			return nil
		}

		auxNodes, auxProviders, err := nvyx.decodeFile(ctx, p, basePath, path)
		if err != nil {
			return err
		}

		nodes = append(nodes, auxNodes...)
		providers = append(providers, auxProviders...)
		return nil
	})
}

func (nvyx *Service) decodeFile(ctx context.Context, p DecodeParams, basePath, path string) (nodes NodeSet, providers []Provider, err error) {
	var aux NodeSet
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	p.Params["stream"] = f

	for _, d := range nvyx.decoders[DecodeTypeURI] {
		if cc, ok := d.(canCheckFile); ok {
			if !cc.CanFile(f) {
				continue
			}
		}

		aux, err = d.Decode(ctx, p)
		if err != nil {
			return
		}
		nodes = append(nodes, aux...)

		_, err = f.Seek(0, 0)
		if err != nil {
			return
		}
	}

	providerIdent := strings.Replace(f.Name(), basePath, "", -1)
	providerIdent = strings.TrimLeft(providerIdent, "/")

	if csv.CanDecodeFile(f) || csv.CanDecodeExt(f.Name()) {
		f.Seek(0, 0)
		aux, err := csv.Decoder(f, providerIdent)
		if err != nil {
			return nil, nil, err
		}
		providers = append(providers, aux)
	}

	if json.CanDecodeFile(f) || json.CanDecodeExt(f.Name()) {
		f.Seek(0, 0)
		aux, err := json.Decoder(f, providerIdent)
		if err != nil {
			return nil, nil, err
		}
		providers = append(providers, aux)
	}

	return
}

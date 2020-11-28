package util

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	es "github.com/cortezaproject/corteza-server/pkg/envoy/store"
	"github.com/cortezaproject/corteza-server/pkg/envoy/yaml"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/goware/statik/fs"
)

type (
	decoder interface {
		CanDecodeExt(string) bool
		CanDecodeFile(io.Reader) bool
		Decode(context.Context, io.Reader, *envoy.DecoderOpts) ([]resource.Interface, error)
	}
)

func StatikResources(ctx context.Context, d decoder, sfs http.FileSystem, root string) ([]resource.Interface, error) {
	nn := make(resource.InterfaceSet, 0, 100)
	err := fs.Walk(sfs, root, func(p string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		f, err := sfs.Open(p)
		if err != nil {
			return err
		}
		defer f.Close()

		if !d.CanDecodeFile(f) && !d.CanDecodeExt(info.Name()) {
			return nil
		}

		if d == nil {
			// no decoder found
			return nil
		}

		_, err = f.Seek(0, 0)
		if err != nil {
			return err
		}

		dir, fn := path.Split(p)
		do := &envoy.DecoderOpts{
			Name: fn,
			Path: dir,
		}

		if dnn, err := d.Decode(ctx, f, do); err != nil {
			return fmt.Errorf("failed to decode %s: %w", info.Name(), err)
		} else {
			nn = append(nn, dnn...)
		}

		return nil
	})

	return nn, err
}

func EncodeStatik(ctx context.Context, s store.Storer, rfs string, root string) error {
	sfs, err := fs.New(rfs)
	if err != nil {
		return err
	}

	yd := yaml.Decoder()
	se := es.NewStoreEncoder(s, &es.EncoderConfig{
		OnExisting: es.MergeLeft,
	})

	nn, err := StatikResources(ctx, yd, sfs, root)
	if err != nil {
		return err
	}

	bld := envoy.NewBuilder(se)
	g, err := bld.Build(ctx, nn...)
	if err != nil {
		return err
	}

	return envoy.Encode(ctx, g, se)
}

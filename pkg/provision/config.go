package provision

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/directory"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	es "github.com/cortezaproject/corteza-server/pkg/envoy/store"
	"github.com/cortezaproject/corteza-server/pkg/envoy/yaml"
	"github.com/cortezaproject/corteza-server/store"
	"go.uber.org/zap"
	"path/filepath"
)

func importConfig(ctx context.Context, log *zap.Logger, s store.Storer) error {
	var (
		yd  = yaml.Decoder()
		nn  = make([]resource.Interface, 0, 200)
		se  = es.NewStoreEncoder(s, &es.EncoderConfig{OnExisting: es.Skip})
		bld = envoy.NewBuilder(se)

		pp, err = filepath.Glob("provision/*")
	)

	if err != nil {
		return err
	}

	for _, path := range pp {
		log.Info("provisioning from path", zap.String("path", path))
		if mm, err := directory.Decode(ctx, path, yd); err != nil {
			return err
		} else {
			nn = append(nn, mm...)
		}
	}

	if g, err := bld.Build(ctx, nn...); err != nil {
		return err
	} else if err = envoy.Encode(ctx, g, se); err != nil {
		return err
	}

	return nil
}

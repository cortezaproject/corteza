package messaging

import (
	"context"
	"io"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"

	"github.com/cortezaproject/corteza-server/messaging/importer"
	"github.com/cortezaproject/corteza-server/messaging/service"
	"github.com/cortezaproject/corteza-server/messaging/types"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	impAux "github.com/cortezaproject/corteza-server/pkg/importer"
	"github.com/cortezaproject/corteza-server/pkg/settings"
	provision "github.com/cortezaproject/corteza-server/provision/messaging"
)

func provisionConfig(ctx context.Context, log *zap.Logger) (err error) {
	log.Debug("running configuration provision")

	var provisioned bool

	// Make sure we have all full access for provisioning
	ctx = auth.SetSuperUserContext(ctx)

	if provisioned, err = isProvisioned(ctx); err != nil {
		return err
	} else if provisioned {
		log.Debug("configuration already provisioned")
	}

	readers, err := impAux.ReadStatic(provision.Asset)
	if err != nil {
		return err
	}

	if provisioned {
		return partialImportSettings(ctx, service.DefaultSettings, readers...)
	}

	return errors.Wrap(
		importer.Import(ctx, readers...),
		"could not provision configuration for messaging service",
	)
}

// Provision ONLY when there are no channels (even if we find delete channels we abort provisioning
func isProvisioned(ctx context.Context) (bool, error) {
	cc, _, err := service.DefaultChannel.With(ctx).Find(types.ChannelFilter{IncludeDeleted: true})
	return len(cc) > 0, err
}

// Partial import of settings from provision files
func partialImportSettings(ctx context.Context, ss settings.Service, ff ...io.Reader) (err error) {
	var (
		// decoded content from YAML files
		aux interface{}

		si = settings.NewImporter()

		// importer w/o permissions & roles
		// we need only settings
		imp = importer.NewImporter(nil, si, nil)

		// current value
		current settings.ValueSet

		// unexisting values
		unex settings.ValueSet
	)

	for _, f := range ff {
		if err = yaml.NewDecoder(f).Decode(&aux); err != nil {
			return
		}

		err = imp.Cast(aux)
		if err != nil {
			return
		}
	}

	// Get all "current" settings storage
	current, err = ss.FindByPrefix(ctx)
	if err != nil {
		return
	}

	// Compare current settings with imported, get all that do not exist yet
	if unex = si.GetValues(); len(unex) > 0 {
		// Store non existing
		err = ss.BulkSet(ctx, current.New(unex))
		if err != nil {
			return
		}
	}

	return nil
}

package envoy

import (
	"context"
	"testing"

	su "github.com/cortezaproject/corteza-server/pkg/envoy/store"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/stretchr/testify/require"
)

func TestProvision_overwriting(t *testing.T) {
	var (
		ctx    = context.Background()
		s, err = initStore(ctx)
	)

	if err != nil {
		t.Fatalf("failed to init sqlite in-memory db: %v", err)
	}

	ni := uint64(0)
	su.NextID = func() uint64 {
		ni++
		return ni
	}

	prepare := func(ctx context.Context, s store.Storer, t *testing.T, suite string, cfg *su.EncoderConfig) (*require.Assertions, error) {
		req := require.New(t)

		nn, err := dd(ctx, suite)
		req.NoError(err)

		return req, encodeC(ctx, s, nn, cfg)
	}

	// Prepare
	s, err = initStore(ctx)
	err = ce(
		err,

		s.TruncateUsers(ctx),
		s.TruncateRoles(ctx),
		s.TruncateRbacRules(ctx),
		s.TruncateActionlogs(ctx),
		s.TruncateApplications(ctx),
		s.TruncateAttachments(ctx),
		s.TruncateComposeAttachments(ctx),
		s.TruncateComposeCharts(ctx),
		s.TruncateComposeNamespaces(ctx),
		s.TruncateComposeModules(ctx),
		s.TruncateComposeModuleFields(ctx),
		s.TruncateComposePages(ctx),
		s.TruncateComposeRecords(ctx, nil),

		storeRole(ctx, s, 1, "everyone"),
		storeRole(ctx, s, 2, "admins"),
	)
	if err != nil {
		t.Fatal(err.Error())
	}

	// Prepare
	req, err := prepare(ctx, s, t, "provision_batch/app_1", nil)
	req.NoError(err)
	store.TruncateComposeRecords(ctx, s, nil)

	req, err = prepare(ctx, s, t, "provision_batch/app_1", &su.EncoderConfig{OnExisting: su.Replace})
	req.NoError(err)
	checkBatchProvision(ctx, t, req, s, "ns1")
	store.TruncateComposeRecords(ctx, s, nil)

	req, err = prepare(ctx, s, t, "provision_batch/app_1", &su.EncoderConfig{OnExisting: su.Skip})
	req.NoError(err)
	checkBatchProvision(ctx, t, req, s, "ns1")
	store.TruncateComposeRecords(ctx, s, nil)

	req, err = prepare(ctx, s, t, "provision_batch/app_1", &su.EncoderConfig{OnExisting: su.MergeLeft})
	req.NoError(err)
	checkBatchProvision(ctx, t, req, s, "ns1")
	store.TruncateComposeRecords(ctx, s, nil)

	req, err = prepare(ctx, s, t, "provision_batch/app_1", &su.EncoderConfig{OnExisting: su.MergeRight})
	req.NoError(err)
	checkBatchProvision(ctx, t, req, s, "ns1")
	store.TruncateComposeRecords(ctx, s, nil)
}

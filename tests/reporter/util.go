package reporter

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"path"

	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/csv"
	"github.com/cortezaproject/corteza-server/pkg/envoy/directory"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	es "github.com/cortezaproject/corteza-server/pkg/envoy/store"
	"github.com/cortezaproject/corteza-server/pkg/envoy/yaml"
	"github.com/cortezaproject/corteza-server/pkg/report"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/service"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	auxReport struct {
		*types.Report

		Frames report.FrameDefinitionSet `json:"frames"`
	}

	valueDef map[string][]string
)

func bmReporterPrepareDM(ctx context.Context, h helper, s store.Storer, suite string) error {
	// cleanup the store
	h.noError(store.TruncateComposeNamespaces(ctx, service.DefaultStore))
	h.noError(store.TruncateComposeModules(ctx, service.DefaultStore))
	h.noError(store.TruncateComposeModuleFields(ctx, service.DefaultStore))
	h.noError(s.TruncateComposeRecords(ctx, nil))

	// prepare the params
	dmDef := path.Join("testdata", suite, "data_model")

	// build the datamodel
	yd := yaml.Decoder()
	cd := csv.Decoder()
	nn, err := directory.Decode(ctx, dmDef, yd, cd)

	crs := resource.ComposeRecordShaper()
	nn, err = resource.Shape(nn, crs)
	h.a.NoError(err)

	// import into the store
	se := es.NewStoreEncoder(s, nil)
	bld := envoy.NewBuilder(se)
	g, err := bld.Build(ctx, nn...)
	if err != nil {
		return err
	}
	err = envoy.Encode(ctx, g, se)
	if err != nil {
		return err
	}

	return nil
}

func bmReporterParseReport(path string) (*auxReport, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	aux := &auxReport{}
	raw, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(raw, &aux)
	if err != nil {
		return nil, err
	}

	return aux, nil
}

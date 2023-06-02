package bundler

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	automationEnvoy "github.com/cortezaproject/corteza/server/automation/envoy"
	"github.com/cortezaproject/corteza/server/automation/types"
	"github.com/cortezaproject/corteza/server/pkg/envoyx"
	"github.com/cortezaproject/corteza/server/pkg/plugin/discovery"
)

type (
	WfService interface {
		Create(context.Context, *types.Workflow) (*types.Workflow, error)
		Search(context.Context, types.WorkflowFilter) (types.WorkflowSet, types.WorkflowFilter, error)
	}

	fi struct {
		i os.FileInfo
		p string
	}

	wf struct {
		s WfService
		e *envoyx.Service
		w []*fi
	}
)

func Workflow(service WfService) *wf {
	evsvc := envoyx.New()
	evsvc.AddDecoder(envoyx.DecodeTypeIO,
		automationEnvoy.YamlDecoder{},
	)

	return &wf{
		s: service,
		e: evsvc,
	}
}

// Validate checks if any of the found subfolders of a bundle contains an automation workflow json file
func (w *wf) Validate(ctx context.Context, b []string) (byte, []string) {
	ss := []string{}

	for _, bb := range b {
		if strings.Contains(strings.ToLower(bb), "automation_workflow") {
			ga, _ := filepath.Glob(bb + "/*.json")
			ss = append(ss, ga...)

			// should we support yaml?

			// ga, _ := filepath.Glob(bb + "/*.yml")
			// ss = append(ss, ga...)

			// if len(ga) == 0 {
			// 	gb, _ := filepath.Glob(bb + "/*.yaml")
			// 	ss = append(ss, gb...)
			// }
		}
	}

	if len(ss) == 0 {
		return 0, ss
	}

	return discovery.AUTOMATION_WORKFLOW, ss
}

func (w *wf) Type() byte {
	return discovery.AUTOMATION_WORKFLOW
}

func (w *wf) Register(ctx context.Context, paths []string) (err error) {

	// json

	// type (
	// 	alias struct {
	// 		Wfs types.WorkflowSet `json:"workflows"`
	// 	}
	// )

	// for _, p := range paths {
	// 	r := &alias{}
	// 	bb, _ := os.ReadFile(p)

	// 	if err = json.Unmarshal(bb, r); err != nil {
	// 		continue
	// 	}

	// 	if len(r.Wfs) != 1 {
	// 		continue
	// 	}

	// 	// find workflow by the same handle
	// 	s, _, err := w.s.Search(ctx, types.WorkflowFilter{Handle: r.Wfs[0].Handle})

	// 	if err != nil {
	// 		continue
	// 	}

	// 	if len(s) == 1 {
	// 		// wf with same handle found, skip
	// 		continue
	// 	}

	// 	_, err = w.s.Create(ctx, r.Wfs[0])

	// 	spew.Dump(err)
	// }

	// yaml

	// for _, f := range w.w {
	// 	fullPath := filepath.Join(f.p, f.i.Name())

	// 	// read file
	// 	// b, err := os.ReadFile(fullPath)
	// 	rr, err := os.Open(fullPath)

	// 	if err != nil {
	// 		spew.Dump("could not read the file ", fullPath, err)
	// 		continue
	// 	}

	// 	ww := &types.Workflow{}

	// 	nodes, _, err := w.e.Decode(context.Background(), envoyx.DecodeParams{
	// 		Type: envoyx.DecodeTypeIO,
	// 		Params: map[string]any{
	// 			"reader": rr,
	// 			"mime":   "text/yaml",
	// 		},
	// 	})

	// 	if err != nil {
	// 		spew.Dump("err on envoy", err)
	// 		continue
	// 	}

	// 	var bb = bytes.NewBuffer(nil)

	// 	prms := envoyx.EncodeParams{
	// 		Type: envoyx.EncodeTypeIo,
	// 		Params: map[string]any{
	// 			"writer": bb,
	// 		},
	// 	}

	// 	gg, err := w.e.Bake(context.TODO(), prms, nil, nodes...)

	// 	spew.Dump("NODES", nodes)

	// 	if err != nil {
	// 		spew.Dump("1", err)
	// 		continue
	// 	}

	// 	err = w.e.Encode(context.TODO(), prms, gg)

	// 	if err != nil {
	// 		spew.Dump("encode err", err)
	// 	}

	// 	spew.Dump("workflow", bb)

	// 	w.s.Create(context.Background(), ww)
	// }

	return nil
}

func (w *wf) Deregister(ctx context.Context, paths []string) error {
	return nil
}

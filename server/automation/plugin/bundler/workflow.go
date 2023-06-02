package bundler

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	automationEnvoy "github.com/cortezaproject/corteza/server/automation/envoy"
	"github.com/cortezaproject/corteza/server/automation/types"
	"github.com/cortezaproject/corteza/server/pkg/envoyx"
)

type (
	WfService interface {
		Create(context.Context, *types.Workflow) (*types.Workflow, error)
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

func New(service WfService) *wf {
	evsvc := envoyx.New()
	evsvc.AddDecoder(envoyx.DecodeTypeIO,
		automationEnvoy.YamlDecoder{},
	)

	return &wf{
		s: service,
		e: evsvc,
	}
}

// Validate checks if any of the found subfolders of a bundle contains an automation workflow yaml file
func (w *wf) Validate(b []string) (byte, []string) {
	ss := []string{}

	for _, bb := range b {
		if strings.Contains(strings.ToLower(bb), "automation_workflow") {
			ga, _ := filepath.Glob(bb + "/*.yml")
			ss = append(ss, ga...)

			if len(ga) == 0 {
				gb, _ := filepath.Glob(bb + "/*.yaml")
				ss = append(ss, gb...)
			}
		}
	}

	if len(ss) == 0 {
		return 0, ss
	}

	return 1, ss
}

func (w *wf) Type() byte {
	return 1
}

func (w *wf) Register(path []string) error {
	// todo - fix

	// spew.Dump("FFFFF", w.w)

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

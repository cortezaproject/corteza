package yaml

import (
	"fmt"
	"strconv"

	atypes "github.com/cortezaproject/corteza/server/automation/types"
	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/envoy/resource"
)

type (
	composePage struct {
		res *types.Page
		ts  *resource.Timestamps

		envoyConfig   *resource.EnvoyConfig
		encoderConfig *EncoderConfig

		children composePageSet
		blocks   composePageBlockSet

		refNamespace string
		relNs        *types.Namespace
		refModule    string
		relMod       *types.Module
		refParent    string
		relParent    *types.Page

		rbac rbacRuleSet
	}
	composePageSet []*composePage

	composePageBlock struct {
		res *types.PageBlock

		relWf    []*atypes.Workflow
		refWf    []string
		relMod   []*types.Module
		refMod   []string
		relChart []*types.Chart
		refChart []string

		cfg *EncoderConfig
	}
	composePageBlockSet []*composePageBlock

	composePageBlockStyle = types.PageBlockStyle
)

func (pg *composePage) matches(i string) bool {
	if pg.res.ID > 0 && strconv.FormatUint(pg.res.ID, 10) == i {
		return true
	} else if pg.res.Handle == i {
		return true
	} else if pg.res.Title == i {
		return true
	}

	return false
}

func (c composePageSet) addComposePage(ii string, p *composePage) error {
	var cpg *composePage
	for _, s := range c {
		if s.matches(ii) {
			cpg = s
			break
		}
	}

	if cpg == nil {
		return composePageErrNotFound(ii)
	}
	if cpg.children == nil {
		cpg.children = make(composePageSet, 0, 1)
	}

	cpg.children = append(cpg.children, p)
	return nil
}

func (nn composePageSet) configureEncoder(cfg *EncoderConfig) {
	for _, n := range nn {
		n.encoderConfig = cfg
	}
}

func (bb composePageBlockSet) configureEncoder(cfg *EncoderConfig) {
	for _, b := range bb {
		b.cfg = cfg
	}
}

func composePageErrNotFound(i string) error {
	return fmt.Errorf("page not found: %v", i)
}

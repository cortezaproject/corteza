package yaml

import (
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
)

type (
	mappingTpl struct {
		resource.MappingTpl `yaml:",inline"`
	}
	mappingTplSet []*mappingTpl

	composeRecordTpl struct {
		Source string `yaml:"from"`

		Key         []string
		Mapping     mappingTplSet
		Defaultable bool
	}

	composeModule struct {
		res    *types.Module
		fields composeModuleFieldSet
		ts     *resource.Timestamps

		envoyConfig   *resource.EnvoyConfig
		encoderConfig *EncoderConfig

		relNs        *types.Namespace
		refNamespace string
		rbac         rbacRuleSet

		recTpl *composeRecordTpl
	}
	composeModuleSet []*composeModule

	composeModuleField struct {
		res  *types.ModuleField
		ts   *resource.Timestamps
		cfg  *EncoderConfig
		expr composeModuleFieldExpr

		relMod *types.Module

		rbac rbacRuleSet
	}
	composeModuleFieldSet []*composeModuleField

	// aux struct for decoding module field expressions
	composeModuleFieldExpr types.ModuleFieldExpr
)

func (nn composeModuleSet) ConfigureEncoder(cfg *EncoderConfig) {
	for _, n := range nn {
		n.encoderConfig = cfg
	}
}

func (ff composeModuleFieldSet) ConfigureEncoder(cfg *EncoderConfig) {
	for _, f := range ff {
		f.cfg = cfg
	}
}

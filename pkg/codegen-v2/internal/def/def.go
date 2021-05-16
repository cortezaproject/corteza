package def

import (
	"github.com/cortezaproject/corteza-server/pkg/codegen-v2/internal/tpl"
	"github.com/cortezaproject/corteza-server/pkg/y7s"
	"github.com/davecgh/go-spew/spew"
	"gopkg.in/yaml.v3"
	"strings"
)

var _ = spew.Dump

type (
	Document struct {
		Skip                bool `yaml:"(skip)"`
		Component           string
		IsComponentResource bool `yaml:"-"`
		Resource            string
		Source              string
		RBAC                *rbac
	}

	rbac struct {
		Schema     string
		Resource   *rbacResource
		Operations rbacOperations
	}

	rbacResource struct {
		Elements []string
	}

	rbacOperations []*rbacOperation

	rbacOperation struct {
		Operation   string
		CanFnName   string `yaml:"canFnName"`
		Description string
	}
)

func (set *rbacOperations) UnmarshalYAML(n *yaml.Node) error {
	return y7s.Each(n, func(k *yaml.Node, v *yaml.Node) (err error) {
		def := rbacOperation{}
		if k != nil {
			def.Operation = k.Value
		}

		*set = append(*set, &def)
		return v.Decode(&def)
	})
}

func (op *rbacOperation) UnmarshalYAML(n *yaml.Node) error {
	if y7s.IsKind(n, yaml.ScalarNode) {
		// @todo handle disabled operations
		// the idea is that when service operations are defined we implicitly define
		// RBAC operations. Here, we'll be able to remove implicitly defined operation
		return nil
	}

	type auxType rbacOperation
	var aux = (*auxType)(op)
	return n.Decode(aux)
}

func RbacOperationCanFnName(res, op string) string {
	// when check function name is not explicitly defined we try
	// to use resource and operation name and generate easy-to-read name
	//
	// <res> + <op>              => Can<Op><Res>
	// <res> + <op:foo.bar.verb> => Can<Verb><Foo><Bar>On<Res>

	if res == "Component" {
		res = ""
	}

	if strings.Contains(op, ".") {
		parts := strings.Split(op, ".")
		l := len(parts)

		parts = append(parts[l-1:], parts[:l-1]...)

		if res != "" {
			// Only append "on" if there is resource
			parts = append(parts, "on")
		}

		op = tpl.Export(parts...)
	}

	return tpl.Export("can", op, res)
}

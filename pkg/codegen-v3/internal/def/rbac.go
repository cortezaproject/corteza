package def

import (
	"fmt"
	"strings"

	"github.com/cortezaproject/corteza-server/pkg/codegen-v3/internal/tpl"
	"github.com/cortezaproject/corteza-server/pkg/y7s"
	"gopkg.in/yaml.v3"
)

type (
	rbac struct {
		// fully qualified resource name
		ResourceType string `yaml:"resourceType"`
		Resource     *rbacResource
		Operations   rbacOperations
	}

	rbacResource struct {
		References []*rbacResourceRef
		Attributes *rbacAttributes
	}

	rbacResourceRef struct {
		Field        string
		ResourceType string
		Resource     string
		Component    string

		custom bool
	}

	rbacOperations []*rbacOperation

	rbacOperation struct {
		Operation   string
		CanFnName   string `yaml:"canFnName"`
		Description string
	}

	rbacAttributes struct {
		Fields []string `yaml:"-"`
	}
)

func (r *rbac) proc(component, resource string) error {
	const (
		defaultNS   = "corteza"
		nsDelimiter = "::"
	)

	if r.ResourceType == "" {
		if strings.ToLower(resource) == "component" {
			r.ResourceType = component
		} else {
			r.ResourceType = fmt.Sprintf("%s:%s", component, resource)
		}

		r.ResourceType = defaultNS + nsDelimiter + r.ResourceType
	}

	if !strings.Contains(r.ResourceType, nsDelimiter) {
		return fmt.Errorf("no namespace prefix found (e.g.: 'corteza::') in resource type")
	}

	for _, op := range r.Operations {
		// Generate all check name
		if op.CanFnName == "" {
			op.CanFnName = RbacOperationCanFnName(resource, op.Operation)
		}
	}

	if r.Resource == nil {
		r.Resource = &rbacResource{References: []*rbacResourceRef{{Field: "ID"}}}
	}

	// check types of each referenced component
	// and prefix non-custom components with corteza::<component>
	// and self-references (field==ID) with own resource type
	for _, rc := range r.Resource.References {
		if !rc.custom {
			if rc.Field == "ID" {
				rc.ResourceType = r.ResourceType
				rc.Component = component
				rc.Resource = resource
			} else {
				rc.ResourceType = defaultNS + nsDelimiter + fmt.Sprintf("%s:%s", component, rc.Field)
				rc.Component = component
				rc.Resource = rc.Field
				rc.Field = rc.Field + "ID"
			}
		}
	}

	return nil
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

func (op *rbacResourceRef) UnmarshalYAML(n *yaml.Node) error {
	if y7s.IsKind(n, yaml.ScalarNode) {
		op.Field = n.Value
		if n.Value != "ID" {
			op.ResourceType = n.Value
			// @todo expand resource & component
		}

		return nil
	}

	type auxType rbacResourceRef
	var aux = (*auxType)(op)
	aux.custom = true
	return n.Decode(aux)
}

func (a *rbacAttributes) UnmarshalYAML(n *yaml.Node) error {
	if y7s.IsKind(n, yaml.ScalarNode) {
		return nil
	}

	// if not scalar, assume we will get list of fields
	a.Fields = make([]string, 0)
	return n.Decode(&a.Fields)
}

func RbacOperationCanFnName(res, op string) string {
	// when check function name is not explicitly defined we try
	// to use resource and operation name and generate easy-to-read name
	//
	// <res> + <op>              => Can<Op><Res>
	// <res> + <op:foo.bar.verb> => Can<Verb><Foo><Bar>On<Res>

	if strings.ToLower(res) == "component" {
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

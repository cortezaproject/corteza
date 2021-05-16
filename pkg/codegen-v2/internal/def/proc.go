package def

import (
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/codegen-v2/internal/tpl"
	"strings"
)

// Preproc preprocesses the document and sets defaults
func (doc *Document) Proc(filename string) error {
	doc.Source = filename

	// filename parts
	fp := strings.Split(filename, ".")
	// trim extension
	fp = fp[:len(fp)-1]
	if len(fp) > 0 {
		// set component from the 1st part
		doc.Component = fp[0]
	}

	if len(fp) > 1 {
		// if there are more parts, set resource
		doc.Resource = fp[1]
	}

	if strings.ToLower(doc.Resource) == "component" {
		return fmt.Errorf("can not use 'component' as a resource name")
	} else if doc.Resource == "" {
		doc.Resource = "component"
		doc.IsComponentResource = true
	}

	doc.Resource = tpl.Export(doc.Resource)

	if err := doc.RBAC.proc(doc.Resource, fmt.Sprintf("corteza+%s", strings.Join(fp, "."))); err != nil {
		return err
	}

	return nil
}

func (r *rbac) proc(resource, schema string) error {
	if r.Schema == "" {
		r.Schema = schema
	}

	for _, op := range r.Operations {
		// Generate all check name
		if op.CanFnName == "" {
			op.CanFnName = RbacOperationCanFnName(resource, op.Operation)
		}
	}

	if r.Resource == nil {
		r.Resource = &rbacResource{Elements: []string{"ID"}}
	}

	return nil
}

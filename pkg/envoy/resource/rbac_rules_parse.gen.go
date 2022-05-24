package resource

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"fmt"
	automationTypes "github.com/cortezaproject/corteza-server/automation/types"
	composeTypes "github.com/cortezaproject/corteza-server/compose/types"
	federationTypes "github.com/cortezaproject/corteza-server/federation/types"
	systemTypes "github.com/cortezaproject/corteza-server/system/types"
	"strings"
)

// Parse generates resource setting logic for each resource
//
// Resources with "envoy: false" are skipped
//
// This function is auto-generated
func ParseRule(res string) (string, *Ref, []*Ref, error) {
	if res == "" {
		return "", nil, nil, fmt.Errorf("empty resource")
	}

	sp := "/"

	res = strings.TrimSpace(res)
	res = strings.TrimRight(res, sp)
	rr := strings.Split(res, sp)

	// only service defined (corteza::system, corteza::compose, ...)
	if len(rr) == 1 {
		return res, nil, nil, nil
	}

	// full thing
	resourceType, path := rr[0], rr[1:]

	for p := 1; p < len(path); p++ {
		if path[p] != "*" && path[p-1] == "*" {
			return "", nil, nil, fmt.Errorf("invalid path wildcard combination for '%s'", res)
		}
	}

	// make the resource provide the slice of parent resources we should nest under
	switch resourceType {
	case systemTypes.ApplicationResourceType:
		if len(path) != 1 {
			return "", nil, nil, fmt.Errorf("expecting 1 reference components in path, got %d", len(path))
		}
		ref, pp, err := SystemApplicationRbacReferences(
			path[0],
		)
		return resourceType, ref, pp, err

	case systemTypes.ApigwRouteResourceType:
		if len(path) != 1 {
			return "", nil, nil, fmt.Errorf("expecting 1 reference components in path, got %d", len(path))
		}
		ref, pp, err := SystemApigwRouteRbacReferences(
			path[0],
		)
		return resourceType, ref, pp, err

	case systemTypes.AuthClientResourceType:
		if len(path) != 1 {
			return "", nil, nil, fmt.Errorf("expecting 1 reference components in path, got %d", len(path))
		}
		ref, pp, err := SystemAuthClientRbacReferences(
			path[0],
		)
		return resourceType, ref, pp, err

	case systemTypes.QueueResourceType:
		if len(path) != 1 {
			return "", nil, nil, fmt.Errorf("expecting 1 reference components in path, got %d", len(path))
		}
		ref, pp, err := SystemQueueRbacReferences(
			path[0],
		)
		return resourceType, ref, pp, err

	case systemTypes.QueueMessageResourceType:
		if len(path) != 1 {
			return "", nil, nil, fmt.Errorf("expecting 1 reference components in path, got %d", len(path))
		}
		ref, pp, err := SystemQueueMessageRbacReferences(
			path[0],
		)
		return resourceType, ref, pp, err

	case systemTypes.ReportResourceType:
		if len(path) != 1 {
			return "", nil, nil, fmt.Errorf("expecting 1 reference components in path, got %d", len(path))
		}
		ref, pp, err := SystemReportRbacReferences(
			path[0],
		)
		return resourceType, ref, pp, err

	case systemTypes.RoleResourceType:
		if len(path) != 1 {
			return "", nil, nil, fmt.Errorf("expecting 1 reference components in path, got %d", len(path))
		}
		ref, pp, err := SystemRoleRbacReferences(
			path[0],
		)
		return resourceType, ref, pp, err

	case systemTypes.TemplateResourceType:
		if len(path) != 1 {
			return "", nil, nil, fmt.Errorf("expecting 1 reference components in path, got %d", len(path))
		}
		ref, pp, err := SystemTemplateRbacReferences(
			path[0],
		)
		return resourceType, ref, pp, err

	case systemTypes.UserResourceType:
		if len(path) != 1 {
			return "", nil, nil, fmt.Errorf("expecting 1 reference components in path, got %d", len(path))
		}
		ref, pp, err := SystemUserRbacReferences(
			path[0],
		)
		return resourceType, ref, pp, err

	case systemTypes.DalConnectionResourceType:
		if len(path) != 1 {
			return "", nil, nil, fmt.Errorf("expecting 1 reference components in path, got %d", len(path))
		}
		ref, pp, err := SystemDalConnectionRbacReferences(
			path[0],
		)
		return resourceType, ref, pp, err

	case composeTypes.ChartResourceType:
		if len(path) != 2 {
			return "", nil, nil, fmt.Errorf("expecting 2 reference components in path, got %d", len(path))
		}
		ref, pp, err := ComposeChartRbacReferences(
			path[0],
			path[1],
		)
		return resourceType, ref, pp, err

	case composeTypes.ModuleResourceType:
		if len(path) != 2 {
			return "", nil, nil, fmt.Errorf("expecting 2 reference components in path, got %d", len(path))
		}
		ref, pp, err := ComposeModuleRbacReferences(
			path[0],
			path[1],
		)
		return resourceType, ref, pp, err

	case composeTypes.ModuleFieldResourceType:
		if len(path) != 3 {
			return "", nil, nil, fmt.Errorf("expecting 3 reference components in path, got %d", len(path))
		}
		ref, pp, err := ComposeModuleFieldRbacReferences(
			path[0],
			path[1],
			path[2],
		)
		return resourceType, ref, pp, err

	case composeTypes.NamespaceResourceType:
		if len(path) != 1 {
			return "", nil, nil, fmt.Errorf("expecting 1 reference components in path, got %d", len(path))
		}
		ref, pp, err := ComposeNamespaceRbacReferences(
			path[0],
		)
		return resourceType, ref, pp, err

	case composeTypes.PageResourceType:
		if len(path) != 2 {
			return "", nil, nil, fmt.Errorf("expecting 2 reference components in path, got %d", len(path))
		}
		ref, pp, err := ComposePageRbacReferences(
			path[0],
			path[1],
		)
		return resourceType, ref, pp, err

	case composeTypes.RecordResourceType:
		if len(path) != 3 {
			return "", nil, nil, fmt.Errorf("expecting 3 reference components in path, got %d", len(path))
		}
		ref, pp, err := ComposeRecordRbacReferences(
			path[0],
			path[1],
			path[2],
		)
		return resourceType, ref, pp, err

	case automationTypes.WorkflowResourceType:
		if len(path) != 1 {
			return "", nil, nil, fmt.Errorf("expecting 1 reference components in path, got %d", len(path))
		}
		ref, pp, err := AutomationWorkflowRbacReferences(
			path[0],
		)
		return resourceType, ref, pp, err

	case federationTypes.NodeResourceType:
		if len(path) != 1 {
			return "", nil, nil, fmt.Errorf("expecting 1 reference components in path, got %d", len(path))
		}
		ref, pp, err := FederationNodeRbacReferences(
			path[0],
		)
		return resourceType, ref, pp, err

	case federationTypes.ExposedModuleResourceType:
		if len(path) != 2 {
			return "", nil, nil, fmt.Errorf("expecting 2 reference components in path, got %d", len(path))
		}
		ref, pp, err := FederationExposedModuleRbacReferences(
			path[0],
			path[1],
		)
		return resourceType, ref, pp, err

	case federationTypes.SharedModuleResourceType:
		if len(path) != 2 {
			return "", nil, nil, fmt.Errorf("expecting 2 reference components in path, got %d", len(path))
		}
		ref, pp, err := FederationSharedModuleRbacReferences(
			path[0],
			path[1],
		)
		return resourceType, ref, pp, err

	case systemTypes.ComponentResourceType:
		if len(path) != 0 {
			return "", nil, nil, fmt.Errorf("expecting 0 reference components in path, got %d", len(path))
		}

		// Component resource, no path
		return resourceType, nil, nil, nil
	case composeTypes.ComponentResourceType:
		if len(path) != 0 {
			return "", nil, nil, fmt.Errorf("expecting 0 reference components in path, got %d", len(path))
		}

		// Component resource, no path
		return resourceType, nil, nil, nil
	case automationTypes.ComponentResourceType:
		if len(path) != 0 {
			return "", nil, nil, fmt.Errorf("expecting 0 reference components in path, got %d", len(path))
		}

		// Component resource, no path
		return resourceType, nil, nil, nil
	case federationTypes.ComponentResourceType:
		if len(path) != 0 {
			return "", nil, nil, fmt.Errorf("expecting 0 reference components in path, got %d", len(path))
		}

		// Component resource, no path
		return resourceType, nil, nil, nil
	}

	// return unhandled resource as-is
	return resourceType, nil, nil, nil
}

package request

import (
	"net/http"

	"github.com/crusttech/human/server/pkg/payload"
)

type (
	SearchResources struct {
		Q             string
		From          int
		Size          int
		Rt            []string
		NamespaceAggs []string
		ModuleAggs    []string
		NamespaceIDs  []string
		ModuleIDs     []string
		SearchMode    string
		DumpRaw       string
	}
)

// NewSearchListResources request
func NewSearchListResources() *SearchResources {
	return &SearchResources{}
}

// Auditable returns all auditable/loggable parameters
func (r SearchResources) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"q":             r.Q,
		"from":          r.From,
		"size":          r.Size,
		"namespaceAggs": r.NamespaceAggs,
		"resourceTypes": r.Rt,
		"moduleAggs":    r.ModuleAggs,
		"namespaceIDs":  r.NamespaceIDs,
		"moduleIDs":     r.ModuleIDs,
		"searchMode":    r.SearchMode,
		"dump":          r.DumpRaw,
	}
}

func (r SearchResources) GetQuery() string {
	return r.Q
}

func (r SearchResources) GetSize() int {
	return r.Size
}

func (r SearchResources) GetFrom() int {
	return r.From
}

func (r SearchResources) GetResourceTypes() []string {
	return r.Rt
}

func (r SearchResources) GetNamespaceAggs() []string {
	return r.NamespaceAggs
}

func (r SearchResources) GetModuleAggs() []string {
	return r.ModuleAggs
}

func (r SearchResources) GetNamespaceIDs() []string {
	return r.NamespaceIDs
}

func (r SearchResources) GetModuleIDs() []string {
	return r.ModuleIDs
}

func (r SearchResources) GetSearchMode() string {
	return r.SearchMode
}

func (r SearchResources) GetDumpRaw() string {
	return r.DumpRaw
}

// Fill processes request and fills internal variables
func (r *SearchResources) Fill(req *http.Request) (err error) {
	{
		// GET params
		tmp := req.URL.Query()

		if val, ok := tmp["q"]; ok && len(val) > 0 {
			r.Q, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := tmp["size"]; ok && len(val) > 0 {
			r.Size, err = payload.ParseInt(val[0]), nil
			if err != nil {
				return err
			}
		}

		if val, ok := tmp["from"]; ok && len(val) > 0 {
			r.From, err = payload.ParseInt(val[0]), nil
			if err != nil {
				return err
			}
		}

		if val, ok := tmp["resourceTypes[]"]; ok {
			r.Rt, err = val, nil
			if err != nil {
				return err
			}
		} else if val, ok = tmp["resourceTypes"]; ok {
			r.Rt, err = val, nil
			if err != nil {
				return err
			}
		}

		if val, ok := tmp["namespaceAggs[]"]; ok {
			r.NamespaceAggs, err = val, nil
			if err != nil {
				return err
			}
		} else if val, ok = tmp["namespaceAggs"]; ok {
			r.NamespaceAggs, err = val, nil
			if err != nil {
				return err
			}
		}

		if val, ok := tmp["moduleAggs[]"]; ok {
			r.ModuleAggs, err = val, nil
			if err != nil {
				return err
			}
		} else if val, ok = tmp["moduleAggs"]; ok {
			r.ModuleAggs, err = val, nil
			if err != nil {
				return err
			}
		}

		if val, ok := tmp["namespaceID[]"]; ok {
			r.NamespaceIDs, err = val, nil
			if err != nil {
				return err
			}
		} else if val, ok = tmp["namespaceID"]; ok {
			r.NamespaceIDs, err = val, nil
			if err != nil {
				return err
			}
		}

		if val, ok := tmp["moduleID[]"]; ok {
			r.ModuleIDs, err = val, nil
			if err != nil {
				return err
			}
		} else if val, ok = tmp["moduleID"]; ok {
			r.ModuleIDs, err = val, nil
			if err != nil {
				return err
			}
		}

		if val, ok := tmp["searchMode"]; ok && len(val) > 0 {
			r.SearchMode, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := tmp["dump"]; ok && len(val) > 0 {
			r.DumpRaw, err = val[0], nil
			if err != nil {
				return err
			}
		}
	}

	return err
}

package scim

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"
)

const (
	urnPatchOp    = "urn:ietf:params:scim:schemas:core:2.0:PatchOp"
	patchOpAdd    = "add"
	patchOpRemove = "remove"
)

type (
	// very crud operations support
	operationsRequest struct {
		Schemas    []string           `json:"schemas"`
		Operations []operationRequest `json:"Operations"`
	}

	operationRequest struct {
		Operation string `json:"op"`
		Path      string `json:"path"`
		Value     map[string]string
	}
)

func (req *operationsRequest) decodeJSON(r io.Reader) error {
	if err := json.NewDecoder(r).Decode(req); err != nil {
		return fmt.Errorf("could not decode operations payload: %w", err)
	}

	// Compiles req.Path and updates req.Value
	for i, op := range req.Operations {
		r := regexp.MustCompile(`([a-z]\w+)\[value eq \"([\w]{8}-[\w]{4}-[\w]{4}-[\w]{4}-[\w]{12})\"\]`)
		strings := r.FindStringSubmatch(op.Path)
		if len(strings) == 3 {
			key := strings[1]
			if op.Value == nil {
				op.Value = make(map[string]string)
			}
			op.Value[key] = strings[2]
			req.Operations[i].Value = op.Value
		}
	}

	return nil
}

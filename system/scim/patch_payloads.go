package scim

import (
	"encoding/json"
	"fmt"
	"io"
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
		Value     []struct {
			Value string `json:"value"`
		} `json:"value"`
	}
)

func (req *operationsRequest) decodeJSON(r io.Reader) error {
	if err := json.NewDecoder(r).Decode(req); err != nil {
		return fmt.Errorf("could not decode operations payload: %w", err)
	}

	return nil
}

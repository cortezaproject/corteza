package scim

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestOperationsRequestDecodeJSON(t *testing.T) {
	const user1Id = `00000000-0000-0000-0000-000000000001`
	var (
		req     = require.New(t)
		payload operationsRequest

		s = strings.NewReader(fmt.Sprintf(
			`{"Operations":[{"op":"add","path":"members[value eq \"%s\"]"}],"schemas":["urn:ietf:params:scim:schemas:core:2.0:PatchOp"]}`,
			user1Id,
		))

		expectedPayload = operationsRequest{
			Schemas: []string{"urn:ietf:params:scim:schemas:core:2.0:PatchOp"},
			Operations: []operationRequest{
				{
					Operation: "add",
					Path:      fmt.Sprintf("members[value eq \"%s\"]", user1Id),
					Value: map[string]string{
						"members": user1Id,
					},
				},
			},
		}
	)

	err := payload.decodeJSON(s)
	req.NoError(err)
	req.Equal(expectedPayload, payload)
}

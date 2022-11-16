package federation

import (
	"strings"
	"testing"
	"time"

	ct "github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/federation/types"
	"github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/cortezaproject/corteza/server/pkg/options"
	"github.com/stretchr/testify/require"
)

func TestEncoder_encodeStructure(t *testing.T) {
	var (
		opts = options.FederationOpt{
			Enabled: true,
			Host:    "example.ltd",
		}

		payload = ListStructurePayload{
			Filter: &types.ExposedModuleFilter{
				Paging: filter.Paging{
					Limit: 42,
				},
			},
			Set: &types.ExposedModuleSet{
				&types.ExposedModule{
					ID:                 123,
					ComposeModuleID:    456,
					ComposeNamespaceID: 789,
					NodeID:             1111,
					Handle:             "test_module",
					Name:               "Test Module",
					CreatedBy:          2222,
					CreatedAt:          time.Date(2020, 10, 10, 10, 10, 10, 10, time.UTC),
					Fields: types.ModuleFieldSet{
						&types.ModuleField{Kind: "String", Name: "test_module_first_field", Label: "Test Module First Field", IsMulti: false},
						&types.ModuleField{Kind: "String", Name: "test_module_second_field", Label: "Test Module Second Field", IsMulti: false},
					},
				},
			},
		}

		tcc = []struct {
			name   string
			format EncodingFormat
			expect string
		}{
			{
				"Encode structure in Activity Streams format",
				ActivityStreamsStructure,
				`{"@context":"https://www.w3.org/ns/activitystreams","itemsPerPage":42,"items":[{"@context":"https://www.w3.org/ns/activitystreams","type":"Module","summary":"Structure for module Test Module on node 1111","name":"Test Module","handle":"test_module","url":"https://example.ltd/nodes/1111/modules/123","node":"1111","federationModule":"123","composeModule":"456","composeNamespace":"789","createdAt":"2020-10-10T10:10:10.00000001Z","createdBy":"2222","attributedTo":[{"@context":"https://www.w3.org/ns/activitystreams","id":"https://example.ltd/system/users/2222","type":"User"}],"fields":[{"kind":"String","name":"test_module_first_field","label":"Test Module First Field","isMulti":false},{"kind":"String","name":"test_module_second_field","label":"Test Module Second Field","isMulti":false}]}]}`,
			},
			{
				"Encode structure in internal Corteza format",
				CortezaInternalStructure,
				`{"response":{"filter":{"nodeID":"0","composeModuleID":"0","composeNamespaceID":"0","lastSync":0,"handle":"","name":"","query":"","limit":42},"set":[{"moduleID":"123","nodeID":"1111","composeModuleID":"456","composeNamespaceID":"789","handle":"test_module","name":"Test Module","fields":[{"kind":"String","name":"test_module_first_field","label":"Test Module First Field","isMulti":false},{"kind":"String","name":"test_module_second_field","label":"Test Module Second Field","isMulti":false}],"createdAt":"2020-10-10T10:10:10.00000001Z","createdBy":"2222"}]}}`,
			},
		}
	)

	for _, tc := range tcc {
		t.Run(tc.name, func(t *testing.T) {
			var (
				req    = require.New(t)
				writer = strings.Builder{}
			)

			encoder := NewEncoder(&writer, opts)
			err := encoder.Encode(payload, tc.format)

			req.NoError(err)
			req.Equal(tc.expect, strings.TrimSuffix(writer.String(), "\n"))
		})
	}
}

func TestEncoder_encodeData(t *testing.T) {
	var (
		opts = options.FederationOpt{
			Enabled: true,
			Host:    "example.ltd",
		}

		payload = ListDataPayload{
			Filter: &ct.RecordFilter{
				ModuleID:    456,
				NamespaceID: 789,
				Paging: filter.Paging{
					Limit: 42,
				},
			},
			Set: &ct.RecordSet{
				&ct.Record{
					ID:          123,
					ModuleID:    456,
					NamespaceID: 789,
					CreatedBy:   2222,
					CreatedAt:   time.Date(2020, 10, 10, 10, 10, 10, 10, time.UTC),
					Values: ct.RecordValueSet{
						&ct.RecordValue{RecordID: 1231111, Name: "First Record First Value", Value: "First Record First Value"},
						&ct.RecordValue{RecordID: 1231112, Name: "First Record Second Value", Value: "First Record Second Value"},
					},
				},
				&ct.Record{
					ID:          124,
					ModuleID:    456,
					NamespaceID: 789,
					CreatedBy:   2222,
					CreatedAt:   time.Date(2020, 10, 10, 10, 10, 10, 10, time.UTC),
					Values: ct.RecordValueSet{
						&ct.RecordValue{RecordID: 1231111, Name: "Second Record First Value", Value: "Second Record First Value"},
						&ct.RecordValue{RecordID: 1231112, Name: "Second Record Second Value", Value: "Second Record Second Value"},
					},
				},
			},
			NodeID:   1111,
			ModuleID: 3333,
		}

		tcc = []struct {
			name   string
			format EncodingFormat
			expect string
		}{
			{
				"Encode structure in Activity Streams format",
				ActivityStreamsData,
				`{"@context":"https://www.w3.org/ns/activitystreams","itemsPerPage":42,"items":[{"@context":"https://www.w3.org/ns/activitystreams","type":"Record","summary":"Data for module 456 on node 1111","url":"https://example.ltd/nodes/1111/modules/3333/records/social/","node":"1111","federationModule":"456","composeModule":"456","composeNamespace":"789","createdAt":"2020-10-10T10:10:10.00000001Z","createdBy":"2222","attributedTo":[{"@context":"https://www.w3.org/ns/activitystreams","id":"https://example.ltd/system/users/2222","type":"User"}],"values":[{"name":"First Record First Value","value":"First Record First Value"},{"name":"First Record Second Value","value":"First Record Second Value"}]},{"@context":"https://www.w3.org/ns/activitystreams","type":"Record","summary":"Data for module 456 on node 1111","url":"https://example.ltd/nodes/1111/modules/3333/records/social/","node":"1111","federationModule":"456","composeModule":"456","composeNamespace":"789","createdAt":"2020-10-10T10:10:10.00000001Z","createdBy":"2222","attributedTo":[{"@context":"https://www.w3.org/ns/activitystreams","id":"https://example.ltd/system/users/2222","type":"User"}],"values":[{"name":"Second Record First Value","value":"Second Record First Value"},{"name":"Second Record Second Value","value":"Second Record Second Value"}]}]}`,
			},
			{
				"Encode structure in internal Corteza format",
				CortezaInternalData,
				`{"response":{"filter":{"moduleID":"456","namespaceID":"789","query":"","deleted":0,"limit":42},"set":[{"recordID":"123","moduleID":"456","values":[{"name":"First Record First Value","value":"First Record First Value"},{"name":"First Record Second Value","value":"First Record Second Value"}],"namespaceID":"789","ownedBy":"0","createdAt":"2020-10-10T10:10:10.00000001Z","createdBy":"2222"},{"recordID":"124","moduleID":"456","values":[{"name":"Second Record First Value","value":"Second Record First Value"},{"name":"Second Record Second Value","value":"Second Record Second Value"}],"namespaceID":"789","ownedBy":"0","createdAt":"2020-10-10T10:10:10.00000001Z","createdBy":"2222"}]}}`,
			},
		}
	)

	for _, tc := range tcc {
		t.Run(tc.name, func(t *testing.T) {
			var (
				req    = require.New(t)
				writer = strings.Builder{}
			)

			encoder := NewEncoder(&writer, opts)
			err := encoder.Encode(payload, tc.format)

			req.NoError(err)
			req.Equal(tc.expect, strings.TrimSuffix(writer.String(), "\n"))
		})
	}
}

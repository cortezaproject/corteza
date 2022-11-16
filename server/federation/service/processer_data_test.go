package service

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	cs "github.com/cortezaproject/corteza/server/compose/service"
	ct "github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/federation/types"
	ss "github.com/cortezaproject/corteza/server/system/service"
	st "github.com/cortezaproject/corteza/server/system/types"
	"github.com/stretchr/testify/require"
)

type (
	testSyncService struct {
		Sync
	}

	testSharedModuleService struct {
		SharedModuleService
	}
	testRecordServicePersistSuccess struct {
		cs.RecordService
	}
	testRecordServiceUpdateSuccess struct {
		cs.RecordService
	}
	testRecordServiceDeleteSuccess struct {
		cs.RecordService
	}
	testRecordServicePersistError struct {
		cs.RecordService
	}
	testUserService struct {
		ss.UserService
	}
	testRoleService struct {
		ss.RoleService
	}
)

func TestProcesserData_persist(t *testing.T) {

	var (
		tcc = []struct {
			name      string
			payload   string
			mappings  string
			persisted int
			err       string
			values    *ct.RecordValueSet
			s         *Sync
		}{
			{
				"successful persist on valid mapping",
				`{"response": {"set": [{"recordID":"1","values":[{"name":"Facebook","value":"foobar"}],"createdAt":"2020-12-05T10:10:10Z"}]}}`,
				`[{"origin":{"kind":"String","name":"Description","label":"Description","isMulti":false},"destination":{"kind":"String","name":"Name","label":"Description","isMulti":false}},{"origin":{"kind":"Url","name":"Facebook","label":"Facebook","isMulti":false},"destination":{"kind":"Url","name":"Fb","label":"Facebook","isMulti":false}}]`,
				1,
				"",
				&ct.RecordValueSet{&ct.RecordValue{Name: "Fb", Value: ""}},
				NewSync(
					&Syncer{},
					&Mapper{},
					&testSharedModuleService{},
					&testRecordServicePersistSuccess{},
					&testUserService{},
					&testRoleService{}),
			},
			{
				"successful update on valid mapping and existing federated record",
				`{"response": {"set": [{"recordID":"1","values":[{"name":"Facebook","value":"foobar"}],"createdAt":"2020-12-05T10:10:10Z", "updatedAt":"2020-12-06T10:10:10Z"}]}}`,
				`[{"origin":{"kind":"String","name":"Description","label":"Description","isMulti":false},"destination":{"kind":"String","name":"Name","label":"Description","isMulti":false}},{"origin":{"kind":"Url","name":"Facebook","label":"Facebook","isMulti":false},"destination":{"kind":"Url","name":"Fb","label":"Facebook","isMulti":false}}]`,
				1,
				"",
				&ct.RecordValueSet{&ct.RecordValue{Name: "Fb", Value: ""}},
				NewSync(
					&Syncer{},
					&Mapper{},
					&testSharedModuleService{},
					&testRecordServiceUpdateSuccess{},
					&testUserService{},
					&testRoleService{}),
			},
			{
				"successful delete on valid mapping and existing federated record",
				`{"response": {"set": [{"recordID":"1","values":[{"name":"Facebook","value":"foobar"}],"createdAt":"2020-12-05T10:10:10Z", "deletedAt":"2020-12-07T10:10:10Z"}]}}`,
				`[{"origin":{"kind":"String","name":"Description","label":"Description","isMulti":false},"destination":{"kind":"String","name":"Name","label":"Description","isMulti":false}},{"origin":{"kind":"Url","name":"Facebook","label":"Facebook","isMulti":false},"destination":{"kind":"Url","name":"Fb","label":"Facebook","isMulti":false}}]`,
				1,
				"",
				&ct.RecordValueSet{&ct.RecordValue{Name: "Fb", Value: ""}},
				NewSync(
					&Syncer{},
					&Mapper{},
					&testSharedModuleService{},
					&testRecordServiceDeleteSuccess{},
					&testUserService{},
					&testRoleService{}),
			},
			{
				"persist error on valid mapping",
				`{"response": {"set": [{"recordID":"1","values":[{"name":"Facebook","value":"foobar"}]}]}}`,
				`[{"origin":{"kind":"String","name":"Description","label":"Description","isMulti":false},"destination":{"kind":"String","name":"Name","label":"Description","isMulti":false}},{"origin":{"kind":"Url","name":"Facebook","label":"Facebook","isMulti":false},"destination":{"kind":"Url","name":"Fb","label":"Facebook","isMulti":false}}]`,
				0,
				"",
				&ct.RecordValueSet{&ct.RecordValue{Name: "Fb", Value: ""}},
				NewSync(
					&Syncer{},
					&Mapper{},
					&testSharedModuleService{},
					&testRecordServicePersistError{},
					&testUserService{},
					&testRoleService{}),
			},
			{
				"persist one record on mixed mappings",
				`{"response": {"set": [{"recordID":"1","values":[{"name":"THIS_FIELD_NAME_IS_NOT_ON_ORIGIN","value":"this value will not be set"}]}, {"recordID":"1","values":[{"name":"Facebook","value":"this value WILL be set"}]}]}}`,
				`[{"origin":{"kind":"String","name":"Description","label":"Description","isMulti":false},"destination":{"kind":"String","name":"Name","label":"Description","isMulti":false}},{"origin":{"kind":"Url","name":"Facebook","label":"Facebook","isMulti":false},"destination":{"kind":"Url","name":"Fb","label":"Facebook","isMulti":false}}]`,
				2,
				"",
				&ct.RecordValueSet{&ct.RecordValue{Name: "Fb", Value: ""}},
				NewSync(
					&Syncer{},
					&Mapper{},
					&testSharedModuleService{},
					&testRecordServicePersistSuccess{},
					&testUserService{},
					&testRoleService{}),
			},
			{
				"no records from payload",
				`{"response": {"set": []}}`,
				`[{"origin":{"kind":"String","name":"Description","label":"Description","isMulti":false},"destination":{"kind":"String","name":"Name","label":"Description","isMulti":false}},{"origin":{"kind":"Url","name":"Facebook","label":"Facebook","isMulti":false},"destination":{"kind":"Url","name":"Fb","label":"Facebook","isMulti":false}}]`,
				0,
				"",
				&ct.RecordValueSet{&ct.RecordValue{Name: "Fb", Value: ""}},
				NewSync(
					&Syncer{},
					&Mapper{},
					&testSharedModuleService{},
					&testRecordServicePersistSuccess{},
					&testUserService{},
					&testRoleService{}),
			},
			{
				"validation error, no persist",
				`{"respon`,
				`[{"origin":{"kind":"String","name":"Description","label":"Description","isMulti":false},"destination":{"kind":"String","name":"Name","label":"Description","isMulti":false}},{"origin":{"kind":"Url","name":"Facebook","label":"Facebook","isMulti":false},"destination":{"kind":"Url","name":"Fb","label":"Facebook","isMulti":false}}]`,
				0,
				"unexpected end of JSON input",
				&ct.RecordValueSet{&ct.RecordValue{Name: "Fb", Value: ""}},
				NewSync(
					&Syncer{},
					&Mapper{},
					&testSharedModuleService{},
					&testRecordServicePersistSuccess{},
					&testUserService{},
					&testRoleService{}),
			},
		}
	)

	for _, tc := range tcc {
		t.Run(tc.name, func(t *testing.T) {
			var (
				ctx = context.Background()
				req = require.New(t)
			)

			mm := &types.ModuleFieldMappingSet{}
			json.Unmarshal([]byte(tc.mappings), mm)

			dp := &dataProcesser{
				ID:                  1,
				ComposeModuleID:     1,
				ComposeNamespaceID:  1,
				NodeBaseURL:         "",
				ModuleMappings:      mm,
				ModuleMappingValues: tc.values,
				SyncService:         tc.s,
				Node:                &types.Node{},
				User:                &st.User{},
			}

			out, err := dp.Process(ctx, []byte(tc.payload))

			if tc.err != "" {
				req.Equal(tc.err, err.Error())
			} else {
				req.NoError(err)
			}

			req.Equal(tc.persisted, out.(dataProcesserResponse).Processed)
		})
	}
}

// create success
func (s testRecordServicePersistSuccess) Create(_ context.Context, record *ct.Record) (*ct.Record, error) {
	return nil, nil
}

func (s testRecordServicePersistSuccess) Find(_ context.Context, filter ct.RecordFilter) (ct.RecordSet, ct.RecordFilter, error) {
	return ct.RecordSet{}, ct.RecordFilter{}, nil
}

// update success
func (s testRecordServiceUpdateSuccess) Update(_ context.Context, record *ct.Record) (*ct.Record, error) {
	return nil, nil
}

func (s testRecordServiceUpdateSuccess) Find(_ context.Context, filter ct.RecordFilter) (ct.RecordSet, ct.RecordFilter, error) {
	return ct.RecordSet{&ct.Record{ID: 2}}, ct.RecordFilter{}, nil
}

// delete success
func (s testRecordServiceDeleteSuccess) DeleteByID(_ context.Context, namespaceID, moduleID uint64, recordID ...uint64) error {
	return nil
}

func (s testRecordServiceDeleteSuccess) Find(_ context.Context, filter ct.RecordFilter) (ct.RecordSet, ct.RecordFilter, error) {
	return ct.RecordSet{&ct.Record{ID: 2, ModuleID: 2, NamespaceID: 2}}, ct.RecordFilter{}, nil
}

// create error
func (s testRecordServicePersistError) Create(_ context.Context, record *ct.Record) (*ct.Record, error) {
	return nil, errors.New("mocked error")
}

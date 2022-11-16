package service

import (
	"context"
	"errors"
	"testing"

	cs "github.com/cortezaproject/corteza/server/compose/service"
	"github.com/cortezaproject/corteza/server/federation/types"
	ss "github.com/cortezaproject/corteza/server/system/service"
	st "github.com/cortezaproject/corteza/server/system/types"
	"github.com/stretchr/testify/require"
)

type (
	testSyncStructureService                     struct{ Sync }
	testSyncStructureServiceCanUpdateFieldsError struct{ Sync }

	testStructureSharedModuleService                   struct{ SharedModuleService }
	testStructureSharedModuleServiceError              struct{ SharedModuleService }
	testStructureSharedModuleServiceCreateNewModule    struct{ SharedModuleService }
	testStructureSharedModuleServiceCreateNewModuleErr struct{ SharedModuleService }
	testStructureSharedModuleServiceUpdateModule       struct{ SharedModuleService }
	testStructureSharedModuleServiceUpdateModuleErr    struct{ SharedModuleService }

	testStructureRecordServicePersistSuccess struct{ cs.RecordService }
	testStructureUserService                 struct{ ss.UserService }
	testStructureRoleService                 struct{ ss.RoleService }
)

func TestProcesserStructure_persist(t *testing.T) {

	var (
		tcc = []struct {
			name      string
			payload   string
			persisted int
			err       string
			s         *Sync
		}{
			{
				"no modules from payload",
				`{"response": {"set": []}}`,
				0,
				"",
				NewSync(
					&Syncer{},
					&Mapper{},
					&testStructureSharedModuleService{},
					&testRecordServicePersistSuccess{},
					&testUserService{},
					&testRoleService{}),
			},
			{
				"validation error, no persist",
				`{"respon`,
				0,
				"unexpected end of JSON input",
				NewSync(
					&Syncer{},
					&Mapper{},
					&testStructureSharedModuleService{},
					&testRecordServicePersistSuccess{},
					&testUserService{},
					&testRoleService{}),
			},
			{
				"lookup existing module error",
				`{"response":{"set": [{"moduleID": "11", "nodeID": "21", "composeModuleID": "31", "composeNamespaceID":"32","handle":"handle","name":"name","fields":[]}]}}`,
				0,
				"db error",
				NewSync(
					&Syncer{},
					&Mapper{},
					&testStructureSharedModuleServiceError{},
					&testRecordServicePersistSuccess{},
					&testUserService{},
					&testRoleService{}),
			},
			{
				"create new module",
				`{"response":{"set": [{"moduleID": "11", "nodeID": "21", "composeModuleID": "31", "composeNamespaceID":"32","handle":"handle","name":"name","fields":[]}]}}`,
				1,
				"",
				NewSync(
					&Syncer{},
					&Mapper{},
					&testStructureSharedModuleServiceCreateNewModule{},
					&testRecordServicePersistSuccess{},
					&testUserService{},
					&testRoleService{}),
			},
			{
				"create new module error",
				`{"response":{"set": [{"moduleID": "11", "nodeID": "21", "composeModuleID": "31", "composeNamespaceID":"32","handle":"handle","name":"name","fields":[]}]}}`,
				0,
				"could not create new module",
				NewSync(
					&Syncer{},
					&Mapper{},
					&testStructureSharedModuleServiceCreateNewModuleErr{},
					&testRecordServicePersistSuccess{},
					&testUserService{},
					&testRoleService{}),
			},
			{
				"could not update module fields error",
				`{"response":{"set": [{"moduleID": "11", "nodeID": "21", "composeModuleID": "31", "composeNamespaceID":"32","handle":"handle","name":"name","fields":[]}]}}`,
				0,
				"module structure changed",
				NewSync(
					&Syncer{},
					&Mapper{},
					&testStructureSharedModuleServiceUpdateModule{},
					&testRecordServicePersistSuccess{},
					&testUserService{},
					&testRoleService{}),
			},
			{
				"update module fields success",
				`{"response":{"set": [{"moduleID": "11", "nodeID": "21", "composeModuleID": "31", "composeNamespaceID":"32","handle":"handle","name":"name","fields":[{"kind":"String","label":"label","isMulti":false,"name":"existingField"}]}]}}`,
				1,
				"",
				NewSync(
					&Syncer{},
					&Mapper{},
					&testStructureSharedModuleServiceUpdateModule{},
					&testRecordServicePersistSuccess{},
					&testUserService{},
					&testRoleService{}),
			},
			{
				"update module error",
				`{"response":{"set": [{"moduleID": "11", "nodeID": "21", "composeModuleID": "31", "composeNamespaceID":"32","handle":"handle","name":"name","fields":[]}]}}`,
				0,
				"could not update module",
				NewSync(
					&Syncer{},
					&Mapper{},
					&testStructureSharedModuleServiceUpdateModuleErr{},
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

			dp := &structureProcesser{
				SyncService: tc.s,
				Node:        &types.Node{},
				User:        &st.User{},
			}

			out, err := dp.Process(ctx, []byte(tc.payload))

			if tc.err != "" {
				if err != nil {
					req.Equal(tc.err, err.Error())
				} else {
					req.Fail("err should not be nil")
				}
			} else {
				req.NoError(err)
			}

			req.Equal(tc.persisted, out.(structureProcesserResponse).Processed)
		})
	}
}

func (s testStructureSharedModuleServiceError) Find(ctx context.Context, filter types.SharedModuleFilter) (types.SharedModuleSet, types.SharedModuleFilter, error) {
	return types.SharedModuleSet{}, types.SharedModuleFilter{}, errors.New("db error")
}

func (s testStructureSharedModuleServiceCreateNewModule) Find(ctx context.Context, filter types.SharedModuleFilter) (types.SharedModuleSet, types.SharedModuleFilter, error) {
	return types.SharedModuleSet{}, types.SharedModuleFilter{}, nil
}

func (s testStructureSharedModuleServiceCreateNewModule) Create(ctx context.Context, new *types.SharedModule) (*types.SharedModule, error) {
	return &types.SharedModule{}, nil
}

func (s testStructureSharedModuleServiceCreateNewModuleErr) Find(ctx context.Context, filter types.SharedModuleFilter) (types.SharedModuleSet, types.SharedModuleFilter, error) {
	return types.SharedModuleSet{}, types.SharedModuleFilter{}, nil
}

func (s testStructureSharedModuleServiceCreateNewModuleErr) Create(ctx context.Context, new *types.SharedModule) (*types.SharedModule, error) {
	return nil, errors.New("could not create new module")
}

func (s testStructureSharedModuleServiceUpdateModule) Find(ctx context.Context, filter types.SharedModuleFilter) (types.SharedModuleSet, types.SharedModuleFilter, error) {
	return types.SharedModuleSet{
			&types.SharedModule{
				ID:     11,
				NodeID: 21,
				Handle: "handle",
				Name:   "name",
				Fields: types.ModuleFieldSet{
					&types.ModuleField{
						Kind:    "String",
						Label:   "label",
						IsMulti: false,
						Name:    "existingField",
					},
				},
			}},
		types.SharedModuleFilter{},
		nil
}

func (s testStructureSharedModuleServiceUpdateModule) Update(ctx context.Context, updated *types.SharedModule) (*types.SharedModule, error) {
	return &types.SharedModule{}, nil
}

func (s testStructureSharedModuleServiceUpdateModuleErr) Find(ctx context.Context, filter types.SharedModuleFilter) (types.SharedModuleSet, types.SharedModuleFilter, error) {
	return types.SharedModuleSet{&types.SharedModule{
			Fields: types.ModuleFieldSet{},
		}},
		types.SharedModuleFilter{},
		nil
}

func (s testStructureSharedModuleServiceUpdateModuleErr) Update(ctx context.Context, updated *types.SharedModule) (*types.SharedModule, error) {
	return nil, errors.New("could not update module")
}

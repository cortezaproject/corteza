package seeder

import (
	"context"
	"fmt"
	cService "github.com/cortezaproject/corteza-server/compose/service"
	"time"

	cTypes "github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/id"
	lTypes "github.com/cortezaproject/corteza-server/pkg/label/types"
	"github.com/cortezaproject/corteza-server/store"
	sTypes "github.com/cortezaproject/corteza-server/system/types"
)

type (
	seeder struct {
		ctx   context.Context
		store storeService
		faker fakerService

		modSvc moduleService
	}
	Params struct {
		// (optional) no record to be generate; Default value will be 1
		Limit int
	}
	RecordParams struct {
		NamespaceID     uint64
		NamespaceHandle string
		ModuleID        uint64
		ModuleHandle    string
		Params
	}

	fakerService interface {
		fakeValueByName(name string) (val string, ok bool)
		fakeValue(name, kind string, opt valueOptions) (val string, err error)
		fakeUserHandle(s string) string
	}

	storeService interface {
		UpsertLabel(ctx context.Context, rr ...*lTypes.Label) error
		SearchUsers(ctx context.Context, f sTypes.UserFilter) (sTypes.UserSet, sTypes.UserFilter, error)
		CreateUser(ctx context.Context, rr ...*sTypes.User) error
		DeleteUser(ctx context.Context, rr ...*sTypes.User) error

		LookupComposeNamespaceBySlug(ctx context.Context, slug string) (*cTypes.Namespace, error)
		LookupComposeNamespaceByID(ctx context.Context, id uint64) (*cTypes.Namespace, error)
		LookupComposeModuleByNamespaceIDHandle(ctx context.Context, namespaceID uint64, name string) (*cTypes.Module, error)
		LookupComposeModuleByID(ctx context.Context, id uint64) (*cTypes.Module, error)
		SearchComposeModuleFields(ctx context.Context, f cTypes.ModuleFieldFilter) (cTypes.ModuleFieldSet, cTypes.ModuleFieldFilter, error)

		SearchComposeRecords(ctx context.Context, _mod *cTypes.Module, f cTypes.RecordFilter) (cTypes.RecordSet, cTypes.RecordFilter, error)
		CreateComposeRecord(ctx context.Context, mod *cTypes.Module, rr ...*cTypes.Record) error
		DeleteComposeRecord(ctx context.Context, m *cTypes.Module, rr ...*cTypes.Record) error
	}

	moduleService interface {
		FindByID(ctx context.Context, namespaceID uint64, moduleID uint64) (*cTypes.Module, error)
		FindByHandle(ctx context.Context, namespaceID uint64, handle string) (*cTypes.Module, error)
		Find(ctx context.Context, filter cTypes.ModuleFilter) (set cTypes.ModuleSet, f cTypes.ModuleFilter, err error)
	}
)

var (
	DefaultStore store.Storer
)

const (
	FakeDataLabel = "generated"
)

func Seeder(ctx context.Context, store store.Storer, faker fakerService) *seeder {
	DefaultStore = store

	return &seeder{
		ctx:    ctx,
		store:  store,
		faker:  faker,
		modSvc: cService.DefaultModule,
	}
}

// getLimit return data generation limit; It will return Default(1) if limit is 0
func (p Params) getLimit() int {
	if p.Limit == 0 {
		return 1
	}
	return p.Limit
}

// CreateLabel return the label for generate data
func (s seeder) CreateLabel(resourceID uint64, kind, name string) *lTypes.Label {
	return &lTypes.Label{
		Kind:       kind,
		ResourceID: resourceID,
		Name:       name,
		Value:      FakeDataLabel,
	}
}

// CreateUser creates given no of users into DB
func (s seeder) CreateUser(params Params) (IDs []uint64, err error) {
	var users []*sTypes.User
	var labels []*lTypes.Label

	for i := 0; i < params.getLimit(); i++ {
		var user sTypes.User
		user.ID = id.Next()
		user.Email, _ = s.faker.fakeValueByName("Email")
		user.Name, _ = s.faker.fakeValueByName("Name")
		user.Handle = s.faker.fakeUserHandle(user.Name)
		user.Kind = sTypes.NormalUser
		user.CreatedAt = time.Now()

		IDs = append(IDs, user.ID)
		users = append(users, &user)
		labels = append(labels, s.CreateLabel(
			user.ID,
			user.LabelResourceKind(),
			FakeDataLabel,
		))
	}

	err = s.store.CreateUser(s.ctx, users...)
	if err != nil {
		return
	}

	err = s.store.UpsertLabel(s.ctx, labels...)
	if err != nil {
		return
	}

	return
}

// DeleteAllUser deletes all the fake users from DB
func (s seeder) DeleteAllUser() (err error) {
	filter := sTypes.UserFilter{
		Labels: map[string]string{
			FakeDataLabel: FakeDataLabel,
		},
	}
	users, _, err := s.store.SearchUsers(s.ctx, filter)
	if err != nil {
		return
	}

	err = s.store.DeleteUser(s.ctx, users...)
	if err != nil {
		return
	}

	return
}

// LookupNamespaceByID will get namespace by ID
func (s seeder) LookupNamespaceByID(ID uint64) (ns *cTypes.Namespace, err error) {
	if ID == 0 {
		err = fmt.Errorf("invalid ID for namespace")
		return
	}
	ns, err = s.store.LookupComposeNamespaceByID(s.ctx, ID)
	if err != nil {
		return
	}
	if ns == nil {
		return ns, fmt.Errorf("namespace not found")
	}
	return
}

// LookupNamespaceByHandle will get namespace by handle
func (s seeder) LookupNamespaceByHandle(handle string) (ns *cTypes.Namespace, err error) {
	if len(handle) == 0 {
		err = fmt.Errorf("invalid handle for namespace")
		return
	}
	ns, err = s.store.LookupComposeNamespaceBySlug(s.ctx, handle)
	if err != nil {
		return
	}
	if ns == nil {
		return ns, fmt.Errorf("namespace not found")
	}
	return
}

// LookupModuleByID will get module by ID
func (s seeder) LookupModuleByID(ID uint64) (mod *cTypes.Module, err error) {
	if ID == 0 {
		err = fmt.Errorf("invalid ID for module")
		return
	}
	mod, err = s.store.LookupComposeModuleByID(s.ctx, ID)
	if err != nil {
		return
	}
	if mod == nil {
		return mod, fmt.Errorf("module not found")
	}

	mod.Fields, _, err = s.store.SearchComposeModuleFields(s.ctx, cTypes.ModuleFieldFilter{ModuleID: []uint64{mod.ID}})
	if err != nil {
		return
	}
	return
}

// LookupModuleByHandle will get module by handle
func (s seeder) LookupModuleByHandle(nsID uint64, handle string) (mod *cTypes.Module, err error) {
	if nsID == 0 {
		err = fmt.Errorf("invalid namespace for module")
		return nil, err
	}
	if len(handle) == 0 {
		err = fmt.Errorf("invalid handle for module")
		return
	}
	mod, err = s.store.LookupComposeModuleByNamespaceIDHandle(s.ctx, nsID, handle)
	if err != nil {
		return
	}
	if mod == nil {
		return mod, fmt.Errorf("module not found")
	}

	mod.Fields, _, err = s.store.SearchComposeModuleFields(s.ctx, cTypes.ModuleFieldFilter{ModuleID: []uint64{mod.ID}})
	if err != nil {
		return
	}

	return
}

// FindModuleByID will get module by ID
func (s seeder) FindModuleByID(nsID, mID uint64) (mod *cTypes.Module, err error) {
	if nsID == 0 {
		err = fmt.Errorf("invalid namespace ID")
		return
	}
	if mID == 0 {
		err = fmt.Errorf("invalid module ID")
		return
	}
	mod, err = s.modSvc.FindByID(s.ctx, nsID, mID)
	if err != nil {
		return
	}
	if mod == nil {
		return mod, fmt.Errorf("module not found")
	}

	return
}

// FindModuleByHandle will get module by handle
func (s seeder) FindModuleByHandle(nsID uint64, handle string) (mod *cTypes.Module, err error) {
	if nsID == 0 {
		err = fmt.Errorf("invalid namespace for module")
		return nil, err
	}
	if len(handle) == 0 {
		err = fmt.Errorf("invalid handle for module")
		return
	}
	mod, err = s.modSvc.FindByHandle(s.ctx, nsID, handle)
	if err != nil {
		return
	}
	if mod == nil {
		return mod, fmt.Errorf("module not found")
	}
	return
}

// GetModuleFromParams will get module of namespace using params
func (s seeder) GetModuleFromParams(params *RecordParams) (mod *cTypes.Module, err error) {
	if params == nil {
		return
	}

	var ns *cTypes.Namespace

	if params.NamespaceID > 0 {
		ns, err = s.LookupNamespaceByID(params.NamespaceID)
		if err != nil {
			return
		}
	}
	if len(params.NamespaceHandle) > 0 {
		ns, err = s.LookupNamespaceByHandle(params.NamespaceHandle)
		if err != nil {
			return
		}
	}

	if params.ModuleID > 0 {
		mod, err = s.LookupModuleByID(params.ModuleID)
		if err != nil {
			return
		}
	}
	if len(params.ModuleHandle) > 0 {
		mod, err = s.LookupModuleByHandle(ns.ID, params.ModuleHandle)
		if err != nil {
			return
		}
	}

	if mod == nil {
		return mod, fmt.Errorf("module not found")
	}
	return
}

// CreateRecord creates given no of record into DB
func (s seeder) CreateRecord(params RecordParams) (IDs []uint64, err error) {
	var (
		records []*cTypes.Record
		labels  []*lTypes.Label
	)

	m, err := s.GetModuleFromParams(&params)
	if err != nil {
		return
	}

	for i := 0; i < params.getLimit(); i++ {
		rec := &cTypes.Record{
			ID:          id.Next(),
			NamespaceID: m.NamespaceID,
			ModuleID:    m.ID,
			CreatedAt:   time.Now(),
		}

		for j, f := range m.Fields {
			err = s.setRecordValues(rec, uint(j), f)
			if err != nil {
				return nil, err
			}
		}

		IDs = append(IDs, rec.ID)
		records = append(records, rec)
		labels = append(labels, s.CreateLabel(
			rec.ID,
			rec.LabelResourceKind(),
			FakeDataLabel,
		))
	}

	err = s.store.CreateComposeRecord(s.ctx, m, records...)
	if err != nil {
		return
	}

	err = s.store.UpsertLabel(s.ctx, labels...)
	if err != nil {
		return
	}

	return
}

// setRecordValues will generate record values from third party and set to record
func (s seeder) setRecordValues(rec *cTypes.Record, place uint, field *cTypes.ModuleField) (err error) {
	// @todo: is it multi-value field? how many values? then create some multi-values

	var value string

	// skip the non required fields
	if !field.Required {
		fmt.Println("Skipping it due to not required")
		return
	}

	// return err for fields without name
	if len(field.Name) == 0 {
		return fmt.Errorf("invalid field name")
	}

	// check the type to fake the data
	if len(field.Kind) > 0 {
		value, err = s.faker.fakeValue(field.Name, field.Kind, valueOptions{})
		if err != nil {
			return fmt.Errorf("coudn't generate the value")
		}
	} else {
		return fmt.Errorf("unknown kind for field")
	}

	rec.Values = rec.Values.Set(&cTypes.RecordValue{
		RecordID: rec.ID,
		Name:     field.Name,
		Value:    value,
		// Ref:      0, // @todo: fake data with ref
		Place: place, // in case of multi-value field this is ++
	})

	return
}

// DeleteAllRecord clear all the fake user from DB
func (s seeder) DeleteAllRecord(mod *cTypes.Module) (err error) {
	filter := cTypes.RecordFilter{
		Labels: map[string]string{
			FakeDataLabel: FakeDataLabel,
		},
	}
	records, _, err := s.store.SearchComposeRecords(s.ctx, mod, filter)
	if err != nil {
		return
	}

	err = s.store.DeleteComposeRecord(s.ctx, mod, records...)
	if err != nil {
		return
	}

	return
}

// DeleteAll will delete all the fake data from the DB
func (s seeder) DeleteAll(params *RecordParams) (err error) {
	err = s.DeleteAllUser()
	if err != nil {
		return
	}

	if params != nil {
		m, err := s.GetModuleFromParams(params)
		if err != nil {
			err = s.DeleteAllRecord(m)
			if err != nil {
				return err
			}
		}
	}

	return
}

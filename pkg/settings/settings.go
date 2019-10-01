package settings

import (
	"context"
)

type (
	service struct {
		repository Repository
	}

	Service interface {
		With(ctx context.Context) Service

		FindByPrefix(prefix string) (vv ValueSet, err error)
		BulkSet(vv ValueSet) (err error)
		Set(v *Value) (err error)
		Get(name string, ownedBy uint64) (out *Value, err error)
		Delete(name string, ownedBy uint64) error
		GetGlobalString(name string) (out string, err error)
		GetGlobalBool(name string) (out bool, err error)
	}
)

func NewService(r Repository) Service {
	svc := &service{
		repository: r,
	}

	return svc
}

func (s service) With(ctx context.Context) Service {
	return &service{
		repository: s.repository.With(ctx),
	}
}

func (s service) FindByPrefix(prefix string) (ValueSet, error) {
	if vv, err := s.repository.Find(Filter{Prefix: prefix}); err != nil {
		return nil, err
	} else {
		return vv, nil
	}
}

func (s service) Set(v *Value) (err error) {
	return s.repository.Set(v)
}

func (s service) BulkSet(vv ValueSet) (err error) {
	return s.repository.BulkSet(vv)
}

func (s service) Get(name string, ownedBy uint64) (out *Value, err error) {
	return s.repository.Get(name, ownedBy)
}

func (s service) Delete(name string, ownedBy uint64) error {
	return s.repository.Delete(name, ownedBy)
}

func (s service) GetGlobalString(name string) (out string, err error) {
	const global = 0
	var v *Value

	if v, err = s.repository.Get(name, global); err == nil && v != nil {
		err = v.Value.Unmarshal(&out)
	}

	return
}

func (s service) GetGlobalBool(name string) (out bool, err error) {
	const global = 0
	var v *Value

	if v, err = s.repository.Get(name, global); err == nil && v != nil {
		err = v.Value.Unmarshal(&out)
	}

	return
}

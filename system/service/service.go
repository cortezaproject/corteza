package service

import (
	"sync"
)

type (
	db interface {
		Transaction(callback func() error) error
	}
)

var (
	o                   sync.Once
	DefaultAuth         AuthService
	DefaultUser         UserService
	DefaultTeam         TeamService
	DefaultOrganisation OrganisationService
)

func Init() {
	o.Do(func() {
		DefaultAuth = Auth()
		DefaultUser = User()
		DefaultTeam = Team()
		DefaultOrganisation = Organisation()
	})
}

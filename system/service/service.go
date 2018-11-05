package service

import (
	"sync"
)

var (
	o                   sync.Once
	DefaultUser         UserService
	DefaultTeam         TeamService
	DefaultOrganisation OrganisationService
)

func Init() {
	o.Do(func() {
		DefaultUser = User()
		DefaultTeam = Team()
		DefaultOrganisation = Organisation()
	})
}

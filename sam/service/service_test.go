package service

import (
	"github.com/crusttech/crust/sam/repository"
)

type (
	// For gomock to generate one-for-all mocked struct
	Repository interface {
		repository.Interfaces
	}
)

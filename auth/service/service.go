package service

import (
	"sync"
)

var (
	o           sync.Once
	DefaultUser UserService
)

func Init() {
	o.Do(func() {
		DefaultUser = User()
	})
}

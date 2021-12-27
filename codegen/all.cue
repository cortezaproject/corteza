package codegen

import (
  "github.com/cortezaproject/corteza-server/codegen/schema"
)

all: [...schema.#codegen] &
	rbacAccessControl +
	rbacTypes +
	[] // placeholder

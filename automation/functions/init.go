package functions

import "github.com/cortezaproject/corteza-server/automation/types/fn"

const (
	baseRef = "base"
)

func List() []*fn.Function {
	return httpSenders()
}

package sam

import (
	"github.com/pkg/errors"
)

func (*Team) Edit(r *teamEditRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Team.edit")
}

func (*Team) Remove(r *teamRemoveRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Team.remove")
}

func (*Team) Read(r *teamReadRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Team.read")
}

func (*Team) Search(r *teamSearchRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Team.search")
}

func (*Team) Archive(r *teamArchiveRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Team.archive")
}

func (*Team) Move(r *teamMoveRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Team.move")
}

func (*Team) Merge(r *teamMergeRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Team.merge")
}

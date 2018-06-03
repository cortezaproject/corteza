package sam

import (
	"github.com/pkg/errors"
)

func (c *Channel) Edit(r *channelEditRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Channel.edit")
}
func (c *Channel) Remove(r *channelRemoveRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Channel.remove")
}
func (c *Channel) Read(r *channelReadRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Channel.read")
}
func (c *Channel) Search(r *channelSearchRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Channel.search")
}
func (c *Channel) Archive(r *channelArchiveRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Channel.archive")
}
func (c *Channel) Move(r *channelMoveRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Channel.move")
}
func (c *Channel) Merge(r *channelMergeRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Channel.merge")
}

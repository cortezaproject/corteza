package service

import (
	"context"
	"github.com/crusttech/crust/sam/repository"
	"github.com/crusttech/crust/sam/types"
)

type (
	channel struct {
		rpo channelRepository
		//
		//sec struct {
		//	ch channelSecurity
		//}
	}

	channelRepository interface {
		repository.Transactionable
		repository.Channel
	}

	//channelSecurity interface {
	//	CanRead(ctx context.Context, ch *types.Channel) bool
	//}
)

func Channel() *channel {
	var svc = &channel{}

	svc.rpo = repository.New()
	//svc.sec.ch = ChannelSecurity(svc.rpo)

	return svc
}

func (svc channel) FindByID(ctx context.Context, id uint64) (ch *types.Channel, err error) {
	ch, err = svc.rpo.FindChannelByID(id)
	if err != nil {
		return
	}

	//if !svc.sec.ch.CanRead(ch) {
	//	return nil, errors.New("Not allowed to access channel")
	//}

	return
}

func (svc channel) Find(ctx context.Context, filter *types.ChannelFilter) ([]*types.Channel, error) {
	// @todo: permission check to return only channels that channel has access to
	// @todo: actual searching not just a full select
	return svc.rpo.FindChannels(filter)
}

func (svc channel) Create(ctx context.Context, mod *types.Channel) (*types.Channel, error) {
	// @todo: topic channelEvent/log entry
	// @todo: channel name cmessage/log entry
	// @todo: permission check if channel can add channel
	return svc.rpo.CreateChannel(mod)
}

func (svc channel) Update(ctx context.Context, mod *types.Channel) (*types.Channel, error) {
	// @todo: load current entry and merge changes
	// @todo: topic change channelEvent/log entry
	// @todo: channel name change channelEvent/log entry
	// @todo: permission check if current user can edit channel
	// @todo: make sure archived & deleted entries can not be edited
	// @todo: handle channel movinga
	// @todo: handle channel archiving

	return svc.rpo.UpdateChannel(mod)
}

func (svc channel) Delete(ctx context.Context, id uint64) error {
	// @todo: make history unavailable
	// @todo: notify users that channel has been removed (remove from web UI)
	// @todo: permissions check if current user can remove channel
	return svc.rpo.DeleteChannelByID(id)
}

func (svc channel) Archive(ctx context.Context, id uint64) error {
	// @todo: make history unavailable
	// @todo: notify users that channel has been removed (remove from web UI)
	// @todo: permissions check if current user can remove channel
	return svc.rpo.ArchiveChannelByID(id)
}

func (svc channel) Unarchive(ctx context.Context, id uint64) error {
	// @todo: make history unavailable
	// @todo: notify users that channel has been removed (remove from web UI)
	// @todo: permissions check if current user can remove channel
	return svc.rpo.UnarchiveChannelByID(id)
}

//// @todo temp location, move this somewhere else
//type (
//	nativeChannelSec struct {
//		rpo struct {
//			ch nativeChannelSecChRepo
//		}
//	}
//
//	nativeChannelSecChRepo interface {
//		FindMember(ctx context.Context, channelId uint64, userId uint64) (*types.User, error)
//	}
//)
//
//func ChannelSecurity(chRpo nativeChannelSecChRepo) channelSecurity {
//	var sec = &nativeChannelSec{}
//
//	sec.rpo.ch = chRpo
//
//	return sec
//}
//
//// Current user can read the channel if he is a member
//func (sec nativeChannelSec) CanRead(ctx context.Context, ch *types.Channel) bool {
//	// @todo check if channel is public?
//
//	var currentUserID = auth.GetIdentityFromContext(ctx).Identity()
//
//	user, err := sec.rpo.FindMember(ch.ID, currentUserID)
//
//	return err != nil && user.Valid()
//}

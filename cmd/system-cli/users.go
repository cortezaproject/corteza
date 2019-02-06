package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	crmRepository "github.com/crusttech/crust/crm/repository"
	"github.com/crusttech/crust/internal/payload"
	messagingRepository "github.com/crusttech/crust/messaging/repository"
	"github.com/crusttech/crust/system/service"
	"github.com/crusttech/crust/system/types"
)

func cliExecUsers(commands ...string) {
	if len(commands) == 0 {
		return
	}

	switch commands[0] {
	case "list":
		cliExecUsersList(commands[1:]...)
	case "merge":
		cliExecUsersMerge(commands[1:]...)
	}
}

func cliExecUsersList(params ...string) {
	var (
		upd, del string

		err error
		uu  types.UserSet
		ctx = context.Background()
	)

	uf := &types.UserFilter{
		OrderBy: "updated_at",
	}

	if uu, err = service.DefaultUser.With(ctx).Find(uf); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("ID                   Updated    Deleted    [email / name / username]")
	for _, u := range uu {
		upd, del = "---- -- --", "---- -- --"

		if u.UpdatedAt != nil {
			upd = u.UpdatedAt.Format("2006-01-02")
		}

		if u.DeletedAt != nil {
			upd = u.DeletedAt.Format("2006-01-02")
		}

		fmt.Printf(
			"%20d %s %s %s\n",
			u.ID,
			upd,
			del,
			u.Email+" / "+u.Name+" / "+u.Username)
	}
}

func cliExecUsersMerge(params ...string) {
	var (
		err  error
		uu   = make([]*types.User, len(params))
		refs = make([]*userRefs, len(params))
		ids  = payload.ParseUInt64s(params)
		ctx  = context.Background()
	)

	if len(ids) < 2 {
		fmt.Printf("Expecting 2+ user IDs (2nd, 3rd ... user ID will be merged into first one\n")
		os.Exit(1)
	}

	for i, id := range ids {
		if id == 0 {
			fmt.Printf("Error: Invalid user ID %q\n", params[i])
			os.Exit(1)
		}

		if uu[i], err = service.DefaultUser.With(ctx).FindByID(id); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	}

	db := messagingRepository.DB(ctx)

	mergers := []struct {
		label string
		count func(userID uint64) (c int, err error)
		merge func(userID, target uint64) (err error)
	}{
		{label: "MsgOw",
			count: messagingRepository.Message(ctx, db).CountOwned,
			merge: messagingRepository.Message(ctx, db).ChangeOwner},
		{label: "MTags",
			count: messagingRepository.Message(ctx, db).CountUserTags,
			merge: messagingRepository.Message(ctx, db).ChangeUserTag},
		{label: "ChCre",
			count: messagingRepository.Channel(ctx, db).CountCreated,
			merge: messagingRepository.Channel(ctx, db).ChangeCreator},
		{label: "Membr",
			count: messagingRepository.ChannelMember(ctx, db).CountMemberships,
			merge: messagingRepository.ChannelMember(ctx, db).ChangeMembership},
		{label: "AttOw",
			count: messagingRepository.Attachment(ctx, db).CountOwned,
			merge: messagingRepository.Attachment(ctx, db).ChangeOwnership},
		{label: "Menti",
			count: messagingRepository.Mention(ctx, db).CountMentions,
			merge: messagingRepository.Mention(ctx, db).ChangeMention},
		{label: "Unrd",
			count: messagingRepository.Unread(ctx, db).CountOwned,
			merge: messagingRepository.Unread(ctx, db).ChangeOwner},
		{label: "CRAut",
			count: crmRepository.Record(ctx, db).CountAuthored,
			merge: crmRepository.Record(ctx, db).ChangeAuthor},
		{label: "CRRef",
			count: crmRepository.Record(ctx, db).CountReferenced,
			merge: crmRepository.Record(ctx, db).ChangeReferences},
	}

	count := func(u *types.User, r *userRefs) (out string) {
		out = fmt.Sprintf(
			"%20d | %-40s",
			u.ID,
			u.Email+" / "+u.Name+" / "+u.Username,
		)

		for _, m := range mergers {
			if count, err := m.count(u.ID); err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			} else {
				out = out + fmt.Sprintf(" | %5d", count)
			}
		}

		return out + fmt.Sprintln()
	}

	stats := fmt.Sprintf(
		"%20s | %40s",
		"ID",
		"Email",
	)

	for _, m := range mergers {
		stats = stats + fmt.Sprintf(" | %5s", m.label)
	}

	stats = stats + fmt.Sprintln() + fmt.Sprintf("Merge %d users:\n", len(uu)-1)
	for i := 1; i < len(uu); i++ {
		stats = stats + count(uu[i], refs[i])
	}
	stats = stats + fmt.Sprintln("Target") + count(uu[0], refs[0])

	fmt.Println(stats)

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Merge [y/N]? ")
	text, _ := reader.ReadByte()
	if "y" != strings.ToLower(string(text)) {
		os.Exit(0)
	}

	for i := 1; i < len(uu); i++ {
		for _, m := range mergers {
			if err := m.merge(uu[i].ID, uu[0].ID); err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
		}
	}

	fmt.Println("Done.")
}

type userRefs struct {
	messagesCreated int
	messagesTagged  int
}

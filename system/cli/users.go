package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/crusttech/crust/system/service"
	"github.com/crusttech/crust/system/types"
)

func UsersList() {
	uf := &types.UserFilter{
		OrderBy: "updated_at",
	}

	users, err := service.DefaultUser.With(context.Background()).Find(uf)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("ID                   Updated    Deleted    [email / name / username]")
	for _, u := range users {
		upd, del := "---- -- --", "---- -- --"

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

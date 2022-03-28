package main

import (
	"context"
	"log"

	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/store/adapters/rdbms/drivers/postgres"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/davecgh/go-spew/spew"
)

func main() {
	var (
		ctx = context.Background()
	)

	//spew.Config.DisablePointerAddresses = true
	spew.Config.MaxDepth = 4
	s, err := postgres.Connect(context.Background(), "postgres://darh@localhost:5432/corteza_2022_3?sslmode=disable&")
	if err != nil {
		log.Fatal("postgresql.Open: ", err)
	}

	f := types.UserFilter{}

	f.Sorting, _ = filter.NewSorting("email")

	spew.Dump(s.SearchUsers(ctx, f))
}

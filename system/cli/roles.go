package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/crusttech/crust/system/internal/repository"
)

func RolesReset() {
	ctx := context.Background()
	db := repository.DB(ctx)

	err := repository.Role(ctx, db).Reset()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Everyone and Administrators role were reset.")
}

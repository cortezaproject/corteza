package app

// Registers all supported store backends
import (
	_ "github.com/crusttech/human/server/store/adapters/api/drivers/rest"
	_ "github.com/crusttech/human/server/store/adapters/rdbms/drivers/mssql"
	_ "github.com/crusttech/human/server/store/adapters/rdbms/drivers/mysql"
	_ "github.com/crusttech/human/server/store/adapters/rdbms/drivers/postgres"
	_ "github.com/crusttech/human/server/store/adapters/rdbms/drivers/sqlite"
)

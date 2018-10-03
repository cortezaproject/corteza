package db

import (
	"context"
	"testing"
	"time"

	"github.com/titpetric/dockertest"
	"github.com/titpetric/factory"
)

func TestMigrations(t *testing.T) {
	args := []string{
		"--rm",
		"-e", "MYSQL_ROOT_PASSWORD=root",
		"-e", "MYSQL_DATABASE=test",
	}

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(30*time.Second))
	defer cancel()

	count := 0

	mysql, err := dockertest.RunContainerContext(ctx, "titpetric/percona-xtrabackup", "3306", func(addr string) error {
		factory.Database.Add("default", "root:root@tcp("+addr+")/test?collation=utf8mb4_general_ci")
		_, err := factory.Database.Get()
		count++
		time.Sleep(time.Second)
		return err
	}, args...)
	defer mysql.Terminate() // Shutdown()
	if err != nil {
		t.Fatalf("Error starting mysql: %#v", err)
	}

	db := factory.Database.MustGet()
	if err := Migrate(db); err != nil {
		t.Fatalf("Unexpected error: %#v", err)
	}
}

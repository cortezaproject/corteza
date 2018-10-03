package db

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"

	_ "github.com/crusttech/crust/sam/db/mysql"
	"github.com/pkg/errors"
	"github.com/rakyll/statik/fs"
	"github.com/titpetric/factory"
)

//go:generate statik -p mysql -Z -f -src=schema/mysql

func statements(contents []byte, err error) ([]string, error) {
	if err != nil {
		return []string{}, err
	}
	return regexp.MustCompilePOSIX(";$").Split(string(contents), -1), nil
}

func Migrate(db *factory.DB) error {
	statikFS, err := fs.New()
	if err != nil {
		return errors.Wrap(err, "Error creating statik filesystem")
	}

	var files []string

	fs.Walk(statikFS, "/", func(filename string, info os.FileInfo, err error) error {
		if len(filename) > 4 && filename[len(filename)-4:] == ".sql" {
			files = append(files, filename)
		}
		return nil
	})

	sort.Strings(files)

	if len(files) == 0 {
		return errors.New("No files encoded for migration, need at least one SQL file")
	}

	// @todo: create table migrations by hand, log each filename status/date
	/*
		create table if not exists migrations (
			filename string, status [ok, fail], message string (may be error text?), stamp datetime
		)
	*/

	for _, filename := range files {
		// @todo: select * from migrations to check if (filename) was processed, skip if yes

		stmts, err := statements(fs.ReadFile(statikFS, filename))
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("Error applying migration %s", filename))
		}

		up := func() error {
			log.Println("Running migration for", filename)
			for _, query := range stmts {
				if _, err := db.Exec(query); err != nil {
					if fmt.Sprintf("%s", err) != "exec query failed: Error 1065: Query was empty" {
						return err
					}
				}
			}
			log.Println("Migration done OK")

			// @todo: insert/update into migrations with (filename)
			return nil
		}

		if err := db.Transaction(up); err != nil {
			return err
		}
	}

	return nil
}

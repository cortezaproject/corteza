package db

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/goware/statik/fs"
	"github.com/pkg/errors"
	"github.com/titpetric/factory"
	"go.uber.org/zap"

	"github.com/cortezaproject/corteza-server/system/db/mysql"
)

func statements(contents []byte, err error) ([]string, error) {
	if err != nil {
		return []string{}, err
	}
	return regexp.MustCompilePOSIX(";$").Split(string(contents), -1), nil
}

func Migrate(db *factory.DB, log *zap.Logger) error {
	log = log.Named("database.migrations")

	statikFS, err := fs.New(mysql.Asset)
	if err != nil {
		return errors.Wrap(err, "error creating statik filesystem")
	}

	var files []string

	fn := func(filename string, info os.FileInfo, err error) error {
		_ = err
		matched, err := filepath.Match("/*.up.sql", filename)
		if matched {
			files = append(files, filename)
		}

		return err
	}

	if err := fs.Walk(statikFS, "/", fn); err != nil {
		return errors.Wrap(err, "error when listing files for migrations")
	}

	sort.Strings(files)

	if len(files) == 0 {
		return errors.New("no files encoded for migration, need at least one SQL file")
	}

	migrate := func(filename string, useLog bool) error {
		status := migration{
			Project:  "system",
			Filename: filename,
		}
		if useLog {
			if err := db.Get(&status, "select * from migrations where project=? and filename=?", status.Project, status.Filename); err != nil {
				return err
			}
			if status.Status == "ok" {
				return nil
			}
		}

		up := func() error {
			stmts, err := statements(fs.ReadFile(statikFS, filename))
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("error reading migration %s", filename))
			}

			log.Debug("Running migration", zap.String("filename", filename))
			for idx, query := range stmts {
				if strings.TrimSpace(query) != "" && idx >= status.StatementIndex {
					status.StatementIndex = idx
					if _, err := db.Exec(query); err != nil {
						log.Debug("migration error ", zap.String("filename", filename), zap.Error(err))
						return err
					}
				}
			}

			status.Status = "ok"
			return nil
		}

		err := db.Transaction(up)
		if err != nil {
			status.Status = err.Error()
		}
		if useLog {
			if err := db.Replace("migrations", status); err != nil {
				return errors.Wrap(err, "migration update failed")
			}
		}
		return err
	}

	if err := migrate("/migrations.sql", false); err != nil {
		return err
	}

	db.Exec("LOCK TABLE migrations WRITE;")
	defer db.Exec("UNLOCK TABLES")

	for _, filename := range files {
		if err := migrate(filename, true); err != nil {
			return err
		}
	}

	return nil
}

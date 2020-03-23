package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/cortezaproject/corteza-server/compose/repository"
	"github.com/cortezaproject/corteza-server/compose/service"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/cli"
	mgg "github.com/cortezaproject/corteza-server/pkg/migrate"
	mgt "github.com/cortezaproject/corteza-server/pkg/migrate/types"
)

func Migrator() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Migrate",

		Run: func(cmd *cobra.Command, args []string) {
			var (
				ctx      = auth.SetSuperUserContext(cli.Context())
				nsFlag   = cmd.Flags().Lookup("namespace").Value.String()
				srcFlag  = cmd.Flags().Lookup("src").Value.String()
				metaFlag = cmd.Flags().Lookup("meta").Value.String()
				ns       *types.Namespace
				err      error

				mg []mgt.Migrateable
			)

			svcNs := service.DefaultNamespace.With(ctx)

			if nsFlag == "" {
				cli.HandleError(errors.New("ns.undefined"))
			}
			if srcFlag == "" {
				cli.HandleError(errors.New("src.undefined"))
			}

			if namespaceID, _ := strconv.ParseUint(nsFlag, 10, 64); namespaceID > 0 {
				ns, err = svcNs.FindByID(namespaceID)
				if err != repository.ErrNamespaceNotFound {
					cli.HandleError(err)
				}
			} else if ns, err = svcNs.FindByHandle(nsFlag); err != nil {
				if err != repository.ErrNamespaceNotFound {
					cli.HandleError(err)
				}
			}

			// load src files
			err = filepath.Walk(srcFlag, func(path string, info os.FileInfo, err error) error {
				if strings.HasSuffix(info.Name(), ".csv") {
					file, err := os.Open(path)
					if err != nil {
						log.Fatal(err)
					}

					ext := filepath.Ext(info.Name())
					name := info.Name()[0 : len(info.Name())-len(ext)]
					mm := migrateableSource(mg, name)
					mm.Name = name
					mm.Path = path
					mm.Source = file

					mg = migrateableAdd(mg, mm)
				}
				return nil
			})

			// load meta files
			if metaFlag != "" {
				err = filepath.Walk(metaFlag, func(path string, info os.FileInfo, err error) error {
					if strings.HasSuffix(info.Name(), ".map.json") {
						file, err := os.Open(path)
						if err != nil {
							log.Fatal(err)
						}

						ext := filepath.Ext(info.Name())
						// @todo improve this!!
						name := info.Name()[0 : len(info.Name())-len(ext)-4]
						mm := migrateableSource(mg, name)
						mm.Name = name
						mm.Map = file

						mg = migrateableAdd(mg, mm)
					} else if strings.HasSuffix(info.Name(), ".join.json") {
						file, err := os.Open(path)
						if err != nil {
							log.Fatal(err)
						}

						ext := filepath.Ext(info.Name())
						// @todo improve this!!
						name := info.Name()[0 : len(info.Name())-len(ext)-5]
						mm := migrateableSource(mg, name)
						mm.Name = name
						mm.Join = file

						mg = migrateableAdd(mg, mm)
					} else if strings.HasSuffix(info.Name(), ".value.json") {
						file, err := os.Open(path)
						if err != nil {
							log.Fatal(err)
						}

						ext := filepath.Ext(info.Name())
						// @todo improve this!!
						name := info.Name()[0 : len(info.Name())-len(ext)-6]
						mm := migrateableSource(mg, name)
						mm.Name = name

						var vmp map[string]map[string]string
						src, _ := ioutil.ReadAll(file)
						err = json.Unmarshal(src, &vmp)
						if err != nil {
							log.Fatal(err)
						}
						mm.ValueMap = vmp

						mg = migrateableAdd(mg, mm)
					}
					return nil
				})

				if err != nil {
					panic(err)
				}
			}

			// clean up migrateables
			hasW := false
			out := make([]mgt.Migrateable, 0)
			for _, m := range mg {
				if m.Source != nil {
					out = append(out, m)
				} else {
					hasW = true
					spew.Dump("[warning] migrationNode.missingSource " + m.Name)
				}
			}

			if hasW {
				var rsp string
				fmt.Print("warnings detected; continue [y/N]? ")
				_, err := fmt.Scanln(&rsp)

				if err != nil || rsp != "y" && rsp != "yes" {
					log.Fatal("migration aborted due to warnings")
				}
			}

			err = mgg.Migrate(out, ns, ctx)
			if err != nil {
				panic(err)
			}

		},
	}

	cmd.Flags().String("namespace", "", "Import into namespace (by ID or string)")
	cmd.Flags().String("src", "", "Directory with migration files")
	cmd.Flags().String("meta", "", "Directory with meta files")

	return cmd
}

// small helper functions for migrateable node management
func migrateableSource(mg []mgt.Migrateable, name string) mgt.Migrateable {
	for _, m := range mg {
		if m.Name == name {
			return m
		}
	}

	return mgt.Migrateable{}
}

func migrateableAdd(mg []mgt.Migrateable, mm mgt.Migrateable) []mgt.Migrateable {
	for i, m := range mg {
		if m.Name == mm.Name {
			mg[i] = mm
			return mg
		}
	}
	return append(mg, mm)
}

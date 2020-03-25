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

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/cortezaproject/corteza-server/compose/repository"
	"github.com/cortezaproject/corteza-server/compose/service"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/cli"
	ngi "github.com/cortezaproject/corteza-server/pkg/ngimporter"
	ngt "github.com/cortezaproject/corteza-server/pkg/ngimporter/types"
)

func NGImporter() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ng-importer",
		Short: "Importer next-gen",

		Run: func(cmd *cobra.Command, args []string) {
			var (
				ctx      = auth.SetSuperUserContext(cli.Context())
				nsFlag   = cmd.Flags().Lookup("namespace").Value.String()
				srcFlag  = cmd.Flags().Lookup("src").Value.String()
				metaFlag = cmd.Flags().Lookup("meta").Value.String()
				ns       *types.Namespace
				err      error

				iss []ngt.ImportSource
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
					mm := importSource(iss, name)
					mm.Name = name
					mm.Path = path
					mm.Source = file

					iss = addImportSource(iss, mm)
				}
				return nil
			})

			// load meta files
			if metaFlag != "" {
				err = filepath.Walk(metaFlag, func(path string, info os.FileInfo, err error) error {
					if strings.HasSuffix(info.Name(), ngt.MetaMapExt) {
						file, err := os.Open(path)
						if err != nil {
							log.Fatal(err)
						}

						// @todo improve this!!
						name := info.Name()[0 : len(info.Name())-len(ngt.MetaMapExt)]
						mm := importSource(iss, name)
						mm.Name = name
						mm.DataMap = file

						iss = addImportSource(iss, mm)
					} else if strings.HasSuffix(info.Name(), ngt.MetaJoinExt) {
						file, err := os.Open(path)
						if err != nil {
							log.Fatal(err)
						}

						// @todo improve this!!
						name := info.Name()[0 : len(info.Name())-len(ngt.MetaJoinExt)]
						mm := importSource(iss, name)
						mm.Name = name
						mm.SourceJoin = file

						iss = addImportSource(iss, mm)
					} else if strings.HasSuffix(info.Name(), ngt.MetaValueExt) {
						file, err := os.Open(path)
						if err != nil {
							log.Fatal(err)
						}

						// @todo improve this!!
						name := info.Name()[0 : len(info.Name())-len(ngt.MetaValueExt)]
						mm := importSource(iss, name)
						mm.Name = name

						var vmp map[string]map[string]string
						src, _ := ioutil.ReadAll(file)
						err = json.Unmarshal(src, &vmp)
						if err != nil {
							log.Fatal(err)
						}
						mm.ValueMap = vmp

						iss = addImportSource(iss, mm)
					}
					return nil
				})

				if err != nil {
					log.Fatal(err)
				}
			}

			// remove nodes with no import source
			hasW := false
			out := make([]ngt.ImportSource, 0)
			for _, is := range iss {
				if is.Source != nil {
					out = append(out, is)
				} else {
					hasW = true
					log.Println("[warning] ImportSOurce.missingSource " + is.Name)
				}
			}

			if hasW {
				var rsp string
				fmt.Print("warnings detected; continue [y/N]? ")
				_, err := fmt.Scanln(&rsp)
				rsp = strings.ToLower(rsp)

				if err != nil || rsp != "y" && rsp != "yes" {
					log.Fatal("import aborted due to warnings")
				}
			}

			err = ngi.Import(ctx, out, ns)
			if err != nil {
				panic(err)
			}

		},
	}

	cmd.Flags().String("namespace", "", "Import into namespace (by ID or string)")
	cmd.Flags().String("src", "", "Directory with import files")
	cmd.Flags().String("meta", "", "Directory with import meta files")

	return cmd
}

// small helper functions for import source node management

func importSource(iss []ngt.ImportSource, name string) ngt.ImportSource {
	for _, is := range iss {
		if is.Name == name {
			return is
		}
	}

	return ngt.ImportSource{}
}

func addImportSource(iss []ngt.ImportSource, nis ngt.ImportSource) []ngt.ImportSource {
	for i, is := range iss {
		if is.Name == nis.Name {
			iss[i] = nis
			return iss
		}
	}
	return append(iss, nis)
}

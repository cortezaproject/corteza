package commands

import (
	"embed"
	"os"
	"path"

	"github.com/cortezaproject/corteza/server/auth"
	"github.com/cortezaproject/corteza/server/pkg/cli"
	"github.com/spf13/cobra"
)

func assets(app serviceInitializer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "assets",
		Short: "Authentication flow assets (styling, images) and templates",
	}

	exportCmd := &cobra.Command{
		Use:   "export",
		Short: "Exports embedded assets into provided path (must exists)",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			assetsRoot := app.Options().Auth.AssetsPath
			if len(args) > 0 {
				assetsRoot = args[0]
			}

			if len(assetsRoot) == 0 {
				cmd.Println("can not export, no path provided and AUTH_ASSETS_PATH is empty")
				return
			}

			if _, err := os.Stat(assetsRoot); err != nil {
				cli.HandleError(err)
			}

			emb := map[string]embed.FS{
				"public":    auth.PublicAssets,
				"templates": auth.Templates,
			}

			var (
				fh        *os.File
				buf       []byte
				dst       string
				src       string
				assetsSub string
			)
			for dir, efs := range emb {
				assetsSub = path.Join(assetsRoot, dir)

				cmd.Printf("exporting auth assets to %s\n", dir)
				if _, err := os.Stat(assetsSub); os.IsNotExist(err) {
					cli.HandleError(os.Mkdir(assetsSub, 0755))
					cmd.Println("directory created")
				} else if err != nil {
					cli.HandleError(err)
				}

				cc, err := efs.ReadDir(path.Join("assets", dir))
				if err != nil {
					cli.HandleError(err)
				}

				for _, c := range cc {
					src = path.Join("assets", dir, c.Name())
					dst = path.Join(assetsSub, c.Name())

					if c.IsDir() {
						cmd.Println("skipping directory:", assetsRoot)
						continue
					}

					cmd.Print("exporting asset ", dst, ": ")

					buf, err = efs.ReadFile(src)
					if err != nil {
						cmd.Println("\n => error:", err.Error())
						continue
					}

					fh, err = os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0666)
					if os.IsExist(err) {
						cmd.Println("exists")
						continue
					}
					if err != nil {
						cmd.Println("\n => error:", err.Error())
						continue
					}

					_, err = fh.Write(buf)
					if err != nil {
						cmd.Println("\n => error:", err.Error())
						continue
					}

					cmd.Println("ok")
				}
			}

		},
	}

	cmd.AddCommand(exportCmd)

	return cmd
}

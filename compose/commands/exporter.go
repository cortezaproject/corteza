package commands

import (
	"context"
	"regexp"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"github.com/cortezaproject/corteza-server/compose/service"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/internal/auth"
	"github.com/cortezaproject/corteza-server/pkg/cli"
	"github.com/cortezaproject/corteza-server/pkg/handle"
)

func Exporter(ctx context.Context, c *cli.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export",

		Run: func(cmd *cobra.Command, args []string) {

			c.InitServices(ctx, c)

			ctx = auth.SetSuperUserContext(ctx)

			mm, _, err := service.DefaultModule.Find(types.ModuleFilter{NamespaceID: 88714882739863655})
			cli.HandleError(err)

			y := yaml.NewEncoder(cmd.OutOrStdout())

			out := Compose{
				Modules: expModules(mm),
			}

			cli.HandleError(y.Encode(out))
		},
	}

	return cmd
}

type (
	Compose struct {
		Modules []Module
	}

	Module struct {
		Name   string
		Handle string
		Meta   string
		Fields map[string]Field
	}

	Field struct {
		Label string
		Kind  string

		Options string

		Private  bool
		Required bool
		Visible  bool
		Multi    bool
		// DefaultValue string
	}
)

func expModules(mm types.ModuleSet) (o []Module) {
	o = make([]Module, len(mm))

	for i, m := range mm {
		o[i] = Module{
			Name:   m.Name,
			Handle: makeHandle(m.Handle, m.Name),
			Meta:   m.Meta.String(),
			Fields: expModuleFields(m.Fields),
		}

		if o[i].Handle == "" {
			h := regexp.MustCompile(`^[^A-Za-z][^0-9A-Za-z_\-.][^A-Za-z0-9]$`).ReplaceAllString(o[i].Name, "")
			if handle.IsValid(h) {
				o[i].Handle = h
			}
		}
	}

	return
}

func expModuleFields(ff types.ModuleFieldSet) (o map[string]Field) {
	o = make(map[string]Field)

	for _, f := range ff {
		o[f.Name] = Field{
			Label:    f.Label,
			Kind:     f.Kind,
			Options:  f.Options.String(),
			Private:  f.Private,
			Required: f.Required,
			Visible:  f.Visible,
			Multi:    f.Multi,
		}
	}

	return
}

func makeHandle(h, n string) string {
	if h == "" {
		h = regexp.MustCompile(`^[^A-Za-z][^0-9A-Za-z_\-.][^A-Za-z0-9]$`).ReplaceAllString(n, "")
		if handle.IsValid(h) {
			return h
		}
	}

	return n
}

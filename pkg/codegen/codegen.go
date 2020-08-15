package codegen

import (
	"github.com/cortezaproject/corteza-server/pkg/cli"
	"strings"
	"text/template"
)

type (
	definitions struct {
		App string

		Rest    []*restDef
		Actions []*actionsDef
		Events  []*eventsDef
		Types   []*typesDef
		Store   []*storeDef
	}
)

func Proc() {
	var (
		err error

		def = &definitions{}

		tpls = template.New("").Funcs(map[string]interface{}{
			"camelCase":       camelCase,
			"pubIdent":        pubIdent,
			"unpubIdent":      unpubIdent,
			"toLower":         strings.ToLower,
			"cc2underscore":   cc2underscore,
			"normalizeImport": normalizeImport,
			"comment": func(text string, skip1st bool) string {
				ll := strings.Split(text, "\n")
				s := 0
				out := ""
				if skip1st {
					s = 1
					out = ll[0] + "\n"
				}

				for ; s < len(ll); s++ {
					out += "// " + ll[s] + "\n"
				}

				return out
			},
		})
	)

	tpls = template.Must(tpls.ParseGlob("pkg/codegen/assets/*.tpl"))

	if def.Actions, err = procActions(); err != nil {
		cli.HandleError(err)
	} else {
		cli.HandleError(genActions(tpls, def.Actions))
	}

	if def.Events, err = procEvents(); err != nil {
		cli.HandleError(err)
	} else {
		cli.HandleError(genEvents(tpls, def.Events))
	}

	if def.Types, err = procTypes(); err != nil {
		cli.HandleError(err)
	} else {
		cli.HandleError(genTypes(tpls, def.Types))
	}

	if def.Rest, err = procRest(); err != nil {
		cli.HandleError(err)
	} else {
		cli.HandleError(genRest(tpls, def.Rest))
	}

	if def.Store, err = procStore(); err != nil {
		cli.HandleError(err)
	} else {
		cli.HandleError(genStore(tpls, def.Store))
	}
}

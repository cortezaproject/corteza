package options

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// pkg/options/template.yaml

type (
	TemplateOpt struct {
		RendererGotenbergAddress string `env:"TEMPLATE_RENDERER_GOTENBERG_ADDRESS"`
		RendererGotenbergEnabled bool   `env:"TEMPLATE_RENDERER_GOTENBERG_ENABLED"`
	}
)

// Template initializes and returns a TemplateOpt with default values
func Template() (o *TemplateOpt) {
	o = &TemplateOpt{
		RendererGotenbergAddress: "",
		RendererGotenbergEnabled: false,
	}

	fill(o)

	// Function that allows access to custom logic inside the parent function.
	// The custom logic in the other file should be like:
	// func (o *Template) Defaults() {...}
	func(o interface{}) {
		if def, ok := o.(interface{ Defaults() }); ok {
			def.Defaults()
		}
	}(o)

	return
}

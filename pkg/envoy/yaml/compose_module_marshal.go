package yaml

import (
	"context"
	"encoding/json"

	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
)

func composeModuleFromResource(r *resource.ComposeModule, cfg *EncoderConfig) *composeModule {
	return &composeModule{
		res:           r.Res,
		encoderConfig: cfg,
	}
}

func (n *composeModule) Prepare(ctx context.Context, state *envoy.ResourceState) (err error) {
	mod, ok := state.Res.(*resource.ComposeModule)
	if !ok {
		return encoderErrInvalidResource(resource.COMPOSE_MODULE_RESOURCE_TYPE, state.Res.ResourceType())
	}

	// Get the related namespace
	n.relNs = resource.FindComposeNamespace(state.ParentResources, mod.RefNs.Identifiers)
	if n.relNs == nil {
		return resource.ComposeNamespaceErrUnresolved(mod.RefNs.Identifiers)
	}

	n.res = mod.Res
	n.refNamespace = relNsToRef(n.relNs)

	n.fields = make(composeModuleFieldSet, 0, len(mod.Res.Fields))
	for _, f := range mod.Res.Fields {
		cmf := &composeModuleField{
			res:  f,
			expr: composeModuleFieldExpr(f.Expressions),
		}

		if f.Kind == "Record" {
			refMod := f.Options.String("module")
			if refMod == "" {
				refMod = f.Options.String("moduleID")
			}
			cmf.relMod = resource.FindComposeModule(state.ParentResources, resource.MakeIdentifiers(refMod))
			if cmf.relMod == nil {
				return resource.ComposeModuleErrUnresolved(resource.MakeIdentifiers(refMod))
			}
		}

		n.fields = append(n.fields, cmf)
	}

	if err != nil {
		return err
	}

	return nil
}

func (n *composeModule) Encode(ctx context.Context, doc *Document, state *envoy.ResourceState) (err error) {
	if n.res.ID <= 0 {
		n.res.ID = nextID()
	}

	if state.Conflicting {
		return nil
	}

	n.ts, err = resource.MakeTimestampsCUDA(&n.res.CreatedAt, n.res.UpdatedAt, n.res.DeletedAt, nil).
		Model(n.encoderConfig.TimeLayout, n.encoderConfig.Timezone)
	if err != nil {
		return err
	}

	// Fields
	for _, f := range n.fields {
		// Timestaps
		n.ts, err = resource.MakeTimestampsCUDA(&f.res.CreatedAt, f.res.UpdatedAt, f.res.DeletedAt, nil).
			Model(n.encoderConfig.TimeLayout, n.encoderConfig.Timezone)
		if err != nil {
			return err
		}
	}

	// @todo skip eval?

	if n.encoderConfig.CompactOutput {
		err = doc.nestComposeModule(n.refNamespace, n)
	} else {
		doc.addComposeModule(n)
	}

	return err
}

func (c *composeModule) MarshalYAML() (interface{}, error) {
	// Get a struct from the raw JSON string
	auxMeta := make(map[string]interface{})
	json.Unmarshal(c.res.Meta, auxMeta)

	nn, err := makeMap(
		"handle", c.res.Handle,
		"name", c.res.Name,
		"meta", auxMeta,
		"labels", c.res.Labels,
	)

	if c.fields != nil && len(c.fields) > 0 {
		c.fields.configureEncoder(c.encoderConfig)

		// Currently only mapped field representation is supported
		//
		// @todo Expand
		nn, err = encodeResource(nn, "fields", c.fields, true, "name")
		if err != nil {
			return nil, err
		}
	}
	if err != nil {
		return nil, err
	}

	nn, err = encodeTimestamps(nn, c.ts)
	if err != nil {
		return nil, err
	}

	return nn, err
}

func (c *composeModuleField) MarshalYAML() (interface{}, error) {

	auxOpt := c.res.Options
	if c.res.Kind == "Record" {
		ref := c.relMod.Handle
		if ref == "" {
			ref = c.relMod.Name
		}

		auxOpt["module"] = ref
		delete(auxOpt, "moduleID")
	}

	if _, has := auxOpt["multiDelimiter"]; has {
		if auxOpt["multiDelimiter"] == "\n" {
			auxOpt["multiDelimiter"] = "\\n"
		}
	}

	nsn, err := makeMap(
		"kind", c.res.Kind,
		"name", c.res.Name,
		"label", c.res.Label,

		"options", auxOpt,

		"private", c.res.Private,
		"required", c.res.Required,
		"visible", c.res.Visible,
		"multi", c.res.Multi,
		"defaultValue", c.res.DefaultValue,

		"expressions", c.expr,
		"labels", c.res.Labels,
	)
	if err != nil {
		return nil, err
	}

	nsn, err = encodeTimestamps(nsn, c.ts)
	if err != nil {
		return nil, err
	}
	return nsn, err
}

func (fe composeModuleFieldExpr) MarshalYAML() (interface{}, error) {
	return makeMap(
		"valueExpr", fe.ValueExpr,
		"sanitizers", fe.Sanitizers,
		"validators", fe.Validators,
		"disableDefaultValidators", fe.DisableDefaultValidators,
		"formatters", fe.Formatters,
		"disableDefaultFormatters", fe.DisableDefaultFormatters,
	)
}

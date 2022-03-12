# Corteza codegen tools and definitions

## History

See [old codegen](../pkh/codegen/README.md).

## Plans

Right now, Corteza is migrating its old YAML definitions to CUE.
We are also simplifying all templates by moving as much data manipulation to Cue as possible.

### Todo
 - options documentation (see assets/templates/docs/options.gen.adoc.tpl)
 

## Intro 

Codegen tools are based on [cuelang](https://cuelang.org/docs/tutorials/) and golang templates.

**What you can find here:**
 - [codegen definitions in `codegen/`](./)
 - [templates in `codegen/asset/templates`](./asset/templates)
 - [schemas in `codegen/schema`](./schema)
 - [template exec tool in `codegen/tool`](./tool)

Platform, component and resource definitions (.cue files) can be found in:
 - `app`
 - `automation`
 - `system`
 - `compose`
 - `federation`

## Running code generator

When a definitions or templates are changed all outputs need to be regenerated

This can be done with the following command (the root of the project)
```shell
make codegen
```

Inside `/codegen` you can run codegen with different options
```shell
# generate all
make 

# generate only server files (golang code, .env.example, ...)  
make server 

# generate only docs (asciidoc files)
make docs

# generate for only one specific codegen definition
make server.options.cue
make docs.options.cue

# or more of them at once
make server.options.cue docs.options.cue
```


## How does it work?

### High-level overview

See [Makefile's `codegen` task](../Makefile)

1. evaluate codegen instructions (see [platform.cue](./platform.cue))
2. output instructions as JSON
3. pipe JSON into [template exec tool](./tool)
4. process instructions, load and execute templates and write to output files.

### Definition structure

#### Codegen instructions

Collection of `#codegen` structs with template + payload + output instructions. Template exec tool iterates over collection and creates output from each one. 
endpoints (unrelated to specific component)

# Troubleshooting

1. Make sure you have the latest version of cue (0.4.2 or newer)

# Corteza codegen tools and definitions

## History

See [old codegen](../pkh/codegen/README.md).

## Plans

Right now, Corteza is migrating its old YAML definitions to CUE.
We are also simplifying all templates by moving as much data manipulation to Cue as possible.

## Intro 

Codegen tools are based on [cuelang](https://cuelang.org/docs/tutorials/) and golang templates.

**What you can find here:**
 - [codegen definitions in `codegen/`](./)
 - [templates in `codegen/asset/templates`](./asset/templates)
 - [schemas in `codegen/schema`](./schema)
 - [template exec tool in `codegen/tool`](./tool)

Platform, component and resource definitions (.cue files) can be found in:
 - `app`
 - `automation` @todo
 - `system`
 - `compose`
 - `federation` @todo

## Running code generator

When a definitions or templates are changed all outputs need to be regenerated

This can be done with the following command
```
make codegen
```
Please note that 


## How does it work?

### High-level overview

See [Makefile's `codegen-cue` task](../Makefile)

1. evaluate codegen instructions (see [platform.cue](./platform.cue))
2. output instructions as JSON
3. pipe JSON into [template exec tool](./tool)
4. process instructions, load and execute templates and write to output files.

### Definition structure

#### Codegen instructions

Collection of `#codegen` structs with template + payload + output instructions. Template exec tool iterates over collection and creates output from each one. 

#### Platform

Main entry point that combines all components

 - @todo options 

 - @todo REST endpoints (unrelated to specific component)

#### Component

Defines component, it's behaviour, RBAC operations and resources

 - @todo REST endpoints (unrelated to specific resources)

#### Resource

Defines resource, it's behaviour, types, RBAC operations, translatable keys

 - @todo events 
 - @todo actions 
 - @todo errors 
 - @todo automation functions
 - @todo expression types 
 - @todo REST endpoints


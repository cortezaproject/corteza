This is 3rd generation of corteza code generator
(1st was written in PHP, apologies 2nd is active and in use)

Goal for the 3rd generation is:
 - improved definition structure that addresses all elements that are used for code generation (types, store, rbac)
 - unify most if not all definitions (stuff in YAML files) into what you can find in /def/...
 - JSON schema validation for YAML files

How to run it (temp):
```
go run pkg/codegen-v3/*.go def
```

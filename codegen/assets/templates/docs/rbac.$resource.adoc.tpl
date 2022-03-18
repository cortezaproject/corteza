= {{ .label }}

[cols="1s,5a,5a"]
|===
| Operation| Description | Default

{{ range .operations }}
| [#{{ .slug }}]#<<{{ .slug }},{{ .label }}>>#
| {{ .description }}
| Deny

{{ end }}
|===

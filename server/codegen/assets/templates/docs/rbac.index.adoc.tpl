= {{ .label }}

:leveloffset: +1

include::./component.gen.adoc[]
{{ range .resources }}
include::./resource.{{ . }}.gen.adoc[]
{{ end }}

:leveloffset: -1


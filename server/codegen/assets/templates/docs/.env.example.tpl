###############################################################################
###############################################################################
#
# This file is auto-generated.
#
# Copy it to .env and uncomment set or adjust any options you need.
#
# There are three ways how this options are loaded when Corteza starts:
#  1. Environmental variables are always loaded and can not be overrided
#     by alternative sources (below) (*)
#  2. If one or more --env-file flag is provided corteza looks in the provided
#     locations (either file or directory containing .env)
#  3. Directory where the corteza binary is located is scanned for .env file)
#
# (*) You can use mechanisms from docker or docker-compose to set or load
#     environmental variables as well
{{- range .groups }}

###############################################################################
###############################################################################
{{ .title }}
#
{{- if .intro }}
{{ .intro }}
#
{{ end -}}

{{ range .options }}

###############################################################################
{{- if .description }}
{{ .description }}
{{- end }}
# Type:    {{ .type }}
# Default: {{ .defaultValue }}
# {{ .env }}={{ .defaultValue }}
{{- end }}
{{- end }}

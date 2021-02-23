{{ range . }}
	<div class="alert alert-{{ .Type }}" role="alert">
		{{ .Text | html }}
	</div>
{{ end }}

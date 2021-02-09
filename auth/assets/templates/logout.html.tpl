{{ template "inc_header.html.tpl" . }}
<div class="card-body">
	<form>
		{{ if .form.error }}
		<div class="alert alert-danger" role="alert">
			{{ .form.error }}
		</div>
		{{ else }}
		<div class="alert alert-primary" role="primary">
			Logout successful
		</div>
		{{ end }}

		<hr />

		<div class="text-center my-3">
			<a href="{{ if .backlink }}{{ .backlink }}{{ else }}{{ links.Login }}{{ end }}">Login</a>
		</div>

	</form>
</div>
{{ template "inc_footer.html.tpl" . }}

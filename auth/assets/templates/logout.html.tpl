{{ template "inc_header.html.tpl" . }}
<div class="card-body">
	<form>
		{{ if .form.error }}
		<div class="text-danger font-weight-bold" role="alert">
			{{ .form.error }}
		</div>
		{{ else }}
		<div class="text-dark font-weight-bold" role="primary">
			Log out successful.
		</div>
		{{ end }}

		<div class="text-center my-3">
			Click here to <a href="{{ if .backlink }}{{ .backlink }}{{ else }}{{ links.Login }}{{ end }}">Log in</a>
		</div>

	</form>
</div>
{{ template "inc_footer.html.tpl" . }}

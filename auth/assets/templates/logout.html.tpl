{{ template "inc_header.html.tpl" . }}
<div class="card-body">
	<form>
		{{ if .form.error }}
		<div class="text-danger font-weight-bold" role="alert">
			{{ .form.error }}
		</div>
		{{ else }}
		<div class="text-dark font-weight-bold" role="primary">
			{{ tr "logout.template.log-out" }}
		</div>
		{{ end }}

		<div class="text-center my-3">
			{{ tr "logout.template.log-in" "link" .link }}
		</div>

	</form>
</div>
{{ template "inc_footer.html.tpl" . }}

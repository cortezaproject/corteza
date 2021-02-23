{{ template "inc_header.html.tpl" . }}
<div class="card-body">
	<h4 class="card-title">Internal error</h4>
	<div class="alert alert-danger" role="alert">
		{{ .error }}
	</div>
</div>
{{ template "inc_footer.html.tpl" . }}

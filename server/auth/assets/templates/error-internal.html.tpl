{{ template "inc_header.html.tpl" . }}
<div class="card-body p-0">
	<h4 class="card-title p-3 border-bottom">{{ tr "error-internal.template.title" }}</h4>
	<div class="text-danger mb-4 font-weight-bold p-3" role="alert">
		{{ .error }}
	</div>
</div>
{{ template "inc_footer.html.tpl" . }}

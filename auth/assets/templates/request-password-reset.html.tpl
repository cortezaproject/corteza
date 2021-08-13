{{ template "inc_header.html.tpl" . }}
<div class="card-body p-0">
	<h1 class="h4 card-title p-3 border-bottom">{{ tr "request_password_reset.template.title" }}</h1>
	<form
		method="POST"
		action="{{ links.RequestPasswordReset }}"
		class="p-3"
	>
		{{ .csrfField }}
		{{ if .form.error }}
		<div class="text-danger font-weight-bold mb-4" role="alert">
			{{ .form.error }}
		</div>
		{{ end }}
		<div class="mb-3">
			<label>
                {{ tr "request_password_reset.template.form.email.label" }}
            </label>
			<input
				type="email"
				class="form-control"
				name="email"
				required
				placeholder="{{ tr "request_password_reset.template.form.email.placeholder" }}"
				autocomplete="username"
				value="{{ if .form }}{{ .form.email }}{{ end }}"
				aria-label="{{ tr "request_password_reset.template.form.email.aria-label" }}">
		</div>
		<div class="text-right">
			<button class="btn btn-primary btn-block btn-lg" type="submit">{{ tr "request_password_reset.template.form.button" }}</button>
		</div>
	</form>
	<div class="text-center my-3">
		<a href="{{ links.Signup }}">{{ tr "request_password_reset.template.new-acc-link" }}</a>
		{{ tr "request_password_reset.template.or" }}
		<a href="{{ links.Login }}">{{ tr "request_password_reset.template.login" }}</a>
	</div>
</div>
{{ template "inc_footer.html.tpl" . }}

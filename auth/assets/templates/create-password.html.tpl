{{ template "inc_header.html.tpl" set . "activeNav" "security" }}
<div class="card-body p-0">
	<h4 class="card-title p-3 border-bottom">{{ tr "create-password.template.title" }}</h4>
	<form
		method="POST"
		action="{{ links.CreatePassword }}"
		class="p-3"
	>
		{{ .csrfField }}
		{{ if .form.error }}
		<div class="text-danger font-weight-bold mb-3" role="alert">
			{{ .form.error }}
		</div>
		{{ end }}
		<div class="mb-3">
		    <label>
                {{ tr "create-password.template.form.email.label" }}
            </label>
			<input
				type="email"
				class="form-control"
				name="email"
				readonly
				value="{{ .user.Email }}"
				aria-label="{{ tr "create-password.template.form.email.label" }}">
		</div>
		<div class="mb-3">
            <label>
                {{ tr "create-password.template.form.password.label" }}
            </label>
			<input
				type="password"
				required
				class="form-control"
				name="password"
				autocomplete="password"
				placeholder="{{ tr "create-password.template.form.password.placeholder" }}"
				aria-label="{{ tr "create-password.template.form.password.label" }}">
		</div>
		<div class="text-right">
			<button class="btn btn-primary btn-block btn-lg" type="submit">{{ tr "create-password.template.form.button.create-password" }}</button>
		</div>
	</form>
</div>
{{ template "inc_footer.html.tpl" . }}

{{ template "inc_header.html.tpl" . }}
<div class="card-body p-0">
	<h1 class="h4 card-title p-3 border-bottom">{{ tr "reset_password.template.title" }}</h1>
	<form
		method="POST"
		action="{{ links.ResetPassword }}"
		class="p-3"
	>
		{{ .csrfField }}
		{{ if .form.error }}
		<div class="text-danger font-weight-bold p-3" role="alert">
			{{ .form.error }}
		</div>
		{{ end }}
		<div class="mb-3">
			<label>
                {{ tr "reset_password.template.form.email.label" }}
            </label>
			<input
				type="email"
				class="form-control"
				name="email"
				readonly
				placeholder="{{ tr "reset_password.template.form.email.placeholder" }}"
				value="{{ .user.Email }}"
				aria-label="{{ tr "reset_password.template.form.email.aria-label" }}">
		</div>
		<div class="mb-3">
            <label>
                {{ tr "reset_password.template.form.new-password.label" }}
            </label>
			<input
				type="password"
				required
				class="form-control"
				name="password"
				autocomplete="new-password"
				placeholder="{{ tr "reset_password.template.form.new-password.placeholder" }}"
				aria-label="{{ tr "reset_password.template.form.new-password.aria-label" }}">
		</div>
		<div class="text-right">
			<button class="btn btn-primary btn-block btn-lg" type="submit">{{ tr "reset_password.template.form.button.change-password" }}</button>
		</div>
	</form>
</div>
{{ template "inc_footer.html.tpl" . }}

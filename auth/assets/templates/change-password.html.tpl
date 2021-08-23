{{ template "inc_header.html.tpl" set . "activeNav" "security" }}
<div class="card-body p-0">
	<h1 class="h4 card-title p-3 border-bottom">{{ tr "change-password.template.title" }}</h1>
	<form
		method="POST"
		action="{{ links.ChangePassword }}"
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
                {{ tr "change-password.template.form.email.label" }}
            </label>
			<input
				type="email"
				class="form-control"
				name="email"
				readonly
				placeholder="{{ tr "change-password.template.form.email.placeholder" }}"
				value="{{ .user.Email }}"
				aria-label="{{ tr "change-password.template.form.email.aria-label" }}">
		</div>
		<div class="mb-3">
            <label>
                {{ tr "change-password.template.form.old-password.label" }}
            </label>
			<input
				type="password"
				required
				class="form-control"
				name="oldPassword"
				autocomplete="current-password"
				placeholder="{{ tr "change-password.template.form.old-password.placeholder" }}"
				aria-label="{{ tr "change-password.template.form.old-password.aria-label" }}">
		</div>
		<div class="mb-3">
            <label>
                {{ tr "change-password.template.form.new-password.label" }}
            </label>
			<input
				type="password"
				required
				class="form-control"
				name="newPassword"
				autocomplete="new-password"
				placeholder="{{ tr "change-password.template.form.new-password.placeholder" }}"
				aria-label="{{ tr "change-password.template.form.new-password.aria-label" }}">
		</div>
		<div class="text-right">
			<button class="btn btn-primary btn-block btn-lg" type="submit">{{ tr "change-password.template.form.button.change-password" }}</button>
		</div>
	</form>
</div>
{{ template "inc_footer.html.tpl" . }}

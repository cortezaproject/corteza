{{ template "inc_header.html.tpl" . }}
<div class="card-body p-0">
	<h4 class="card-title p-3 border-bottom">Request password reset link</h4>
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
                E-mail *
            </label>
			<input
				type="email"
				class="form-control"
				name="email"
				required
				placeholder="email@domain.ltd"
				autocomplete="username"
				value="{{ if .form }}{{ .form.email }}{{ end }}"
				aria-label="Email">
		</div>
		<div class="text-right">
			<button class="btn btn-primary btn-block btn-lg" type="submit">Request password reset link via email</button>
		</div>
	</form>
	<div class="text-center my-3">
		<a href="{{ links.Signup }}">Create new account</a>
		or
		<a href="{{ links.Login }}">Log in</a>
	</div>
</div>
{{ template "inc_footer.html.tpl" . }}

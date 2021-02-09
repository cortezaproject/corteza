{{ template "inc_header.html.tpl" . }}
<div class="card-body">
	{{ template "inc_alerts.html.tpl" .alerts }}
	<h4 class="card-title">Request password reset link</h4>
	<form
		method="POST"
		onsubmit="buttonDisabler()"
		action="{{ links.RequestPasswordReset }}"
	>
		{{ .csrfField }}
		{{ if .form.error }}
		<div class="alert alert-danger" role="alert">
			{{ .form.error }}
		</div>
		{{ end }}
		<div class="input-group mb-3">
			<span class="input-group-text">
			  <i class="bi bi-envelope"></i>
			</span>
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
			<button class="btn btn-primary btn-block" type="submit">Request password reset link via email</button>
		</div>
	</form>
	<div class="text-center my-3">
		<a href="{{ links.Signup }}">Create new account</a>
		|
		<a href="{{ links.Login }}">Login</a>
	</div>
</div>
{{ template "inc_footer.html.tpl" . }}

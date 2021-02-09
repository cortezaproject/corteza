{{ template "inc_header.html.tpl" . }}
<div class="card-body">
	{{ template "inc_alerts.html.tpl" .alerts }}
	<h4 class="card-title">Login</h4>
	{{ if .settings.LocalEnabled }}
	<form
		method="POST"
		onsubmit="buttonDisabler()"
		action="{{ links.Login }}"
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
				value="{{ if .form }}{{ .form.email }}{{ end }}"
				autocomplete="username"
				aria-label="Email">
		</div>
		<div class="input-group mb-3">
			<span class="input-group-text">
			  <i class="bi bi-key-fill"></i>
			</span>
			<input
				type="password"
				required
				class="form-control"
				name="password"
				placeholder="Password"
				autocomplete="current-password"
				aria-label="Password">
		</div>
		<div class="row">
			<div class="col">
			{{ if .settings.PasswordResetEnabled }}
				<a href="{{ links.RequestPasswordReset }}" class="small">Forgot your password?</a>
			{{ end }}
			</div>
			<div class="col text-right">
				<button
					class="btn btn-primary"
					name="keep-session"
					value="true"
					type="submit"
				>
					Log in and remember me
				</button>
				<button
					class="btn btn-secondary"
					type="submit"
				>
					Log in
				</button>
			</div>
		</div>
	</form>
	{{ if .settings.SignupEnabled }}
	<div class="text-center my-5">
		<a href="{{ links.Signup }}">Create new account</a>
	</div>
	{{ end }}
	{{ end }}
	{{ if .settings.ExternalEnabled }}
	<hr>
	<div>
	{{ range .providers }}
		<a href="{{ links.External }}/{{ .Handle }}" class="btn btn-outline-dark btn-block text-left mb-2">
			<i class="bi bi-{{ .Icon }} mr-2"></i>
			<small>Login with {{ coalesce .Label .Handle }}</small>
		</a>
	{{ end }}
	</div>
	{{ end }}
</div>
{{ template "inc_footer.html.tpl" . }}

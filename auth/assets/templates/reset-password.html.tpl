{{ template "inc_header.html.tpl" . }}
<div class="card-body">
	{{ template "inc_alerts.html.tpl" .alerts }}
	<h4 class="card-title">Reset your password</h4>
	<form
		method="POST"
		onsubmit="buttonDisabler()"
		action="{{ links.ResetPassword }}"
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
				readonly
				placeholder="email@domain.ltd"
				value="{{ .user.Email }}"
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
				autocomplete="new-password"
				placeholder="Set new password"
				aria-label="Set new password">
		</div>
		<div class="text-right">
			<button class="btn btn-primary btn-block" type="submit">Change your password</button>
		</div>
	</form>
</div>
{{ template "inc_footer.html.tpl" . }}

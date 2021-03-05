{{ template "inc_header.html.tpl" . }}
<div class="card-body p-0">
	{{ template "inc_alerts.html.tpl" .alerts }}
	<h4 class="card-title p-3 border-bottom">Change your password</h4>
	<form
		method="POST"
		action="{{ links.ChangePassword }}"
		onsubmit="buttonDisabler()"
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
                E-mail *
            </label>
			<input
				type="email"
				class="form-control"
				name="email"
				readonly
				placeholder="email@domain.ltd"
				value="{{ .user.Email }}"
				aria-label="Email">
		</div>
		<div class="mb-3">
            <label>
                Old Password *
            </label>
			<input
				type="password"
				required
				class="form-control"
				name="oldPassword"
				autocomplete="current-password"
				placeholder="Enter your old password"
				aria-label="Old password">
		</div>
		<div class="mb-3">
            <label>
                New Password *
            </label>
			<input
				type="password"
				required
				class="form-control"
				name="newPassword"
				autocomplete="new-password"
				placeholder="Enter your new password"
				aria-label="New password">
		</div>
		<div class="text-right">
			<button class="btn btn-primary btn-block btn-lg" type="submit">Change your password</button>
		</div>
	</form>
</div>
{{ template "inc_footer.html.tpl" . }}

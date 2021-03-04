{{ template "inc_header.html.tpl" set . "activeNav" "security" }}
<div class="card-body p-0">
	<h4 class="card-title p-3 border-bottom">Security</h4>
	<form
		method="POST"
		action="{{ links.Security }}"
		class="p-3"
	>

	{{ if .settings.LocalEnabled }}
	<h5>Password</h5>
	<a href="{{ links.ChangePassword }}">Change your password</a>
	{{ end }}

	<hr />

	<div>
		{{ .csrfField }}
		<h5>Multi-factor authentication</h5>
		{{ if or .settings.MultiFactor.TOTP.Enabled .settings.MultiFactor.EmailOTP.Enabled }}
			{{ if .settings.MultiFactor.TOTP.Enabled }}
			<div class="py-4">
				<h6>Additional security with mobile app (time-based one-time-password)</h6>
				<div class="row">
					<div class="col-1 text-left">
						{{ if .totpEnforced }}
						<h4><i class="bi bi-check-square text-primary"></i></h4>
						{{ else }}
						<h4><i class="bi bi-exclamation-triangle-fill text-danger"></i></h4>
						{{ end }}
					</div>
					<div class="col-7 pt-2">
						{{ if .totpEnforced }}
							Configured and required on login.
						{{ else }}
							Currently disabled.
						{{ end }}
					</div>
					<div class="col-5">
						{{ if .totpEnforced }}
							{{ if not .settings.MultiFactor.TOTP.Enforced }}
							<button name="action" value="disableTOTP" class="btn btn-link text-danger">Disable</button>
							{{ end }}
						{{ else }}
							<button name="action" value="configureTOTP" class="btn btn-link text-primary">Configure</button>
						{{ end }}
					</div>

				</div>
			</div>
			{{ end }}

			{{ if .settings.MultiFactor.EmailOTP.Enabled }}
			<div class="py-4">
				<h6>Additional security with one-time-password over email</h6>
				<div class="row">
					<div class="col-1 text-left">
						{{ if .emailOtpEnforced }}
						<h4><i class="bi bi-check-square text-primary"></i></h4>
						{{ else }}
						<h4><i class="bi bi-exclamation-triangle-fill text-danger"></i></h4>
						{{ end }}
					</div>
					<div class="col-7 pt-2">
					{{ if .emailOtpEnforced }}
						Enabled and required on login.
					{{ else }}
						Currently disabled.
					{{ end }}
					</div>
					<div class="col-4">
					{{ if .emailOtpEnforced }}
						{{ if not .settings.MultiFactor.EmailOTP.Enforced }}
						<button name="action" value="disableEmailOTP" class="btn btn-link text-danger">Disable</button>
						{{ end }}
					{{ else }}
						<button name="action" value="enableEmailOTP" class="btn btn-link text-primary">Enable</button>
					{{ end }}
					</div>

				</div>
			</div>
			{{ end }}
		{{ else }}
			<div class="text-danger font-weight-bold mb-4" role="alert">
				All MFA methods are currently disabled. Ask your administrator to enable them.
			</div>
		{{ end }}
		</div>
	</form>
</div>
{{ template "inc_footer.html.tpl" . }}

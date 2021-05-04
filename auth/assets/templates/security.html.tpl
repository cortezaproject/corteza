{{ template "inc_header.html.tpl" set . "activeNav" "security" }}
<div class="card-body p-0">
	<h1 class="h4 card-title p-3 border-bottom">Security</h1>
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
					<div class="col-10 pt-2">
						{{ if .totpEnforced }}
						<i class="bi bi-check-circle text-success h5 mr-1"></i>
						{{ else }}
						<i class="bi bi-exclamation-circle-fill text-danger h5 mr-1"></i>
						{{ end }}
						{{ if .totpEnforced }}
							Configured and required on login.
						{{ else }}
							Currently disabled.
						{{ end }}
					</div>
					<div class="col-md-2 col-sm-12">
						{{ if .totpEnforced }}
							{{ if not .settings.MultiFactor.TOTP.Enforced }}
							<button name="action" value="disableTOTP" class="btn btn-danger float-right">Disable</button>
							{{ end }}
						{{ else }}
							<button name="action" value="configureTOTP" class="btn btn-primary float-right">Configure</button>
						{{ end }}
					</div>
				</div>
			</div>
			{{ end }}

			{{ if .settings.MultiFactor.EmailOTP.Enabled }}
			<div class="py-4">
				<h6>Additional security with one-time-password over email</h6>
				<div class="row">
					<div class="col-10 pt-2">
                    {{ if .emailOtpEnforced }}
                    <i class="bi bi-check-circle text-success h5 mr-1"></i>
                    {{ else }}
                    <i class="bi bi-exclamation-circle-fill text-danger h5 mr-1"></i>
                    {{ end }}
					{{ if .emailOtpEnforced }}
						Enabled and required on login.
					{{ else }}
						Currently disabled.
					{{ end }}
					</div>
					<div class="col-md-2 col-sm-12">
					{{ if .emailOtpEnforced }}
						{{ if not .settings.MultiFactor.EmailOTP.Enforced }}
						<button name="action" value="disableEmailOTP" class="btn btn-danger float-right">Disable</button>
						{{ end }}
					{{ else }}
						<button name="action" value="enableEmailOTP" class="btn btn-primary float-right">Enable</button>
					{{ end }}
					</div>

				</div>
			</div>
			{{ end }}
		{{ else }}
			<div class="mb-4 font-italic" role="alert">
				All MFA methods are currently disabled. Ask your administrator to enable them.
			</div>
		{{ end }}
		</div>
	</form>
</div>
{{ template "inc_footer.html.tpl" . }}

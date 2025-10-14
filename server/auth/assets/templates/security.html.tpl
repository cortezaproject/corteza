{{ template "inc_header.html.tpl" set . "activeNav" "security" }}
<div class="card-body p-0">
	<form
		method="POST"
		action="{{ links.Security }}"
		class="p-3"
	>

	{{ if .settings.LocalEnabled }}
	<h5>{{ tr "security.template.password.title" }}</h5>
	<a
		data-test-id="link-change-password"
		href="{{ links.ChangePassword }}"
		>
			{{ tr "security.template.password.change-link" }}
		</a>
	{{ end }}

	<hr />

	<div>
		{{ .csrfField }}
		<h5 class="mb-3">{{ tr "security.template.mfa.title" }}</h5>
		{{ if or .settings.MultiFactor.TOTP.Enabled .settings.MultiFactor.EmailOTP.Enabled }}
			{{ if .settings.MultiFactor.TOTP.Enabled }}
			<div class="mb-4">
				<label class="text-primary">{{ tr "security.template.mfa.totp.title" }}</label>
				<div class="d-flex align-items-center">
					<div>
						{{ if .totpEnforced }}
						  <i class="bi bi-check-circle text-success h5 mr-1"></i>
						{{ else }}
						  <i class="bi bi-exclamation-circle-fill text-danger h5 mr-1"></i>
						{{ end }}
						{{ if .totpEnforced }}
							{{ tr "security.template.mfa.totp.enforced" }}
						{{ else }}
							{{ tr "security.template.mfa.totp.disabled" }}
						{{ end }}
					</div>

					<div class="ml-auto">
						{{ if .totpEnforced }}
							{{ if not .settings.MultiFactor.TOTP.Enforced }}
                <button
                  data-test-id="button-disable-totp"
                  name="action"
                  value="disableTOTP"
                  class="btn btn-danger"
                >
                  {{ tr "security.template.mfa.totp.disable" }}
                </button>
							{{ end }}
						{{ else }}
							<button
								data-test-id="button-configure-totp"
								name="action"
								value="configureTOTP"
								class="btn btn-primary"
							>
								{{ tr "security.template.mfa.totp.configure" }}
							</button>
						{{ end }}
					</div>
				</div>
			</div>
			{{ end }}

			{{ if .settings.MultiFactor.EmailOTP.Enabled }}
			<div class="mb-3">
				<label class="text-primary">{{ tr "security.template.mfa.email.title" }}</label class="text-primary">
				<div class="d-flex align-items-center">
					<div>
            {{ if .emailOtpEnforced }}
              <i class="bi bi-check-circle text-success h5 mr-1"></i>
            {{ else }}
              <i class="bi bi-exclamation-circle-fill text-danger h5 mr-1"></i>
            {{ end }}

            {{ if .emailOtpEnforced }}
              {{ tr "security.template.mfa.email.enforced" }}
            {{ else }}
              {{ tr "security.template.mfa.email.disabled" }}
            {{ end }}
					</div>

					<div class="ml-auto">
					{{ if .emailOtpEnforced }}
						{{ if not .settings.MultiFactor.EmailOTP.Enforced }}
						<button
							data-test-id="button-disable-email-otp"
							name="action"
							value="disableEmailOTP"
							class="btn btn-danger"
						>
							{{ tr "security.template.mfa.email.disable" }}
						</button>
						{{ end }}
					{{ else }}
						<button
							data-test-id="button-enable-email-otp"
							name="action"
							value="enableEmailOTP"
							class="btn btn-primary"
						>
							{{ tr "security.template.mfa.email.enable" }}
						</button>
					{{ end }}
					</div>

				</div>
			</div>
			{{ end }}
		{{ else }}
			<div class="mb-1 font-italic" role="alert">
				{{ tr "security.template.mfa.all-disabled" }}
			</div>
		{{ end }}
		</div>
	</form>
</div>
{{ template "inc_footer.html.tpl" . }}

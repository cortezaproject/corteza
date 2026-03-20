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
						  <svg class="text-success mr-1" xmlns="http://www.w3.org/2000/svg" width="20" height="20" fill="currentColor" viewBox="0 0 16 16"><path d="M8 15A7 7 0 1 1 8 1a7 7 0 0 1 0 14m0 1A8 8 0 1 0 8 0a8 8 0 0 0 0 16"/><path d="m10.97 4.97-.02.022-3.473 4.425-2.093-2.094a.75.75 0 0 0-1.06 1.06L6.97 11.03a.75.75 0 0 0 1.079-.02l3.992-4.99a.75.75 0 0 0-1.071-1.05"/></svg>
						{{ else }}
						  <svg class="text-danger mr-1" xmlns="http://www.w3.org/2000/svg" width="20" height="20" fill="currentColor" viewBox="0 0 16 16"><path d="M16 8A8 8 0 1 1 0 8a8 8 0 0 1 16 0M8 4a.905.905 0 0 0-.9.995l.35 3.507a.552.552 0 0 0 1.1 0l.35-3.507A.905.905 0 0 0 8 4m.002 6a1 1 0 1 0 0 2 1 1 0 0 0 0-2"/></svg>
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
              <svg class="text-success mr-1" xmlns="http://www.w3.org/2000/svg" width="20" height="20" fill="currentColor" viewBox="0 0 16 16"><path d="M8 15A7 7 0 1 1 8 1a7 7 0 0 1 0 14m0 1A8 8 0 1 0 8 0a8 8 0 0 0 0 16"/><path d="m10.97 4.97-.02.022-3.473 4.425-2.093-2.094a.75.75 0 0 0-1.06 1.06L6.97 11.03a.75.75 0 0 0 1.079-.02l3.992-4.99a.75.75 0 0 0-1.071-1.05"/></svg>
            {{ else }}
              <svg class="text-danger mr-1" xmlns="http://www.w3.org/2000/svg" width="20" height="20" fill="currentColor" viewBox="0 0 16 16"><path d="M16 8A8 8 0 1 1 0 8a8 8 0 0 1 16 0M8 4a.905.905 0 0 0-.9.995l.35 3.507a.552.552 0 0 0 1.1 0l.35-3.507A.905.905 0 0 0 8 4m.002 6a1 1 0 1 0 0 2 1 1 0 0 0 0-2"/></svg>
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

{{ template "inc_header.html.tpl" set . "activeNav" "security" }}
<div class="card-body p-0">
	<h1 class="h4 card-title p-3 border-bottom">{{ tr "security.template.title" }}</h1>
	<form
		method="POST"
		action="{{ links.Security }}"
		class="p-3"
	>

	{{ if .settings.LocalEnabled }}
	<h5>{{ tr "security.template.form.title" }}</h5>
	<a href="{{ links.ChangePassword }}">{{ tr "security.template.form.passw-change" }}</a>
	{{ end }}

	<hr />

	<div>
		{{ .csrfField }}
		<h5>{{ tr "security.template.security.title" }}</h5>
		{{ if or .settings.MultiFactor.TOTP.Enabled .settings.MultiFactor.EmailOTP.Enabled }}
			{{ if .settings.MultiFactor.TOTP.Enabled }}
			<div class="py-4">
				<h6>{{ tr "security.template.security.description" }}</h6>
				<div class="row">
					<div class="col-10 pt-2">
						{{ if .totpEnforced }}
						<i class="bi bi-check-circle text-success h5 mr-1"></i>
						{{ else }}
						<i class="bi bi-exclamation-circle-fill text-danger h5 mr-1"></i>
						{{ end }}
						{{ if .totpEnforced }}
							{{ tr "security.template.security.req-login" }}
						{{ else }}
							{{ tr "security.template.security.disable" }}
						{{ end }}
					</div>
					<div class="col-md-2 col-sm-12">
						{{ if .totpEnforced }}
							{{ if not .settings.MultiFactor.TOTP.Enforced }}
							<button name="action" value="disableTOTP" class="btn btn-danger float-right">{{ tr "security.template.button.disable" }}</button>
							{{ end }}
						{{ else }}
							<button name="action" value="configureTOTP" class="btn btn-primary float-right">{{ tr "security.template.button.configure" }}</button>
						{{ end }}
					</div>
				</div>
			</div>
			{{ end }}

			{{ if .settings.MultiFactor.EmailOTP.Enabled }}
			<div class="py-4">
				<h6>{{ tr "security.template.aditional-security.title" }}</h6>
				<div class="row">
					<div class="col-10 pt-2">
                    {{ if .emailOtpEnforced }}
                    <i class="bi bi-check-circle text-success h5 mr-1"></i>
                    {{ else }}
                    <i class="bi bi-exclamation-circle-fill text-danger h5 mr-1"></i>
                    {{ end }}
					{{ if .emailOtpEnforced }}
						{{ tr "security.template.aditional-security.email-otp-enforced" }}
					{{ else }}
						{{ tr "security.template.aditional-security.disable" }}
					{{ end }}
					</div>
					<div class="col-md-2 col-sm-12">
					{{ if .emailOtpEnforced }}
						{{ if not .settings.MultiFactor.EmailOTP.Enforced }}
						<button name="action" value="disableEmailOTP" class="btn btn-danger float-right">{{ tr "security.template.button.disable" }}</button>
						{{ end }}
					{{ else }}
						<button name="action" value="enableEmailOTP" class="btn btn-primary float-right">{{ tr "security.template.button.enable" }}</button>
					{{ end }}
					</div>

				</div>
			</div>
			{{ end }}
		{{ else }}
			<div class="mb-4 font-italic" role="alert">
				{{ tr "security.template.alert" }}
			</div>
		{{ end }}
		</div>
	</form>
</div>
{{ template "inc_footer.html.tpl" . }}

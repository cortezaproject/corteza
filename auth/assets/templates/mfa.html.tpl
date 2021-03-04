{{ template "inc_header.html.tpl" set . "hideNav" true }}
<div class="card-body p-0">
	<h4 class="card-title p-3 border-bottom">Multi-factor authentication</h4>

	{{ if .emailOtpPending }}
	<form
		class="p-3"
		method="POST"
		action="{{ links.Mfa }}"
	>
		<h5>Check your inbox and enter the received code</h5>

		{{ if .form.emailOtpError }}
		<div class="alert alert-danger" role="alert">
			{{ .form.emailOtpError }}
		</div>
		{{ end }}
		{{ .csrfField }}


		<div class="input-group my-3">
			<input
				type="text"
				required
				class="form-control text-center mfa-code-mask"
				name="code"
				maxlength="6"
				minlength="6"
				aria-required="true"
				placeholder="000 000"
				autocomplete="off"
				aria-label="Code">
		</div>

		<button
			class="btn btn-primary btn-block btn-lg"
			name="action"
			value="verifyEmailOtp"
			type="submit"
		>
			Verify
		</button>

		<a
			href="{{ links.Mfa }}?action=resendEmailOtp"
			class="btn btn-link btn-block btn-lg"
			name="action"
			value="resendEmailOtp"
		>
			Resend
		</a>
	</form>
	{{ else if not .emailOtpDisabled }}
		<p class="p-3">
			<i class="bi bi-check text-primary"></i> Email OTP confirmed
		</p>
	{{ end }}

	{{ if .totpPending }}
	<form
		class="p-3"
		method="POST"
		action="{{ links.Mfa }}"
	>
		<h5>Check your TOTP application and enter the received code</h5>

		{{ if .form.totpError }}
		<div class="alert alert-danger" role="alert">
			{{ .form.totpError }}
		</div>
		{{ end }}
		{{ .csrfField }}


		<div class="input-group my-3">
			<input
				type="text"
				required
				class="form-control text-center mfa-code-mask"
				name="code"
				maxlength="6"
				minlength="6"
				aria-required="true"
				placeholder="000 000"
				autocomplete="off"
				aria-label="Code">
		</div>

		<button
			class="btn btn-primary btn-block btn-lg"
			type="submit"
			name="action"
			value="verifyTotp"
		>
			Verify
		</button>
	</form>
	{{ else if not .totpDisabled }}
		<p class="p-3">
			<i class="bi bi-check text-primary"></i> TOTP confirmed
		</p>
	{{ end }}
</div>
{{ template "inc_footer.html.tpl" . }}

{{ template "inc_header.html.tpl" set . "hideNav" true }}
<div class="card-body p-0 mb-2">
	<h4 class="card-title p-3 border-bottom">{{ tr "mfa.template.title" }}</h4>

	{{ if .emailOtpPending }}
	<form
		class="p-3"
		method="POST"
		action="{{ links.Mfa }}"
	>
		<h5>{{ tr "mfa.template.email.instructions" }}</h5>

		{{ if .form.emailOtpError }}
		<div class="text-danger my-4 font-weight-bold" role="alert">
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
				style="letter-spacing:5px;font-size:20px;"
				aria-label="{{ tr "mfa.template.email.code" }}">
		</div>

		<button
			class="btn btn-primary btn-block btn-lg"
			name="action"
			value="verifyEmailOtp"
			type="submit"
		>
			{{ tr "mfa.template.email.verify" }}
		</button>

		<a
			href="{{ links.Mfa }}?action=resendEmailOtp"
			class="btn btn-light btn-block btn-lg text-dark"
			name="action"
			value="resendEmailOtp"
		>
			{{ tr "mfa.template.email.resend" }}
		</a>
	</form>
	{{ else if not .emailOtpDisabled }}
		<p class="p-3 mb-0">
			<i class="bi bi-check-circle text-success h5 mr-1"></i> {{ tr "mfa.template.email.confirmed" }}
		</p>
	{{ end }}

	{{ if .totpPending }}
	<form
		class="p-3"
		method="POST"
		action="{{ links.Mfa }}"
	>
		<h5>{{ tr "mfa.template.totp.instructions" }}</h5>

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
				style="letter-spacing:5px;font-size:20px;"
				aria-label="{{ tr "mfa.template.totp.code" }}">
		</div>

		<button
			class="btn btn-primary btn-block btn-lg"
			type="submit"
			name="action"
			value="verifyTotp"
		>
			{{ tr "mfa.template.email.verify" }}
		</button>
	</form>
	{{ else if not .totpDisabled }}
		<p class="px-3 pt-3 pb-2 mb-0">
			<i class="bi bi-check-circle text-success h5 mr-1"></i> {{ tr "mfa.template.totp.confirmed" }}
		</p>
	{{ end }}
</div>
{{ template "inc_footer.html.tpl" . }}

{{ template "inc_header.html.tpl" set . "hideNav" true }}
<div class="card-body p-0 mb-2">
	<h4 class="mb-0 p-3 border-bottom">{{ tr "mfa.template.title" }}</h4>

	{{ if .emailOtpPending }}
	<form
		class="p-3"
		method="POST"
		action="{{ links.Mfa }}"
	>
		<p>{{ tr "mfa.template.email.instructions" }}</p>

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
			<svg class="text-success mr-1" xmlns="http://www.w3.org/2000/svg" width="20" height="20" fill="currentColor" viewBox="0 0 16 16"><path d="M8 15A7 7 0 1 1 8 1a7 7 0 0 1 0 14m0 1A8 8 0 1 0 8 0a8 8 0 0 0 0 16"/><path d="m10.97 4.97-.02.022-3.473 4.425-2.093-2.094a.75.75 0 0 0-1.06 1.06L6.97 11.03a.75.75 0 0 0 1.079-.02l3.992-4.99a.75.75 0 0 0-1.071-1.05"/></svg> {{ tr "mfa.template.email.confirmed" }}
		</p>
	{{ end }}

	{{ if .totpPending }}
	<form
		class="p-3"
		method="POST"
		action="{{ links.Mfa }}"
	>
		<p>{{ tr "mfa.template.totp.instructions" }}</p>

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
			{{ tr "mfa.template.totp.verify" }}
		</button>
	</form>
	{{ else if and (not .totpDisabled) (not .totpPending) (not .totpUnconfigured) }}
		<p class="px-3 pt-3 pb-2 mb-0">
			<svg class="text-success mr-1" xmlns="http://www.w3.org/2000/svg" width="20" height="20" fill="currentColor" viewBox="0 0 16 16"><path d="M8 15A7 7 0 1 1 8 1a7 7 0 0 1 0 14m0 1A8 8 0 1 0 8 0a8 8 0 0 0 0 16"/><path d="m10.97 4.97-.02.022-3.473 4.425-2.093-2.094a.75.75 0 0 0-1.06 1.06L6.97 11.03a.75.75 0 0 0 1.079-.02l3.992-4.99a.75.75 0 0 0-1.071-1.05"/></svg> {{ tr "mfa.template.totp.confirmed" }}
		</p>
	{{ end }}
</div>
{{ template "inc_footer.html.tpl" . }}

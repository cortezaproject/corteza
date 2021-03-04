{{ template "inc_header.html.tpl"  set . "hideNav" true }}
<div class="card-body p-0">
	<h4 class="card-title p-3 border-bottom">Disable two-factor authentication with TOTP</h4>

	<form
		class="p-3"
		method="POST"
		action="{{ links.MfaTotpDisable }}"
	>
		Disable by entering existing code.

		{{ if .form.error }}
		<div class="alert alert-danger" role="alert">
			{{ .form.error }}
		</div>
		{{ end }}

		{{ .csrfField }}
		<div class="input-group my-3">
			<input
				type="text"
				required
				class="form-control lg text-center mfa-code-mask"
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
			name="keep-session"
			value="true"
			type="submit"
		>
			Remove
		</button>
	</form>
</div>
{{ template "inc_footer.html.tpl" . }}

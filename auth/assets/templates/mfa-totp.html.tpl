{{ template "inc_header.html.tpl"  set . "hideNav" true }}
<div class="card-body p-0">
	<h1 class="h4 card-title p-3 border-bottom">Configure two-factor authentication with TOTP</h1>

	{{ if .enforced }}
	<p class="p-3 text-danger mb-0 font-weight-bold">
		TOTP multi factor authentication is enforced by Corteza administrator.
		Please configure it right away.
	</p>
	{{ end }}


	<div class="container p-3 m-0">
		<div class="row">
			<div class="col-12 col-sm-6 p-0 mb-3">
				<pre class="h5 px-4">{{ .secret }}</pre>
				<img style="width: 280px" class="d-block m-auto pb-2" src="{{ if .devQRImage }}{{ .devQRImage }}{{ else }}{{ links.MfaTotpQRImage }}{{ end }}" />

				<form
					class="px-3"
					method="POST"
								action="{{ links.MfaTotpNewSecret }}"
				>
					Complete the configuration by entering code from the authenticator application:

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
							style="letter-spacing:5px;font-size:20px;"
							aria-label="Code">
					</div>

					<button
						class="btn btn-primary btn-block btn-lg"
						name="keep-session"
						value="true"
						type="submit"
					>
						Submit
					</button>
				</form>
			</div>
			<div class="col-12 col-sm-6">
				<p class="text-justify">
					Corteza uses time based one time passwords (TOTP) as one of the
					underlying technologies for two-factor authentication.
					Use one of the applications listed below and type in the secret or scan the QR code.
				</p>
				<p>
					This will enable additional security for your account.
				</p>
				<p>
					You can use one of the following applications:
				</p>
				<ul>
					<li>
						<a target="_blank" href="https://lastpass.com/auth/">LastPass Authenticator</a>
					</li>
					<li>
						Google Authenticator in <br />
						<a target="_blank" href="https://play.google.com/store/apps/details?id=com.google.android.apps.authenticator2">Android</a>
						or
						<a target="_blank" href="https://apps.apple.com/us/app/google-authenticator/id388497605">App Store</a>
					</li>
					<li>
						<a target="_blank" href="https://authy.com">Authy</a>
					</li>
				</ul>
			</div>
		</div>
	</div>


</div>
{{ template "inc_footer.html.tpl" . }}

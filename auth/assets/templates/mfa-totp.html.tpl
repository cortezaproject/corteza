{{ template "inc_header.html.tpl"  set . "hideNav" true }}
<div class="card-body p-0">
	<h1 class="h4 card-title p-3 border-bottom">Configure two-factor authentication with TOTP</h1>

	{{ if .enforced }}
	<p class="p-3 text-danger mb-0 font-weight-bold">
		{{ tr "mfa_totp.template.paragraph" }}
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
					{{ tr "mfa_totp.template.form.title-2" }}

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
						{{ tr "mfa_totp.template.form.button" }}
					</button>
				</form>
			</div>
			<div class="col-12 col-sm-6">
				<p class="text-justify">
					{{ tr "mfa_totp.template.form.paragraph-one" }}
				</p>
				<p>
					{{ tr "mfa_totp.template.form.paragraph-two" }}
				</p>
				<p>
					{{ tr "mfa_totp.template.form.paragraph-three" }}
				</p>
				<ul>
					<li>
						<a target="_blank" href="https://lastpass.com/auth/">{{ tr "mfa_totp.template.form.ul-links.auth" }}</a>
					</li>
					<li>
						{{ tr "mfa_totp.template.form.ul-links.text" }} <br />
						<a target="_blank" href="https://play.google.com/store/apps/details?id=com.google.android.apps.authenticator2">{{ tr "mfa_totp.template.form.ul-links.android" }}</a>
						{{ tr "mfa_totp.template.form.ul-links.or" }}
						<a target="_blank" href="https://apps.apple.com/us/app/google-authenticator/id388497605">{{ tr "mfa_totp.template.form.ul-links.store" }}</a>
					</li>
					<li>
						<a target="_blank" href="https://authy.com">{{ tr "mfa_totp.template.form.ul-links.authy" }}</a>
					</li>
				</ul>
			</div>
		</div>
	</div>


</div>
{{ template "inc_footer.html.tpl" . }}

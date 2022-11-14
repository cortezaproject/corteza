{{ template "inc_header.html.tpl"  set . "hideNav" true }}
<div class="card-body p-0">
	<h4 class="card-title p-3 border-bottom">{{ tr "mfa-totp.template.title" }}</h4>

	{{ if .enforced }}
	<p class="p-3 text-danger mb-0 font-weight-bold">
		{{ tr "mfa-totp.template.enforced" }}
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
					{{ tr "mfa-totp.template.form.title" }}

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
						{{ tr "mfa-totp.template.form.button" }}
					</button>
				</form>
			</div>
			<div class="col-12 col-sm-6">
				<p class="text-justify">
				{{ tr "mfa-totp.template.instructions" }}
				</p>
				<ul>
					<li>
						{{ tr "mfa-totp.template.lastpass" "lastpass" "https://lastpass.com/auth/" }}
					</li>
					<li>
						{{ tr "mfa-totp.template.gauth" "android" "https://play.google.com/store/apps/details?id=com.google.android.apps.authenticator2" "iphone" "https://apps.apple.com/us/app/google-authenticator/id388497605" }}
					</li>
					<li>
						{{ tr "mfa-totp.template.authy" "link" "https://authy.com" }}
					</li>
				</ul>
			</div>
		</div>
	</div>
</div>
{{ template "inc_footer.html.tpl" . }}

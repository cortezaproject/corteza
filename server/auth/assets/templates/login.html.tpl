{{ template "inc_header.html.tpl" . }}
<div class="card-body p-0">
	<h4 class="login-title mb-0 p-3 border-bottom">{{ tr "login.template.title" }}</h4>
	{{ if .settings.LocalEnabled }}
	<form
		method="POST"
		action="{{ links.Login }}"
		class="p-3"
	>
		{{ .csrfField }}
		{{ if .form.error }}
		<div
			data-test-id="error"
			class="text-danger mb-4 font-weight-bold"
			role="alert"
		>
			{{ .form.error }}
		</div>
		{{ end }}
		<div class="mb-3">
		    <label>
                {{ tr "login.template.form.email.label" }}
            </label>
			<input
				data-test-id="input-email"
				type="email"
				class="form-control"
				name="email"
				required
				placeholder="{{ tr "login.template.form.email.placeholder" }}"
				value="{{ if .form }}{{ .form.email }}{{ end }}"
				autocomplete="username"
				aria-label="{{ tr "login.template.form.email.label" }}">
		</div>
		{{ if not .form.splitCredentialsCheck }}
		<div class="mb-3">
            <label>
                {{ tr "login.template.form.password.label" }}
            </label>
			<input
				data-test-id="input-password"
				type="password"
				required
				class="form-control"
				name="password"
				placeholder="{{ tr "login.template.form.password.placeholder" }}"
				autocomplete="current-password"
				aria-label="{{ tr "login.template.form.password.label" }}">
		</div>
		<div class="row">
			<div class="col text-right">
				{{ if .enableRememberMe }}
				<button
					data-test-id="button-login-and-remember"
					class="btn btn-primary btn-block btn-lg"
					name="keep-session"
					value="true"
					type="submit"
				>
					{{ tr "login.template.form.button.login-and-remember" }}
				</button>
				{{ end }}
				<button
					data-test-id="button-login"
					class="btn btn-light btn-block btn-lg mt-2"
					type="submit"
				>
					{{ tr "login.template.form.button.login" }}
				</button>
			</div>
		</div>
		{{ else }}
		<div class="row">
			<div class="col text-right">
				<button
					data-test-id="button-continue"
					class="btn btn-primary btn-block btn-lg"
					name="keep-session"
					value="true"
					type="submit"
				>
					{{ tr "login.template.form.button.continue" }}
				</button>
			</div>
		</div>
		{{ end }}
	</form>
	<div class="d-flex text-center px-3 pb-3 justify-content-around">
        {{ if .settings.PasswordResetEnabled }}
        <div>
            <a
							data-test-id="link-request-password-reset"
							href="{{ links.RequestPasswordReset }}"
						>
							{{ tr "login.template.links.request-password-reset" }}
						</a>
        </div>
        {{ end }}
        {{ if .settings.SignupEnabled }}
        <div>
            <a
							data-test-id="link-signup"
							href="{{ links.Signup }}"
						>
							{{ tr "login.template.links.signup" }}
						</a>
        </div>
        {{ end }}
	</div>
	{{ end }}

	{{ if .settings.ExternalEnabled }}
	<div class="px-3 pb-3">
		{{ range .providers }}
			<a href="{{ links.External }}/{{ .Handle }}" class="btn btn-light btn-block btn-lg mb-2 mt-1 text-dark">
				<svg class="mr-1" xmlns="http://www.w3.org/2000/svg" width="18" height="18" fill="currentColor" viewBox="0 0 16 16"><path d="M0 8a4 4 0 0 1 7.465-2H14a.5.5 0 0 1 .354.146l1.5 1.5a.5.5 0 0 1 0 .708l-1.5 1.5a.5.5 0 0 1-.708 0L13 9.207l-.646.647a.5.5 0 0 1-.708 0L11 9.207l-.646.647a.5.5 0 0 1-.708 0L9 9.207l-.646.647A.5.5 0 0 1 8 10h-.535A4 4 0 0 1 0 8m4-1a1 1 0 1 0 0 2 1 1 0 0 0 0-2"/></svg>
				{{ tr "login.template.links.external.login-with" "idp" (coalesce .Label .Handle) }}
			</a>
		{{ end }}
	</div>
	{{ end }}
</div>
{{ template "inc_footer.html.tpl" . }}

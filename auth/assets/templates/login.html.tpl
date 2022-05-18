{{ template "inc_header.html.tpl" . }}
<div class="card-body p-0">
	<h4 class="card-title p-3 border-bottom">{{ tr "login.template.title" }}</h4>
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
                {{ tr "login.template.form.email.label" }} *
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
                {{ tr "login.template.form.password.label" }} *
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
				<button
					data-test-id="button-login-and-remember"
					class="btn btn-primary btn-block btn-lg"
					name="keep-session"
					value="true"
					type="submit"
				>
					{{ tr "login.template.form.button.login-and-remember" }}
				</button>
				<button
					data-test-id="button-login"
					class="btn btn-light btn-block"
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
	<div class="row text-center pb-3">
        {{ if .settings.PasswordResetEnabled }}
        <div class="col cols-6">
            <a
							data-test-id="link-request-password-reset"
							href="{{ links.RequestPasswordReset }}"
						>
							{{ tr "login.template.links.request-password-reset" }}
						</a>
        </div>
        {{ end }}
        {{ if .settings.SignupEnabled }}
        <div class="col cols-6">
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
	<div class="pb-3">
	{{ range .providers }}
		<a href="{{ links.External }}/{{ .Handle }}" class="btn btn-light btn-block btn-lg mb-2 mt-1 text-dark">
			<i class="bi bi-{{ .Icon }} mr-1"></i>
			{{ tr "login.template.links.external.login-with" "idp" (coalesce .Label .Handle) }}
		</a>
	{{ end }}
	</div>
	{{ end }}
</div>
{{ template "inc_footer.html.tpl" . }}

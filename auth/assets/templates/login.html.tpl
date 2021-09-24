{{ template "inc_header.html.tpl" . }}
<div class="card-body p-0">
	<h1 class="h4 card-title p-3 border-bottom">{{ tr "login.template.title" }}</h1>
	{{ if .settings.LocalEnabled }}
	<form
		method="POST"
		action="{{ links.Login }}"
		class="p-3"
	>
		{{ .csrfField }}
		{{ if .form.error }}
		<div class="text-danger mb-4 font-weight-bold" role="alert">
			{{ .form.error }}
		</div>
		{{ end }}
		<div class="mb-3">
		    <label>
                {{ tr "login.template.form.email.label" }} *
            </label>
			<input
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
				type="password"
				required
				class="form-control"
				name="password"
				placeholder="{{ tr "login.template.form.password.placeholder" }}"
				autocomplete="current-password"
				aria-label="{{ tr "login.template.form.password.labels" }}">
		</div>
		<div class="row">
			<div class="col text-right">
				<button
					class="btn btn-primary btn-block btn-lg"
					name="keep-session"
					value="true"
					type="submit"
				>
					{{ tr "login.template.form.button.login-and-remember" }}
				</button>
				<button
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
	<div class="row text-center">
        {{ if .settings.PasswordResetEnabled }}
        <div class="col cols-6">
            <a href="{{ links.RequestPasswordReset }}">{{ tr "login.template.links.request-password-reset" }}</a>
        </div>
        {{ end }}
        {{ if .settings.SignupEnabled }}
        <div class="col cols-6">
            <a href="{{ links.Signup }}">{{ tr "login.template.links.signup" }}</a>
        </div>
        {{ end }}
	</div>
	{{ end }}
	{{ if .settings.ExternalEnabled }}
	<div class="p-2">
	{{ range .providers }}
		<a href="{{ links.External }}/{{ .Handle }}" class="btn btn-light btn-block btn-lg mb-2 mt-1 text-dark">
			<i class="bi bi-{{ .Icon }} mr-1"></i>
			{{ tr "login.template.external.login-with" "idp" (coalesce .Label .Handle) }}
		</a>
	{{ end }}
	</div>
	{{ end }}
</div>
{{ template "inc_footer.html.tpl" . }}

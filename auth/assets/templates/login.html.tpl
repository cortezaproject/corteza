{{ template "inc_header.html.tpl" . }}
<div class="card-body p-0">
	<h4 class="card-title p-3 border-bottom">Log in</h4>
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
                E-mail *
            </label>
			<input
				type="email"
				class="form-control"
				name="email"
				required
				placeholder="email@domain.ltd"
				value="{{ if .form }}{{ .form.email }}{{ end }}"
				autocomplete="username"
				aria-label="Email">
		</div>
		{{ if not .form.splitCredentialsCheck }}
		<div class="mb-3">
            <label>
                Password *
            </label>
			<input
				type="password"
				required
				class="form-control"
				name="password"
				placeholder="Password"
				autocomplete="current-password"
				aria-label="Password">
		</div>
		<div class="row">
			<div class="col text-right">
				<button
					class="btn btn-primary btn-block btn-lg"
					name="keep-session"
					value="true"
					type="submit"
				>
					Log in and remember me
				</button>
				<button
					class="btn btn-light btn-block"
					type="submit"
				>
					Log in
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
					Continue
				</button>
			</div>
		</div>
		{{ end }}
	</form>
	<div class="row text-center">
        {{ if .settings.PasswordResetEnabled }}
        <div class="col cols-6 mb-5">
            <a href="{{ links.RequestPasswordReset }}">Forgot your password?</a>
        </div>
        {{ end }}
        {{ if .settings.SignupEnabled }}
        <div class="col cols-6 mb-5">
            <a href="{{ links.Signup }}">Create a new account</a>
        </div>
        {{ end }}
	</div>
	{{ end }}
	{{ if .settings.ExternalEnabled }}
	<div class="p-3">
	{{ range .providers }}
		<a href="{{ links.External }}/{{ .Handle }}" class="btn btn-light btn-block btn-lg mb-2 text-dark">
			<i class="bi bi-{{ .Icon }} mr-1"></i>
			Login with {{ coalesce .Label .Handle }}
		</a>
	{{ end }}
	</div>
	{{ end }}
</div>
{{ template "inc_footer.html.tpl" . }}

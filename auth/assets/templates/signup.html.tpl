{{ template "inc_header.html.tpl" . }}
<div class="card-body p-0">
	{{ template "inc_alerts.html.tpl" .alerts }}
	<h4 class="card-title p-3 border-bottom">Sign up</h4>
	<form
		method="POST"
		onsubmit="buttonDisabler()"
		action="{{ links.Signup }}"
		class="p-3"
	>
		{{ .csrfField }}
		{{ if .form.error }}
		<div class="text-danger font-weight-bold mb-3" role="alert">
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
                autocomplete="username"
                value="{{ .form.email }}"
                aria-label="Email">
        </div>
        <div class="mb-3">
            <label>
                Password *
            </label>
			<input
				type="password"
				class="form-control"
				name="password"
				required
				placeholder="Password"
				autocomplete="new-password"
				aria-label="Password">
        </div>
        <div class="mb-3">
            <label>
                Full Name
            </label>
			<input
				type="text"
				class="form-control"
				name="name"
				placeholder="Your full name"
				value="{{ .form.name }}"
				autocomplete="name"
				aria-label="Full name">
        </div>
        <div class="mb-3">
            <label>
                Short name, nickname or handle
            </label>
			<input
				type="text"
				class="form-control"
				name="handle"
				placeholder="Short name, nickname or handle"
				value="{{ .form.handle }}"
				autocomplete="handle"
				aria-label="Handle">
        </div>
		<div>
			<button
				id="submit"
				class="btn btn-primary btn-block btn-lg"
				type="submit"
			>Submit</button>
		</div>
	</form>
	<div class="text-center my-3">Already have an account?
		<a href="{{ links.Login }}">Log in here</a>
	</div>
</div>
{{ template "inc_footer.html.tpl" . }}

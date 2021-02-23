{{ template "inc_header.html.tpl" . }}
<div class="card-body">
	{{ template "inc_alerts.html.tpl" .alerts }}
	<h4 class="card-title">Sign up</h4>
	<form
		method="POST"
		onsubmit="buttonDisabler()"
		action="{{ links.Signup }}"
	>
		{{ .csrfField }}
		{{ if .form.error }}
		<div class="alert alert-danger" role="alert">
			{{ .form.error }}
		</div>
		{{ end }}

		<div class="input-group mb-3">
		<span class="input-group-text">
			<i class="bi bi-envelope"></i>
		</span>
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

		<div class="input-group mb-3">
		<span class="input-group-text">
			<i class="bi bi-key-fill"></i>
		</span>
			<input
				type="password"
				class="form-control"
				name="password"
				required
				placeholder="Password"
				autocomplete="new-password"
				aria-label="Password">
		</div>
		<div class="input-group mb-3">
		<span class="input-group-text">
		  <i class="bi bi-person-fill"></i>
		</span>
			<input
				type="text"
				class="form-control"
				name="name"
				placeholder="Your full name"
				value="{{ .form.name }}"
				autocomplete="name"
				aria-label="Full name">
		</div>
		<div class="input-group mb-3">
		<span class="input-group-text">
			<i class="bi bi-emoji-smile"></i>
		</span>
			<input
				type="text"
				class="form-control"
				name="handle"
				placeholder="Short name, nickname or handle"
				value="{{ .form.handle }}"
				autocomplete="handle"
				aria-label="Handle">
		</div>
		<div class="text-right">
			<button
				id="submit"
				class="btn btn-primary btn-block"
				type="submit"
			>Submit</button>
		</div>
	</form>
	<div class="text-center my-3">Already have an account?
		<a href="{{ links.Login }}">Login</a>
	</div>
</div>
{{ template "inc_footer.html.tpl" . }}

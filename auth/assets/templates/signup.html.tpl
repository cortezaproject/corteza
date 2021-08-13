{{ template "inc_header.html.tpl" . }}
<div class="card-body p-0">
	<h1 class="h4 card-title p-3 border-bottom">{{ tr "singup.template.title" }}</h1>
	<form
		method="POST"
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
				{{ tr "singup.template.form.email.label" }}
            </label>
            <input
                type="email"
                class="form-control"
                name="email"
                required
                placeholder="{{ tr "singup.template.form.email.placeholder" }}"
                autocomplete="username"
                value="{{ .form.email }}"
                aria-label="{{ tr "singup.template.form.email.aria-label" }}">
        </div>
        <div class="mb-3">
            <label>
                {{ tr "singup.template.form.password.label" }}
            </label>
			<input
				type="password"
				class="form-control"
				name="password"
				required
				placeholder="{{ tr "singup.template.form.password.placeholder" }}"
				autocomplete="new-password"
				aria-label="{{ tr "singup.template.form.password.aria-label" }}">
        </div>
        <div class="mb-3">
            <label>
                {{ tr "singup.template.form.name.label" }}
            </label>
			<input
				type="text"
				class="form-control"
				name="name"
				placeholder="{{ tr "singup.template.form.name.placeholder" }}"
				value="{{ .form.name }}"
				autocomplete="name"
				aria-label="{{ tr "singup.template.form.name.aria-label" }}">
        </div>
        <div class="mb-3">
            <label>
                {{ tr "singup.template.form.nickname.label" }}
            </label>
			<input
				type="text"
				class="form-control handle-mask"
				name="handle"
				placeholder="{{ tr "singup.template.form.nickname.placeholder" }}"
				value="{{ .form.handle }}"
				autocomplete="handle"
				aria-label="{{ tr "singup.template.form.nickname.aria-label" }}">
        </div>
		<div>
			<button
				id="submit"
				class="btn btn-primary btn-block btn-lg"
				type="submit"
			>{{ tr "singup.template.form.button.sign-up" }}</button>
		</div>
	</form>
	<div class="text-center my-3">{{ tr "singup.template.form.link.alternative" }}
		<a href="{{ links.Login }}">{{ tr "singup.template.form.link.login" }}</a>
	</div>
</div>
{{ template "inc_footer.html.tpl" . }}

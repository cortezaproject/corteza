{{ template "inc_header.html.tpl" set . "activeNav" "profile" }}
<div class="card-body p-0">
	<h1 class="h4 card-title p-3 border-bottom">{{ tr "profile.template.title" }}</h1>
	<form
		method="POST"
		action="{{ links.Profile }}"
		class="p-3"
	>
		{{ .csrfField }}

		{{ if .form.error }}
		<div class="text-danger mb-4 font-weight-bold" role="alert">
			{{ .form.error }}
		</div>
		{{ end }}

        <div class="mb-3">
            <label for="profileFormEmail">{{ tr "profile.template.form.email.label" }}</label>
            <input
                type="email"
                class="form-control"
                name="email"
                id="profileFormEmail"
                placeholder="{{ tr "profile.template.form.email.placeholder" }}"
                autocomplete="username"
                readonly
                value="{{ .form.email }}"
                aria-label="{{ tr "profile.template.form.email.aria-label" }}">
            <div>
                {{ if .emailConfirmationRequired }}
                <div class="form-text text-danger">
                    {{ tr "profile.template.form.text-danger-1" }} <a href="{{ links.PendingEmailConfirmation }}?resend">{{ tr "profile.template.form.text-danger-2" }}</a>
                </div>
                {{ end }}
            </div>
        </div>

		<div class="mb-3">
			<label for="profileFormName">{{ tr "profile.template.form.name.label" }}</label>
            <input
                type="text"
                class="form-control"
                name="name"
                id="profileFormName"
                placeholder="{{ tr "profile.template.form.name.placeholder" }}"
                value="{{ .form.name }}"
                autocomplete="name"
                aria-label="{{ tr "profile.template.form.name.aria-label" }}">
		</div>

		<div class="mb-3">
			<label for="profileFormHandle">{{ tr "profile.template.form.handle.label" }}</label>
            <input
                type="text"
                class="form-control handle-mask"
                name="handle"
                id="profileFormHandle"
                placeholder="{{ tr "profile.template.form.handle.placeholder" }}"
                value="{{ .form.handle }}"
                autocomplete="handle"
                aria-label="{{ tr "profile.template.form.handle.aria-label" }}">
		</div>

        <div>
            <button
                type="submit"
                class="btn btn-primary btn-block btn-lg"
            >
                {{ tr "profile.template.form.button" }}
            </button>
        </div>
	</form>
</div>
{{ template "inc_footer.html.tpl" . }}

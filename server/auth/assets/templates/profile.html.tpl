{{ template "inc_header.html.tpl" set . "activeNav" "profile" }}
<div class="card-body p-0">
	<h4 class="card-title p-3 border-bottom">{{ tr "profile.template.title" }}</h4>
	<form
		method="POST"
		action="{{ links.Profile }}"
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
            <label for="profileFormEmail">{{ tr "profile.template.form.email.label" }}</label>
            <input
                data-test-id="input-email"
                type="email"
                class="form-control"
                name="email"
                id="profileFormEmail"
                placeholder="email@domain.ltd"
                autocomplete="username"
                readonly
                value="{{ .form.email }}"
                aria-label="{{ tr "profile.template.form.email.label" }}"
            >
            <div>
                {{ if .emailConfirmationRequired }}
                <div class="form-text text-danger">
                	{{ tr "profile.template.form.email.resend-confirmation-link" "link" links.PendingEmailConfirmation }}
                </div>
                {{ end }}
            </div>
        </div>

		<div class="mb-3">
			<label for="profileFormName">{{ tr "profile.template.form.name.label" }}</label>
            <input
                data-test-id="input-name"
                type="text"
                class="form-control"
                name="name"
                id="profileFormName"
                placeholder="{{ tr "profile.template.form.name.placeholder" }}"
                value="{{ .form.name }}"
                autocomplete="name"
                aria-label="{{ tr "profile.template.form.name.label" }}"
            >
		</div>

		<div class="mb-3">
			<label for="profileFormHandle">{{ tr "profile.template.form.handle.label" }}</label>
            <input
                data-test-id="input-handle"
                type="text"
                class="form-control handle-mask"
                name="handle"
                id="profileFormHandle"
                placeholder="{{ tr "profile.template.form.handle.placeholder" }}"
                value="{{ .form.handle }}"
                autocomplete="handle"
                aria-label="{{ tr "profile.template.form.handle.label" }}"
            >
		</div>


		<div class="mb-3">
			<label for="profileFormPreferredLanguage">{{ tr "profile.template.form.preferred-language.label" }}</label>
			<select
                data-test-id="select-language"
                class="form-control"
				name="preferredLanguage"
                id="profileFormPreferredLanguage"
                aria-label="{{ tr "profile.template.form.preferred-language.label" }}"
                value="{{ .form.preferredLanguage }}"
			>
			{{ $prefLang := .form.preferredLanguage }}
			{{ range .languages }}
				<option
					value="{{ .Tag }}"
					{{ if eq $prefLang .Tag.String }}selected{{ end }}
				>
					{{ .LocalizedName }} ({{ .Name }})
				</option>
			{{ end }}
			</select>
		</div>

        <div>
            <button
                data-test-id="button-submit"
                type="submit"
                class="btn btn-primary btn-block btn-lg"
            >
                {{ tr "profile.template.form.buttons.submit" }}
            </button>
        </div>
	</form>
</div>
{{ template "inc_footer.html.tpl" . }}

{{ template "inc_header.html.tpl" set . "activeNav" "profile" }}
<div class="card-body p-0">
	<form
		method="POST"
		action="{{ links.Profile }}"
        enctype="multipart/form-data"
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

        {{ if .avatarEnabled }}
        <hr/>
        <div class="mb-3">
            <label>{{ tr "profile.template.form.avatar.label" }}</label>
            <div class="d-flex align-items-center" style="gap: 1rem;">
                {{ if .isAvatar }}
                <div class="avatar-preview">
                    <img src="{{ .form.avatarUrl }}" alt="Profile Photo">
                </div>
                {{ else }}
                <div class="avatar-preview avatar-initials" id="avatarInitials"
                     style="background-color: {{ if .form.initialBgColor }}{{ .form.initialBgColor }}{{ else }}#e4e4e7{{ end }}; color: {{ if .form.initialTextColor }}{{ .form.initialTextColor }}{{ else }}#3f3f46{{ end }};">
                </div>
                {{ end }}
                <div class="d-flex align-items-center" style="gap: 0.5rem;">
                    <label for="avatar" class="btn btn-light btn-sm" style="margin-bottom: 0; color: var(--auth-text-color); cursor: pointer;">
                    {{ tr "profile.template.form.avatar.upload" }}
                    </label>
                    <input id="avatar" name="avatar" value="avatar" type="file" class="sr-only" accept="image/*">

                    {{  if .isAvatar }}
                    <button
                        name="avatar-delete"
                        value="avatar-delete"
                        class="btn btn-danger btn-sm"
                    >
                        {{ tr "profile.template.form.avatar.delete" }}
                    </button>
                    {{ end }}
                </div>
            </div>
        </div>


        {{ end }}

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

{{ template "inc_header.html.tpl" set . "activeNav" "profile" }}
<div class="card-body p-0">
	<h1 class="h4 card-title p-3 border-bottom">Your profile</h1>
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
            <label for="profileFormEmail">Email</label>
            <input
                type="email"
                class="form-control"
                name="email"
                id="profileFormEmail"
                placeholder="email@domain.ltd"
                autocomplete="username"
                readonly
                value="{{ .form.email }}"
                aria-label="Email">
            <div>
                {{ if .emailConfirmationRequired }}
                <div class="form-text text-danger">
                    Email is not verified, <a href="{{ links.PendingEmailConfirmation }}?resend">resend confirmation link.</a>
                </div>
                {{ end }}
            </div>
        </div>

		<div class="mb-3">
			<label for="profileFormName">Full name</label>
            <input
                type="text"
                class="form-control"
                name="name"
                id="profileFormName"
                placeholder="Your full name"
                value="{{ .form.name }}"
                autocomplete="name"
                aria-label="Full name">
		</div>

		<div class="mb-3">
			<label for="profileFormHandle">Handle</label>
            <input
                type="text"
                class="form-control handle-mask"
                name="handle"
                id="profileFormHandle"
                placeholder="Short name, nickname or handle"
                value="{{ .form.handle }}"
                autocomplete="handle"
                aria-label="Handle">
		</div>

        <div>
            <button
                type="submit"
                class="btn btn-primary btn-block btn-lg"
            >
                Update profile
            </button>
        </div>
	</form>
</div>
{{ template "inc_footer.html.tpl" . }}

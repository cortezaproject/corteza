{{ template "inc_header.html.tpl" . }}
<div class="card-body p-0">
  <h1 class="h4 card-title p-3 border-bottom">{{ tr "signup.template.title" }}</h1>
  <form
		method="POST"
		action="{{ links.Signup }}"
		class="p-3"
	>
    {{ .csrfField }}
    {{ if .form.error }}
    <div class="custom-row text-danger font-weight-bold mb-3 pl-3" role="alert">
      {{ .form.error }}
    </div>
    {{ end }}
		
		<div class="custom-row">
			<div class="custom-col">
				<div class="mb-3">
					<label>
						{{ tr "signup.template.form.email.label" }}
					</label>
					<input
						type="email"
						class="form-control text-truncate"
						name="email"
						required
						placeholder="{{ tr "signup.template.form.email.placeholder" }}"
						autocomplete="username"
						value="{{ .form.email }}"
						aria-label="{{ tr "signup.template.form.email.label" }}">
				</div>
				<div class="mb-3 margin-auto">
					<label>
						{{ tr "signup.template.form.password.label" }}
					</label>
					<input
						type="password"
						class="form-control text-truncate"
						name="password"
						required
						placeholder="{{ tr "signup.template.form.password.placeholder" }}"
						autocomplete="new-password"
						aria-label="{{ tr "signup.template.form.password.label" }}">
				</div>
			</div>
			<div class="custom-col margin-auto">
				<div class="mb-3">
					<label>
						{{ tr "signup.template.form.name.label" }}
					</label>
					<input
						type="text"
						class="form-control text-truncate"
						name="name"
						placeholder="{{ tr "signup.template.form.name.placeholder" }}"
						value="{{ .form.name }}"
						autocomplete="name"
						aria-label="{{ tr "signup.template.form.name.label" }}">
				</div>
				<div class="mb-3">
					<label>
						{{ tr "signup.template.form.nickname.label" }}
					</label>
					<input
						type="text"
						class="form-control handle-mask text-truncate"
						name="handle"
						placeholder="{{ tr "signup.template.form.nickname.placeholder" }}"
						value="{{ .form.handle }}"
						autocomplete="handle"
						aria-label="{{ tr "signup.template.form.nickname.label" }}">
				</div>
			</div>
		</div>
		<div class="custom-row">
      <button
				id="submit"
				class="btn btn-primary btn-block btn-lg"
				type="submit"
			>{{ tr "signup.template.form.button.sign-up" }}</button>
		</div>
  </form>
  <div class="text-center mb-3">
    {{ tr "signup.template.log-in" "link" links.Login }}
  </div>
</div>

{{ template "inc_footer.html.tpl" . }}

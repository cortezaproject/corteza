{{ template "inc_header.html.tpl" . }}
<div class="card-body">
	{{ template "inc_alerts.html.tpl" .alerts }}
	<h4 class="card-title">Your profile</h4>
	<form
		method="POST"
		action="{{ links.Profile }}"
		onsubmit="buttonDisabler()"
	>
		{{ .csrfField }}

		{{ if .form.error }}
		<div class="alert alert-danger" role="alert">
			{{ .form.error }}
		</div>
		{{ end }}

		<div class="form-group row">
			<label for="profileFormEmail" class="col-sm-2 col-form-label">Email</label>
			<div class="col-sm-10">
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

				{{ if .emailConfirmationRequired }}
				<small class="form-text text-danger">
					Email is not verified, <a href="{{ links.PendingEmailConfirmation }}?resend">resend confirmation link.</a>
				</small>
				{{ end }}
			</div>
		</div>

		<div class="form-group row">
			<label for="profileFormName" class="col-sm-2 col-form-label">Full name</label>
			<div class="col-sm-10">
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
		</div>

		<div class="form-group row">
			<label for="profileFormHandle" class="col-sm-2 col-form-label">Handle</label>
			<div class="col-sm-10">
				<input
					type="text"
					class="form-control"
					name="handle"
					id="profileFormHandle"
					placeholder="Short name, nickname or handle"
					value="{{ .form.handle }}"
					autocomplete="handle"
					aria-label="Handle">
			</div>
		</div>

		<div class="form-group row">
			<div class="col-12 text-right">
				<button
					type="submit"
					class="btn btn-primary"
				>
					Update profile
				</button>
			</div>
		</div>
	</form>
</div>
{{ template "inc_footer.html.tpl" . }}

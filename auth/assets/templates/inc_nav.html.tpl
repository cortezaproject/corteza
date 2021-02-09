<div class="card-header p-0">
	<div class="row">
		<div class="col cols-6">
			<a href="{{ links.Profile }}">
				<img class="logo m-3" src="{{ links.Assets }}/logo.png">
			</a>
		</div>
		<div class="col cols-6 text-right small m-3">
			{{ if .user }}
			Logged-in as
			<a href="{{ links.Profile }}">{{ coalesce .user.Name .user.Handle .user.Email }}</a>
			|
			<a href="{{ links.Logout }}">Logout</a>
			{{ end }}
		</div>
	</div>

{{ if not .hideNav }}
	{{ if and .user .client }}
	<div class="py-1 px-3">
		<a href="{{ links.OAuth2AuthorizeClient }}">&raquo; Continue with authorization of {{ .client.Name }}</a>
	</div>
	{{ else if .user }}
	<ul class="nav ml-1">
		<li class="nav-item">
			<a class="nav-link" href="{{ links.Profile }}">Your profile</a>
		</li>
		<li class="nav-item">
			<a class="nav-link" href="{{ links.Sessions }}">Login sessions</a>
		</li>
		{{ if .settings.ExternalEnabled }}
		<li class="nav-item">
			<a class="nav-link" href="{{ links.AuthorizedClients }}">Authorized clients</a>
		</li>
		{{ end }}
		{{ if .settings.LocalEnabled }}
		<li class="nav-item">
			<a class="nav-link" href="{{ links.ChangePassword }}">Change password</a>
		</li>
		{{ end }}
	</ul>
	{{ end }}
{{ end }}
</div>

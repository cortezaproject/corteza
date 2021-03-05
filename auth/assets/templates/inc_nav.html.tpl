<div class="card-header border-0">
	<div class="row">
		<div class="text-center w-100">
			<a href="{{ links.Profile }}">
				<img class="logo m-2" src="{{ links.Assets }}/logo.png">
			</a>
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

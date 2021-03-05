<div class="card-header border-0 p-0">
    <div class="text-center w-100">
        <a href="{{ links.Profile }}">
            <img class="logo m-3" src="{{ links.Assets }}/logo.png">
        </a>
    </div>

{{ if not .hideNav }}
	{{ if and .user .client }}
	<div class="py-1 px-3">
		<a href="{{ links.OAuth2AuthorizeClient }}" class="text-danger">
		 Finalize the authorization of {{ .client.Name }}
		 <i class="bi bi-chevron-double-right"></i>
		 </a>
	</div>
	{{ else if .user }}
	<ul class="nav ml-1 d-flex justify-content-around">
	   {{/* @TO-DO Denis -- apply active class to selected nav item */}}
		<li class="nav-item active">
			<a class="nav-link" href="{{ links.Profile }}">Your profile</a>
		</li>
		{{ if .settings.LocalEnabled }}
		<li class="nav-item">
			<a class="nav-link" href="{{ links.ChangePassword }}">Security</a>
		</li>
		{{ end }}
		<li class="nav-item">
			<a class="nav-link" href="{{ links.Sessions }}">Login sessions</a>
		</li>
		{{ if .settings.ExternalEnabled }}
		<li class="nav-item">
			<a class="nav-link" href="{{ links.AuthorizedClients }}">Authorized clients</a>
		</li>
		{{ end }}
	</ul>
	{{ end }}
{{ end }}
</div>

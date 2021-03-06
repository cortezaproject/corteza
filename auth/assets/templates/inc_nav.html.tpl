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
		<li class="nav-item {{ if eq .activeNav "profile" }}active{{ end  }}">
			<a class="nav-link" href="{{ links.Profile }}">Your profile</a>
		</li>
		<li class="nav-item {{ if eq .activeNav "security" }}active{{ end  }}">
			<a class="nav-link" href="{{ links.Security }}">Security</a>
		</li>
		<li class="nav-item {{ if eq .activeNav "sessions" }}active{{ end  }}">
			<a class="nav-link" href="{{ links.Sessions }}">Login sessions</a>
		</li>
		{{ if .settings.ExternalEnabled }}
		<li class="nav-item {{ if eq .activeNav "clients" }}active{{ end  }}">
			<a class="nav-link" href="{{ links.AuthorizedClients }}">Authorized clients</a>
		</li>
		{{ end }}
	</ul>
	{{ end }}
{{ end }}
</div>

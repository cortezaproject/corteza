<div class="card-header border-0 p-0">
    <div class="text-center w-100">
        <a href="{{ links.Profile }}">
            <img class="logo m-3" alt="Company logo" src="{{ links.Assets }}/logo.png">
        </a>
    </div>

{{ $activeNav := default "" .activeNav }}

{{ if not .hideNav }}
	{{ if and .user .client }}
	<div class="py-1 px-3">
		<a href="{{ links.OAuth2AuthorizeClient }}" class="text-danger">
		 {{ tr "inc_nav.template.authorize-client" }} {{ .client.Name }}
		 <i class="bi bi-chevron-double-right"></i>
		 </a>
	</div>
	{{ else if .user }}
	<ul class="nav ml-1 d-flex justify-content-around">
		<li class="nav-item {{ if eq $activeNav "profile" }}active{{ end  }}">
			<a class="nav-link" href="{{ links.Profile }}">{{ tr "inc_nav.template.class.your-profile" }}</a>
		</li>
		<li class="nav-item {{ if eq $activeNav "security" }}active{{ end  }}">
			<a class="nav-link" href="{{ links.Security }}">{{ tr "inc_nav.template.class.security" }}</a>
		</li>
		<li class="nav-item {{ if eq $activeNav "sessions" }}active{{ end  }}">
			<a class="nav-link" href="{{ links.Sessions }}">{{ tr "inc_nav.template.class.login-session" }}</a>
		</li>
		<li class="nav-item {{ if eq $activeNav "clients" }}active{{ end  }}">
			<a class="nav-link" href="{{ links.AuthorizedClients }}">{{ tr "inc_nav.template.class.authorized-clients" }}</a>
		</li>
	</ul>
	{{ end }}
{{ end }}
</div>

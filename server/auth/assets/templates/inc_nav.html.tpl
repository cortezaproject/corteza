<div class="header card-header border-0 p-0">
    <div class="text-center w-100 my-4">
        <a href="{{ links.Profile }}">
            <img
							data-test-id="img-corteza-logo"
							class="logo"
							alt="Company logo"
							src="{{ links.Assets }}/logo.png"
						>
        </a>
    </div>
{{ $activeNav := default "" .activeNav }}

{{ if not .hideNav }}
	{{ if and .user .client }}
	<div class="py-1 px-3">
		<a
			data-test-id="link-finalize-authorization"
			class="text-danger"
			href="{{ links.OAuth2AuthorizeClient }}"
		>
		 {{ tr "inc_nav.template.authorize-client" }} {{ .client.Name }}
		 <i class="bi bi-chevron-double-right"></i>
		 </a>
	</div>
	{{ else if .user }}
	<ul class="nav ml-1 d-flex justify-content-around">
		<li class="nav-item {{ if eq $activeNav "profile" }}active{{ end  }}">
			<a
				data-test-id="link-tab-profile"
				class="nav-link"
				href="{{ links.Profile }}"
			>
				{{ tr "inc_nav.template.class.your-profile" }}
			</a>
		</li>
		<li class="nav-item {{ if eq $activeNav "security" }}active{{ end  }}">
			<a
				data-test-id="link-tab-security"
				class="nav-link"
				href="{{ links.Security }}"
			>
				{{ tr "inc_nav.template.class.security" }}
			</a>
		</li>
		<li class="nav-item {{ if eq $activeNav "sessions" }}active{{ end  }}">
			<a
				data-test-id="link-tab-login-session"
				class="nav-link"
				href="{{ links.Sessions }}"
			>
				{{ tr "inc_nav.template.class.login-session" }}
			</a>
		</li>
		<li class="nav-item {{ if eq $activeNav "clients" }}active{{ end  }}">
			<a
				data-test-id="link-tab-authorized-clients"
				class="nav-link"
				href="{{ links.AuthorizedClients }}"
			>
				{{ tr "inc_nav.template.class.authorized-clients" }}
			</a>
		</li>
	</ul>
	{{ end }}
{{ end }}
</div>

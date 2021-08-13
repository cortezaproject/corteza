{{ template "inc_header.html.tpl" set . "activeNav" "clients" }}
<div class="card-body p-0">
    <h1 class="h4 card-title p-3 border-bottom">{{ tr "authorized-clients.template.title" }}</h1>
	<form
		method="POST"
		class="clearfix"
		action="{{ links.AuthorizedClients }}"
		class="p-3"
	>

		{{ .csrfField }}

	{{ range .authorizedClients}}
        <div class="p-3">
            <div class="text-primary font-weight-bold">{{ .Name }}</div>
            <div>
                {{ tr "authorized-clients.template.authorized" }}
                <time
                    datetime="{{ .ConfirmedAt | date "2006-01-02T15:04:05Z07:00" }}"
                >
                    {{ .ConfirmedAt | date "Mon, 02 Jan 2006 15:04 MST" }}
                </time>
            </div>
            <button
                type="submit"
                name="revoke"
                value="{{ .ID }}"
                class="btn btn-sm btn-danger"
            >
                {{ tr "authorized-clients.template.button.revoke-access" }}
            </button>
        </div>
	{{ else }}
		<div class="text-center m-3 mb-5">
			<i>{{ tr "authorized-clients.template.no-auth-clients" }}</i>
		</div>
	{{ end }}
	</form>
</div>
{{ template "inc_footer.html.tpl" . }}

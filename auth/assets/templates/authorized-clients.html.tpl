{{ template "inc_header.html.tpl" . }}
<div class="card-body p-0">
	{{ template "inc_alerts.html.tpl" .alerts }}
    <h4 class="card-title p-3 border-bottom">Authorized clients</h4>
	<form
		method="POST"
		class="clearfix"
		action="{{ links.AuthorizedClients }}"
		onsubmit="buttonDisabler()"
		class="p-3"
	>

		{{ .csrfField }}

	{{ range .authorizedClients}}
        <div class="p-3">
            <div class="text-primary font-weight-bold">{{ .Name }}</div>
            <div>
                Authorized on
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
                Revoke access
            </button>
        </div>
	{{ else }}
		<div class="text-center m-3">
			<i>No authorized clients found</i>
		</div>
	{{ end }}
	</form>
</div>
{{ template "inc_footer.html.tpl" . }}

{{ template "inc_header.html.tpl" set . "activeNav" "clients" }}
<div class="card-body p-0">
    <h4 class="card-title p-3 border-bottom">{{ tr "authorized-clients.template.title" }}</h4>
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
                {{ tr "authorized-clients.template.list.authorized-on" }}
                <time
                    datetime="{{ .ConfirmedAt | date "2006-01-02T15:04:05Z07:00" }}"
                >
                    {{ .ConfirmedAt | date "Mon, 02 Jan 2006 15:04 MST" }}
                </time>
            </div>
            <button
                data-test-id="button-revoke-access"
                type="submit"
                name="revoke"
                value="{{ .ID }}"
                class="btn btn-sm btn-danger"
            >
                {{ tr "authorized-clients.template.list.buttons.revoke" }}
            </button>
        </div>
        {{ end }}
    </form>

    {{ if not .authorizedClients}}
    <div
        data-test-id="text-empty-list"
        class="text-center m-3 mb-3"
    >
        <i>{{ tr "authorized-clients.template.list.empty" }}</i>
    </div>
    {{ end }}
</div>
{{ template "inc_footer.html.tpl" . }}

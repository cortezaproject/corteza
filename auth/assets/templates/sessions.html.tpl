{{ template "inc_header.html.tpl" set . "activeNav" "sessions" }}
<div class="card-body p-0">
	<form
		method="POST"
		action="{{ links.Sessions }}"
		class="p-3"
	>
		{{ .csrfField }}
		<div class="card-title p-3 border-bottom">
			<h4 class="card-title d-inline">{{ tr "sessions.template.title" }}</h4>
			{{- if and .sessions (gt (len .sessions) 1) }}
			<button
				type="submit"
				name="delete-all-but-current"
				value="true"
				class="btn btn-sm btn-danger float-right"
			>
				{{ tr "sessions.template.delete-all" }}
			</button>
		{{ end }}
		</div>
	</form>
	<form
		method="POST"
		action="{{ links.Sessions }}"
		class="p-3"
	>
		{{ .csrfField }}

	{{ range .sessions}}
		<div class="mb-3 border-bottom">
			{{ if .Current }}
                <h5>{{ tr "sessions.template.list.current" }}</h5>
			{{ end }}
            {{ if not .Current }}
                <button
                    type="submit"
                    name="delete"
                    value="{{ .ID }}"
                    class="btn btn-sm btn-link text-danger float-right"
                >
                    {{ tr "sessions.template.list.delete" }}
                </button>
            {{ end }}
			<label class="mb-0 d-block">{{ tr "sessions.template.list.authorized-on" }}</label>
			<p class="w-75 d-inline-block">
                <time datetime="{{ .CreatedAt }}">{{ .CreatedAt | date "Mon, 02 Jan 2006 15:04 MST" }}</time>
            </p>
			<label class="mb-0 d-block">{{ tr "sessions.template.list.ip-address" }}</label>
			<p class="w-75 d-inline-block">{{ .RemoteAddr }}</p>
			{{ if .SameRemoteAddr}}
                <span class="badge badge-light float-right">{{ tr "sessions.template.list.same-machine" }}</span>
            {{ end }}
			<label class="mb-0 d-block">{{ tr "sessions.template.list.expires" }}</label>
			<p class="w-75 d-inline-block">
			    <time datetime="{{ .ExpiresAt }}">{{ .ExpiresAt | date "Mon, 02 Jan 2006 15:04 MST" }}</time>
			</p>
			{{ if .Expired }}
                <span class="badge badge-warning float-right">{{ tr "sessions.template.list.expired" }}</span>
                {{ else if eq .ExpiresIn 0 }}
                <span class="badge badge-warning float-right">{{ tr "sessions.template.list.today" }}</span>
                {{ else if eq .ExpiresIn 1 }}
                <span class="badge badge-light float-right">{{ tr "sessions.template.list.tomorrow" }}</span>
                {{ else }}
                <span class="badge badge-light float-right">{{ tr "sessions.template.list.soon" "days" .ExpiresIn }}</span>
                {{ end }}
			<label class="mb-0 d-block">{{ tr "sessions.template.list.browser" }}</label>
			<p class="small w-75 d-inline-block">
			    {{ .UserAgent }}
			</p>
			{{ if .SameUserAgent}}
                <span class="badge badge-light float-right">{{ tr "sessions.template.list.same-browser" }}</span>
            {{ end }}
		</div>
	{{ end }}
	</form>
</div>
{{ template "inc_footer.html.tpl" . }}

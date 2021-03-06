{{ template "inc_header.html.tpl" set . "activeNav" "sessions" }}
<div class="card-body p-0">
	    <div class="card-title p-3 border-bottom">
            <h4 class="card-title d-inline">Your sessions</h4>
            {{ if .sessions}}
            <button
                type="submit"
                name="delete-all-but-current"
                value="true"
                class="btn btn-sm btn-danger float-right"
            >
                Delete all sessions
            </button>
            {{ end }}
	    </div>
	<form
		method="POST"
		action="{{ links.Sessions }}"
		class="p-3"
	>
		{{ .csrfField }}

	{{ range .sessions}}
		<div class="mb-3 border-bottom">
			{{ if .Current }}
                <h5>Current session</h5>
			{{ end }}
            {{ if not .Current }}
                <button
                    type="submit"
                    name="delete"
                    value="{{ .ID }}"
                    class="btn btn-sm btn-link text-danger float-right"
                >
                    Delete this session
                </button>
            {{ end }}
			<label class="mb-0 d-block">Authorized on</label>
			<p class="w-75 d-inline-block">
                <time datetime="{{ .CreatedAt }}">{{ .CreatedAt | date "Mon, 02 Jan 2006 15:04 MST" }}</time>
            </p>
			<label class="mb-0 d-block">IP Address</label>
			<p class="w-75 d-inline-block">{{ .RemoteAddr }}</p>
			{{ if .SameRemoteAddr}}
                <span class="badge badge-light float-right">This machine</span>
            {{ end }}
			<label class="mb-0 d-block">Expires</label>
			<p class="w-75 d-inline-block">
			    <time datetime="{{ .ExpiresAt }}">{{ .ExpiresAt | date "Mon, 02 Jan 2006 15:04 MST" }}</time>
			</p>
			{{ if .Expired }}
                <span class="badge badge-warning float-right">expired</span>
                {{ else if eq .ExpiresIn 0 }}
                <span class="badge badge-warning float-right">today</span>
                {{ else if eq .ExpiresIn 1 }}
                <span class="badge badge-light float-right">in 1 day</span>
                {{ else }}
                <span class="badge badge-light float-right">in {{ .ExpiresIn }} days</span>
                {{ end }}
			<label class="mb-0 d-block">Browser</label>
			<p class="small w-75 d-inline-block">
			    {{ .UserAgent }}
			</p>
			{{ if .SameUserAgent}}
                <span class="badge badge-light float-right">This browser</span>
            {{ end }}
		</div>
	{{ end }}
	</form>
</div>
{{ template "inc_footer.html.tpl" . }}

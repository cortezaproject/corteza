{{ template "inc_header.html.tpl" . }}
<div class="card-body">
	{{ template "inc_alerts.html.tpl" .alerts }}

	<form
		method="POST"
		action="{{ links.Sessions }}"
		onsubmit="buttonDisabler()"
	>
		{{ if .sessions}}
		<button
			type="submit"
			name="delete-all-but-current"
			value="true"
			class="btn btn-sm btn-link text-danger float-right"
		>
			Delete all
		</button>
		{{ end }}

		<h4 class="card-title">Your sessions</h4>
		{{ .csrfField }}

	{{ range .sessions}}
		<div class="card mb-3 shadow-sm rounded-lg">
			{{ if .Current }}
			<div class="card-header">
				<b>Current session</b>
			</div>
			{{ end }}

			<div class="card-body bg-light">

				<div class="row clearfix">
					<div class="col">
						<dt>
							Authorized on
						</dt>
						<dd>
							<time datetime="{{ .CreatedAt }}">{{ .CreatedAt | date "Mon, 02 Jan 2006 15:04 MST" }}</time>
						</dd>
					</div>
					<div class="col">
						<dt>
							{{ if .Expired }}
							<span class="badge badge-warning float-right">expired</span>
							{{ else if eq .ExpiresIn 0 }}
							<span class="badge badge-warning float-right">today</span>
							{{ else if eq .ExpiresIn 1 }}
							<span class="badge badge-warning float-right">in 1 day</span>
							{{ else }}
							<span class="badge badge-info float-right">in {{ .ExpiresIn }} days</span>
							{{ end }}
							Expires
						</dt>
						<dd>
							<time datetime="{{ .ExpiresAt }}">{{ .ExpiresAt | date "Mon, 02 Jan 2006 15:04 MST" }}</time>
						</dd>
					</div>
				</div>
				<div class="row">
					<div class="col">
						<dt>
							{{ if .SameRemoteAddr}}
								<span class="badge badge-info float-right">This machine</span>
							{{ end }}
							IP address
						</dt>
						<dd>
							{{ .RemoteAddr }}
						</dd>
					</div>
					<div class="col">
						<dt>
							{{ if .SameUserAgent}}
								<span class="badge badge-info float-right">This browser</span>
							{{ end }}
							Browser
						</dt>
						<dd class="small">
							{{ .UserAgent }}
						</dd>
					</div>
				</div>

			</div>
			{{ if not .Current }}
			<div class="card-footer">
				<button
					type="submit"
					name="delete"
					value="{{ .ID }}"
					class="btn btn-sm btn-link text-danger"
				>
					Delete this session
				</button>
			</div>
			{{ end }}
		</div>
	{{ end }}
	</form>
</div>
{{ template "inc_footer.html.tpl" . }}

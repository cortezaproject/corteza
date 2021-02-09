{{ template "inc_header.html.tpl" . }}
<div class="card-body">
	{{ template "inc_alerts.html.tpl" .alerts }}
	<form
		method="POST"
		class="clearfix"
		action="{{ links.AuthorizedClients }}"
		onsubmit="buttonDisabler()"
	>
		<h4 class="card-title">Authorized clients</h4>

		{{ .csrfField }}

	{{ range .authorizedClients}}
			<div class="card-body bg-light">
				<div class="row">
					<div class="col-9">
						{{ .Name }}
						<div class="small">
							Authorized on
							<time
								datetime="{{ .ConfirmedAt | date "2006-01-02T15:04:05Z07:00" }}"
							>
								{{ .ConfirmedAt | date "Mon, 02 Jan 2006 15:04 MST" }}
							</time>
						</div>
					</div>
					<div class="col-3 text-right">
						<button
							type="submit"
							name="revoke"
							value="{{ .ID }}"
							class="btn btn-sm btn-link text-danger"
						>
							Revoke access
						</button>
					</div>
				</div>
			</div>
	{{ else }}
		<div class="text-center m-3">
			<i>No authorized clients found</i>
		</div>
	{{ end }}
	</form>
</div>
{{ template "inc_footer.html.tpl" . }}

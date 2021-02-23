{{/* setting hideNav=true to root scope and passign it to header template */}}
{{ template "inc_header.html.tpl" set . "hideNav" true }}
<div class="card-body">
	{{ template "inc_alerts.html.tpl" .alerts }}

	{{ if .invalidUser }}
	<div class="alert alert-danger" role="alert">
		{{ .invalidUser }}
	</div>
	{{ end }}

	<h4 class="card-title">Authorize "{{ coalesce .client.Name }}"</h4>

	<form
		action="{{ links.OAuth2AuthorizeClient }}"
		method="POST"
		onsubmit="buttonDisabler()"
	>
	  {{ .csrfField }}
	  <p>
	  	Hello {{ coalesce .user.Name .user.Handle .user.Email }},
	  </p>
	  <p>
		  <b>{{ .client.Name }}</b> would like to perform actions on this Corteza server on your behalf.
	  </p>

	  <p class="text-center">
		<button
		  type="submit"
		  name="deny"
		  class="btn btn-secondary btn-lg m-2"
		  style="width:180px;"
		>
		  Deny
		</button>
		<button
		  type="submit"
		  name="allow"
		  {{ if .disabled }}disabled{{ end }}
		  class="btn btn-{{ if .disabled }}secondary{{ else }}primary{{ end }} btn-lg m-2"
		  style="width:180px;"
		>
		  Allow
		</button>
	  </p>

	  <hr />

	  If this is a mistake, <a href="{{ links.Logout }}">logout</a>.

	</form>
</div>
 {{ template "inc_footer.html.tpl" . }}

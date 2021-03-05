{{/* setting hideNav=true to root scope and passign it to header template */}}
{{ template "inc_header.html.tpl" set . "hideNav" true }}
<div class="card-body p-0">
	{{ template "inc_alerts.html.tpl" .alerts }}

	{{ if .invalidUser }}
	<div class="text-danger font-weight-bold p-3" role="alert">
		{{ .invalidUser }}
	</div>
	{{ end }}

	<h4 class="card-title p-3 border-bottom">Authorize "{{ coalesce .client.Name }}"</h4>

	<form
		action="{{ links.OAuth2AuthorizeClient }}"
		method="POST"
		onsubmit="buttonDisabler()"
		class="p-3"
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
		  name="allow"
		  {{ if .disabled }}disabled{{ end }}
		  class="btn btn-{{ if .disabled }}secondary{{ else }}primary{{ end }} btn-lg m-2"
		  style="width:250px;"
		>
		  Allow
		</button>
		<button
		  type="submit"
		  name="deny"
		  class="btn btn-danger btn-lg m-2"
		  style="width:250px;"
		>
		  Deny
		</button>
	  </p>
      <div class="text-center">
	    If this is a mistake, please <a href="{{ links.Logout }}">log out</a>.
      </div>

	</form>
</div>
 {{ template "inc_footer.html.tpl" . }}

{{/* setting hideNav=true to root scope and passign it to header template */}}
{{ template "inc_header.html.tpl" set . "hideNav" true }}
<div class="card-body p-0">
	{{ if .invalidUser }}
	<div class="text-danger font-weight-bold p-3" role="alert">
		{{ .invalidUser }}
	</div>
	{{ end }}

	<h1 class="h4 card-title p-3 border-bottom">{{ tr "oauth2_authorize_client.template.title" }} "{{ coalesce .client.Name }}"</h1>

	<form
		action="{{ links.OAuth2AuthorizeClient }}"
		method="POST"
		class="p-3"
	>
	  {{ .csrfField }}
	  <p>
	  	{{ tr "oauth2_authorize_client.template.form.greeting-paragraph" }} {{ coalesce .user.Name .user.Handle .user.Email }},
	  </p>
	  <p>
		  <b>{{ .client.Name }}</b> {{ tr "oauth2_authorize_client.template.form.question-for-client" }}
	  </p>

	  <p class="text-center">
		<button
		  type="submit"
		  name="allow"
		  {{ if .disabled }}disabled{{ end }}
		  class="btn btn-{{ if .disabled }}secondary{{ else }}primary{{ end }} btn-lg m-2"
		  style="width:250px;"
		>
		  {{ tr "oauth2_authorize_client.template.form.button-allow" }}
		</button>
		<button
		  type="submit"
		  name="deny"
		  class="btn btn-danger btn-lg m-2"
		  style="width:250px;"
		>
		  {{ tr "oauth2_authorize_client.template.form.button-deny" }}
		</button>
	  </p>
      <div class="text-center">
	    {{ tr "oauth2_authorize_client.template.form.mistake" }} <a href="{{ links.Logout }}">{{ tr "oauth2_authorize_client.template.form.log-out-link" }}</a>.
      </div>

	</form>
</div>
 {{ template "inc_footer.html.tpl" . }}

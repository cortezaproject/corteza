{{ template "inc_header.html.tpl" . }}
<div class="card-body p-0">
	<h1 class="h4 card-title p-3 border-bottom">{{ tr "password-reset-requested.template.title" }}</h1>
	<div class="p-3" role="alert">
	    {{ tr "password-reset-requested.template.instructions" }}
	</div>
    <div class="text-center my-3">
		{{ tr "password-reset-requested.template.links" "signup" links.Signup "login" links.Login }}</a>
    </div>
</div>
{{ template "inc_footer.html.tpl" . }}

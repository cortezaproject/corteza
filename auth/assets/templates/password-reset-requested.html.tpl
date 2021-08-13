{{ template "inc_header.html.tpl" . }}
<div class="card-body p-0">
	<h1 class="h4 card-title p-3 border-bottom">{{ tr "password_reset_requested.template.title" }}</h1>
	<div class="p-3" role="alert">
	    {{ tr "password_reset_requested.template.alert" }}
	</div>
    <div class="text-center my-3">
        <a href="{{ links.Signup }}">{{ tr "password_reset_requested.template.new-account-link" }}</a>
        or
        <a href="{{ links.Login }}">{{ tr "password_reset_requested.template.login" }}</a>
    </div>
</div>
{{ template "inc_footer.html.tpl" . }}

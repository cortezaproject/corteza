{{ template "inc_header.html.tpl" . }}
<div class="card-body p-0">
	<h1 class="h4 card-title p-3 border-bottom">{{ tr "pending_email_confirmation.template.title" }}</h1>
	<div class="p-3" role="alert">
		{{ tr "pending_email_confirmation.template.alert" }}
	</div>
    <div class="text-center my-3">
        <a href="{{ links.Signup }}">{{ tr "pending_email_confirmation.template.new-acc-link" }}</a>
        or
        <a href="{{ links.Login }}">{{ tr "pending_email_confirmation.template.login" }}</a>
    </div>
</div>
{{ template "inc_footer.html.tpl" . }}

{{ template "inc_header.html.tpl" . }}
<div class="card-body p-0">
	<h4 class="card-title p-3 border-bottom">{{ tr "pending-email-confirmation.template.title" }}</h4>
	<div class="p-3" role="alert">
		{{ tr "pending-email-confirmation.template.instructions" }}
	</div>
    <div class="text-center my-3">
		{{ tr "pending-email-confirmation.template.links" "signup" links.Signup "login" links.Login }}</a>
    </div>
</div>
{{ template "inc_footer.html.tpl" . }}

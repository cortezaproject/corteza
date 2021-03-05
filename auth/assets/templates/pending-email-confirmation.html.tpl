{{ template "inc_header.html.tpl" . }}
<div class="card-body p-0">
	{{ template "inc_alerts.html.tpl" .alerts }}
	<h4 class="card-title p-3 border-bottom">Confirm your email</h4>
	<div class="text-primary p-3" role="alert">
		You should receive email confirmation link to your inbox in a few moments.
	</div>
    <div class="text-center my-3">
        <a href="{{ links.Signup }}">Create new account</a>
        or
        <a href="{{ links.Login }}">Log in</a>
    </div>
</div>
{{ template "inc_footer.html.tpl" . }}

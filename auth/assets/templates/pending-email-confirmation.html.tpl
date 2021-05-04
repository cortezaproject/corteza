{{ template "inc_header.html.tpl" . }}
<div class="card-body p-0">
	<h1 class="h4 card-title p-3 border-bottom">Confirm your email</h1>
	<div class="p-3" role="alert">
		You should receive email confirmation link to your inbox in a few moments.
	</div>
    <div class="text-center my-3">
        <a href="{{ links.Signup }}">Create new account</a>
        or
        <a href="{{ links.Login }}">Log in</a>
    </div>
</div>
{{ template "inc_footer.html.tpl" . }}

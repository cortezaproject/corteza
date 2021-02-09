{{ template "inc_header.html.tpl" . }}
<div class="card-body">
	{{ template "inc_alerts.html.tpl" .alerts }}
	<h4 class="card-title">Password reset requested</h4>
	<div class="alert alert-primary" role="alert">
		If entered email is found in our database you should receive password reset link to your inbox in a few moments.
	</div>
	<div class="text-center my-3 small">
		<a href="{{ links.Signup }}">Create new account</a>
		|
		<a href="{{ links.Login }}">Login</a>
	</div>
</div>
{{ template "inc_footer.html.tpl" . }}

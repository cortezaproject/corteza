{{ template "inc_header.html.tpl" set . "activeNav" "security" }}
<div class="card-body p-0">
	<h4 class="card-title p-3 border-bottom">Security</h4>
	<ul>
		{{ if .settings.LocalEnabled }}
		<li><a href="{{ links.ChangePassword }}">Change your password</a></li>
		{{ end }}
	</ul>
</div>
{{ template "inc_footer.html.tpl" . }}

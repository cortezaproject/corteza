<!doctype html>
<html lang="en">
	<head>
		<meta charset="utf-8">
		<meta name="viewport" content="width=device-width, initial-scale=1">
		<!-- Bootstrap core CSS -->
		<link href="https://cdn.jsdelivr.net/npm/bootstrap@4.6.0/dist/css/bootstrap.min.css" rel="stylesheet">
		<script type="application/javascript">
			window.addEventListener('load', function() {
				document.getElementById('preview').src=location.hash.substring(1)
			})
		</script>
	</head>
	<body class="bg-dark text-light">
		<div class="container-fluid">
			<div class="row m-0 p-0">
				<div class="col-3 pl-0">
					<sidebar class="vh-100 position-fixed w-25 overflow-auto">
					{{ range .templates }}
                        <div class="pt-2 w-75">
                            <code>{{ .Template }}.html.tpl</code>
                            <ul class="nav flex-column">
                                {{ range .Scenes }}
                                <li class="nav-item border-top">
                                    <a target="preview"
                                        class="nav-link text-light"
                                        href="/auth/dev/scenarios?template={{ .Template }}&scene={{ .Name }}"
                                        onclick="location.hash=this.href"
                                    >
                                        <div>
                                            {{ .Name }}
                                        </div>
                                    </a>
                                </li>
                                {{ end }}
                            </ul>
                        </div>
					{{ end }}
                    <sidebar>
				</div>
				<div class="col-9 m-0 p-0">
					<iframe
						name="preview"
						id="preview"
						class="w-100 border-0"
						style="height: 100vh;"
						src="about:blank"
					/>
				</div>
			</div>
		</div>
	</body>
</html>

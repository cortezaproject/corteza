<!doctype html>
<html lang="{{ language }}">

<head>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<!-- Bootstrap core CSS -->
	<link href="https://cdn.jsdelivr.net/npm/bootstrap@4.6.0/dist/css/bootstrap.min.css" rel="stylesheet">
	<!-- Bootstrap icons -->
	<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bootstrap-icons@1.3.0/font/bootstrap-icons.css">
	<link rel="preconnect" href="https://fonts.googleapis.com">
	<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
	<link href="https://fonts.googleapis.com/css2?family=Poppins:wght@500;700&display=swap" rel="stylesheet">
	<!-- Custom CSS -->
	<link href="{{ links.AuthAssets }}/style.css?{{ buildtime }}" rel="stylesheet">
	<title>Corteza</title>
	<style>
		body {
			{{ safeCSS .authBg }}
			background-size: cover;
			background-attachment: fixed;
		}
	</style>
</head>

<body>
	<header>
		{{ if .user }}
		<div class="float-right text-white m-2">
			<a class="font-weight-bold text-white" href="{{ links.Base }}"><i
					class="bi bi-grid-3x2-gap-fill text-white mr-1 align-middle" style="font-size: 1.4rem;"></i></a>
			{{ tr "inc_header.logged-in-as" }}
			<a data-test-id="link-redirect-to-profile" class="font-weight-bold text-white"
				href="{{ links.Profile }}">{{ coalesce .user.Name .user.Handle .user.Email }}</a>
			|
			<a data-test-id="link-logout" class="font-weight-bold text-white"
				href="{{ links.Logout }}">{{ tr "inc_header.logout" }}</a>
		</div>
		{{ end }}
	</header>

	<main class="auth mt-sm-5">
		<div class="card">
{{ template "inc_nav.html.tpl" . }}
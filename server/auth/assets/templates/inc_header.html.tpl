<!doctype html>
<html data-color-mode="{{ .theme }}" lang="{{ language }}">

<head>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1">

	<!-- Bootstrap icons -->
	<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bootstrap-icons@1.3.0/font/bootstrap-icons.css">
	<link rel="preconnect" href="https://fonts.googleapis.com">
	<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>

	<!-- Fonts -->
	<link href="{{ links.AuthAssets }}/fonts.css" rel="stylesheet">

	<!-- Custom CSS -->
	<link href="/custom.css" rel="stylesheet">
	<link href="{{ links.AuthAssets }}/style.css?{{ buildtime }}" rel="stylesheet">

	<title>Corteza</title>
	<style>
		body {
			font-size: 1rem !important;
			{{ .authBg }}
			background-size: cover;
			background-attachment: fixed;
		}
	</style>
</head>

<body>
	<header>
		{{ if .user }}
		<div class="d-flex justify-content-end align-items-center text-white m-2">
			<a class="text-white font-weight-bold" href="{{ links.Base }}">
				{{ tr "inc_header.menu" }}
			</a>
			<span class="mx-2">
			|
			</span>
			{{ tr "inc_header.logged-in-as" }}
			<a data-test-id="link-redirect-to-profile" class="font-weight-bold text-white mx-2"
			 href="{{ links.Profile }}">{{ coalesce .user.Name .user.Handle .user.Email }}
			</a>
			|
			<a data-test-id="link-logout" class="font-weight-bold text-white ml-2" href="{{ links.Logout }}">
			 {{ tr "inc_header.logout" }}
			</a>
		</div>
		{{ end }}
	</header>

	<main class="auth mt-sm-5">
		<div class="tabs card">
{{ template "inc_nav.html.tpl" . }}

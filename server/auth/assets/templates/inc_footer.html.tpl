		</main>
		{{ template "inc_toasts.html.tpl" .alerts }}
		<footer class="d-flex align-items-end justify-content-center text-white py-4">
			{{ tr "inc_footer.code-link" }}
			<a data-test-id="link-github" href="https://github.com/crusttech/human" target="_blank" class="text-white ml-1">GitHub</a>
		</footer>
	</body>
	<script src="{{ links.AuthAssets }}/script.js?{{ buildtime }}"></script>
</html>

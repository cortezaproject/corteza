		</main>
		{{ template "inc_toasts.html.tpl" .alerts }}
		<footer class="d-flex align-items-end justify-content-center text-white p-2 pt-4">
			{{ tr "inc_footer.code-link" }}<a data-test-id="link-github" href="https://github.com/cortezaproject/" target="_blank" class="text-white ml-2">GitHub</a>
		</footer>
	</body>
	<script src="https://code.jquery.com/jquery-3.5.1.slim.min.js" integrity="sha384-DfXdz2htPH0lsSSs5nCTpuj/zy4C+OGpamoFVy38MVBnE+IbbVYUew+OrCXaRkfj" crossorigin="anonymous"></script>
	<script src="https://cdn.jsdelivr.net/npm/bootstrap@4.6.0/dist/js/bootstrap.bundle.min.js" integrity="sha384-Piv4xVNRyMGpqkS2by6br4gNJ7DXjqk09RmUpJ8jgGtD7zP9yug3goQfGII0yAns" crossorigin="anonymous"></script>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/jquery.mask/1.14.16/jquery.mask.js" integrity="sha512-0XDfGxFliYJPFrideYOoxdgNIvrwGTLnmK20xZbCAvPfLGQMzHUsaqZK8ZoH+luXGRxTrS46+Aq400nCnAT0/w==" crossorigin="anonymous"></script>
	<script src="{{ links.AuthAssets }}/script.js?{{ buildtime }}"></script>
</html>

		</main>
		{{ template "inc_toasts.html.tpl" .alerts }}
        {{ if .user }}
            <div class="position-absolute user text-white m-2">
                You're logged-in as
                <a class="font-weight-bold text-white" href="{{ links.Profile }}">{{ coalesce .user.Name .user.Handle .user.Email }}</a>
                |
                <a class="font-weight-bold text-white" href="{{ links.Logout }}">Logout</a>
            </div>
        {{ end }}
        <div class="footer col text-center position-absolute mb-4">
            <a href="https://cortezaproject.org/" target="_blank" class="text-white mr-2">cortezaproject.org</a>
            <a href="https://github.com/cortezaproject/" target="_blank" class="text-white ml-2">GitHub</a>
        </div>
		<i class="p-1 small text-white position-absolute version mr-3 mb-3">
			version {{ version }}
		</small>
	</body>
	<script src="https://code.jquery.com/jquery-3.5.1.slim.min.js" integrity="sha384-DfXdz2htPH0lsSSs5nCTpuj/zy4C+OGpamoFVy38MVBnE+IbbVYUew+OrCXaRkfj" crossorigin="anonymous"></script>
	<script src="https://cdn.jsdelivr.net/npm/bootstrap@4.6.0/dist/js/bootstrap.bundle.min.js" integrity="sha384-Piv4xVNRyMGpqkS2by6br4gNJ7DXjqk09RmUpJ8jgGtD7zP9yug3goQfGII0yAns" crossorigin="anonymous"></script>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/jquery.mask/1.14.16/jquery.mask.js" integrity="sha512-0XDfGxFliYJPFrideYOoxdgNIvrwGTLnmK20xZbCAvPfLGQMzHUsaqZK8ZoH+luXGRxTrS46+Aq400nCnAT0/w==" crossorigin="anonymous"></script>
	<script src="{{ links.Assets }}/script.js?{{ buildtime }}"></script>
</html>

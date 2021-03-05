            {{ if .user }}
                <div class="col cols-6 text-right small m-3">
                    Logged-in as
                    <a href="{{ links.Profile }}">{{ coalesce .user.Name .user.Handle .user.Email }}</a>
                    |
                    <a href="{{ links.Logout }}">Logout</a>
                </div>
            {{ end }}
		</main>
        <div class="footer col text-center position-absolute mb-4">
            <a href="https://cortezaproject.org/" target="_blank" class="text-white mr-2">cortezaproject.org</a>
            <a href="https://github.com/cortezaproject/" target="_blank" class="text-white ml-2">GitHub</a>
        </div>
		<small class="p-1 text-secondary position-absolute version">
			version {{ version }}
		</small>
	</body>
	<script src="{{ links.Assets }}/script.js?{{ buildtime }}"></script>
</html>

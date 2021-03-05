		</main>
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
		<small class="p-1 text-secondary position-absolute version mt-2">
			version {{ version }}
		</small>
	</body>
	<script src="{{ links.Assets }}/script.js?{{ buildtime }}"></script>
</html>

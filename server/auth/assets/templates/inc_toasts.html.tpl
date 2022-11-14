<div class="position-fixed bottom-0 right-0 p-3 mb-2" style="z-index: 5; right: 0; bottom: 0;">
{{ range . }}
	<div class="toast align-items-center bg-dark text-white" role="alert" aria-live="polite" aria-atomic="true" data-delay="5000" style="min-width:200px;">
	  <div class="d-flex toast-body border-top border-{{ .Type }}">
		{{ .Text | html }}
		<button type="button" class="ml-auto text-white btn p-0" data-dismiss="toast" aria-label="Close">
		<i class="bi bi-x"></i>
		</button>
	  </div>
	</div>
{{ end }}
</div>

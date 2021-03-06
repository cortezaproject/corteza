<div class="position-fixed bottom-0 right-0 p-3" style="z-index: 5; right: 0; bottom: 0;">
{{ range . }}
	<div class="toast" role="alert" aria-live="polite" aria-atomic="true" data-delay="5000">
	  <div class="toast-header">
		<span class="badge badge-{{ .Type }} mr-1 px-2">&nbsp;</span>
		<strong class="mr-auto">Corteza</strong>
		<button type="button" class="ml-2 mb-1 close" data-dismiss="toast" aria-label="Close">
		  <span aria-hidden="true">&times;</span>
		</button>
	  </div>
	  <div class="toast-body p-3">
		{{ .Text | html }}
	  </div>
	</div>
{{ end }}
</div>

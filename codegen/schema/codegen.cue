package schema

#codegen: {
	template: string
	output:   string

	syntax: string | *"go"
	if output =~ "\\.adoc$" {
		syntax: "adoc"
	}

	payload: _
}

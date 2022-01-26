package schema

#codegen: {
	template: string
	output:   string

	syntax: "go"
	if output =~ "\\.adoc" {
		syntax: "adoc"
	}

	payload: _
}

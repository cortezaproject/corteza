package schema

#_ioSpec: {
	template: string
	output:   string

	syntax: string | *"go"
	if output =~ "\\.adoc$" {
		syntax: "adoc"
	}
}

#codegen: {
	bulk?: [...#_ioSpec]
	if bulk == _|_ {
		#_ioSpec
	}

	payload: _
}

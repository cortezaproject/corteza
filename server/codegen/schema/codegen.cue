package schema

#_ioSpec: {
	template: string
	output:   string

	syntax: string | *"go"
	if output =~ "\\.md$" {
		syntax: "md"
	}
}

#codegen: {
	#_ioSpec
	payload: _
} | {
	bulk?: [...#_ioSpec]
	payload: _
}

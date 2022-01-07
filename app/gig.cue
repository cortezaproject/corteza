package app

// Decoder definitions

_decoders: [{
	ident:       "noop"
	description: "Noop does nothing."
	struct: {
		source: {
			goType:   "uint64"
			exprType: "String"
			required: true
		}
	}
},
	{
		ident:       "archive"
		description: "Extracts the contents of the archive into sepparate sources; extraction is not recursive."
		struct: {
			source: {
				goType:   "uint64"
				exprType: "String"
				required: true
			}
		}
	}]

// Postprocessor definitions

_postprocessors: [
	{ident: "noop", description:    "Noop does nothing."},
	{ident: "discard", description: "Discards the resulting sources."},
	{ident: "save", description:    "Saves the resulting sources to perminant storage; not implemented."},
	{
		ident:       "archive"
		transformer: "postprocessorArchiveTransformer"
		description: "Compresses the resulting sources into an archive."
		struct: {
			encoding: {
				goType:   "archive"
				exprType: "String"
				castFunc: "archiveFromParams"
				required: true
			}
			name: {}
		}
	},
]

// Preprocessor definitions

_noopPreprocessors: [
	{ident: "noop", description: "Noop does nothing."},
]

_attachmentPreprocessors: [
	{
		ident:       "attachmentRemove"
		description: "Removes the attachment."
		struct: {
			mimeType: {
				required: true
			}
		}
	},
	{
		ident:       "attachmentTransform"
		description: "Applies the specified transformations."
		struct: {
			width: {
				goType:   "int"
				exprType: "Number"
			}
			height: {
				goType:   "int"
				exprType: "Number"
			}
		}
	},
]

_envoyPreprocessors: [
	{
		ident:       "resourceRemove"
		transformer: "preprocessorResourceRemoveTransformer"
		description: "Removes the specified resource."
		struct: {
			resource: {
				required: true
			}
			identifier: {}
		}
	},
	{
		ident:       "resourceLoad"
		description: "Loads the specified resource from internal storage."
		struct: {
			resource: {
				required: true
			}
			id: {
				goType:   "uint64"
				exprType: "string"
			}
			handle: {
			}
			query: {
			}
		}
	},
	{
		ident:       "namespaceLoad"
		description: "Loads the namespace with a predefined set of sub-resources."
		struct: {
			id: {
				goType:   "uint64"
				exprType: "string"
			}
			handle: {
			}
		}
	},
]

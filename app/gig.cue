package app

// Decoder definitions

_decoders: [{
	ident:       "noop"
	description: "Noop does nothing."
	struct: {
		source: {
			goType:   "uint64"
			exprType: "String"
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
		ident:       "experimentalExport"
		description: "Loads the namespace along with some sub-resources (modules, pages, charts, ...)"
		struct: {
			id: {
				goType:   "uint64"
				exprType: "string"
			}
			handle: {
			}

			inclRBAC: {
				goType: "bool"
			}
			inclRoles: {
				goType: "[]string"
				castFunc: "cast.ToStringSlice"
			}
			exclRoles: {
				goType: "[]string"
				castFunc: "cast.ToStringSlice"
			}

			inclTranslations: {
				goType: "bool"
			}
			inclLanguage: {
				goType: "[]string"
				castFunc: "cast.ToStringSlice"
			}
			exclLanguage: {
				goType: "[]string"
				castFunc: "cast.ToStringSlice"
			}
		}
	},
]

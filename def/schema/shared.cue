package schema

// Resource definition identifier
#ident: =~"^[a-z][a-zA-Z0-9_]*$"

// Exported identifier
#expIdent: =~"^[A-Z][a-zA-Z0-9]*$"

// More liberal then identifier, allows underscores and dots
#handle: =~"^[A-Za-z][a-zA-Z0-9_\\-\\.]*[a-zA-Z0-9]+$"


// More liberal then identifier, allows underscores and dots
#baseHandle: =~"^[a-z][a-z0-9-]*[a-z0-9]+$"

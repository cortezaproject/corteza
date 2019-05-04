# QueryLanguage (ql) package

_This package was written mainly to assist with column 
conversion (wrapping non-physical columns into functions that 
extract values from a JSON value)_

Provides lexer and ast parser to convert a set of simple
instructions (see `ast_parser_test.go` for examples) that closely 
resemble SQL into syntax trees.

It provides handlers for identifiers and functions that can
validate and modify nodes and tree at parse time.

Tree can than be converted back to SQL or to structs that
assist Squirrel select builder.

## Pending improvements

### Lexer

 - Operator validation
 
### AST Parser

 - simplify / combine ASTNode vs ASTSet vs Columns
 - improve resilience and detect basic syntax errors
 - properly handling COUNT(DISTINCT ...) ident syntax
 - parsing complex expressions (eg: `year(created_at) > year(NOW()) - 2`)

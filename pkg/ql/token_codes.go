package ql

type (
	tokenCode int
)

const (
	CHAR_WHITELIST_WHITESPACE = " \n\t"
	CHAR_WHITELIST_OPERATORS  = "!+-/*=<>"
	CHAR_WHITELIST_QUOTES     = "'"
)

const (
	// Special tokens
	ILLEGAL tokenCode = iota
	EOF
	WS // 2
	IDENT
	LNULL
	LBOOL // 4
	LNUMBER
	LSTRING
	COMMA    // ,
	DOT      // .
	OPERATOR // + - / *
	PARENTHESIS_OPEN
	PARENTHESIS_CLOSE
	KEYWORD
)

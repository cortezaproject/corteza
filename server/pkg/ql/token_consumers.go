package ql

import (
	"bytes"
	"strings"
)

type (
	TokenConsumerWS       struct{}
	TokenConsumerIdent    struct{}
	TokenConsumerOperator struct{}
	TokenConsumerComma    struct{}
	TokenConsumerString   struct{}
	TokenConsumerNumber   struct{}
	TokenConsumerGeneric  struct {
		token     tokenCode
		whitelist string
		maxLength int
	}
)

func in(ch rune, wl string) bool {
	var w rune
	for _, w = range wl {
		if ch == w {
			return true
		}
	}
	return false
}

func (g TokenConsumerGeneric) Test(ch rune) bool {
	return in(ch, g.whitelist)
}

func (g TokenConsumerGeneric) Consume(s RuneReader) Token {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent whitespace character into the buffer.
	// Non-whitespace characters and EOF will cause the loop to exit.
	for {
		if g.maxLength > 0 && buf.Len() >= g.maxLength {
			// Length control
			break
		}

		if ch := s.read(); ch == eof {
			break
		} else if !g.Test(ch) {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return Token{code: g.token, literal: buf.String()}
}

func (i TokenConsumerIdent) Test(ch rune) bool {
	return isLetter(ch)
}

// Consumes the current rune and all contiguous ident runes.
func (TokenConsumerIdent) Consume(s RuneReader) Token {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent ident character into the buffer.
	// Non-ident characters and EOF will cause the loop to exit.
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isLetter(ch) && !isDigit(ch) && ch != '_' {
			s.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	lit := strings.ToUpper(buf.String())

	switch lit {
	case "NULL":
		return Token{code: LNULL}
	case "TRUE", "FALSE":
		return Token{code: LBOOL, literal: lit}
	case "IS", "LIKE", "NOT", "AND", "OR", "XOR", "IN", "BETWEEN":
		return Token{code: OPERATOR, literal: lit}
	case "DESC", "ASC", "INTERVAL":
		return Token{code: KEYWORD, literal: lit}
	}

	// Otherwise return as a regular identifier.
	return Token{code: IDENT, literal: buf.String()}
}

func (str TokenConsumerString) Test(ch rune) bool {
	return in(ch, CHAR_WHITELIST_QUOTES)
}

// Consumes entire string (skipping quotes)
func (str TokenConsumerString) Consume(s RuneReader) Token {
	var buf bytes.Buffer
	var escaping = false
	var ch = s.read() // skip quite

	for {
		if ch = s.read(); ch == eof {
			break
		} else if !escaping && str.Test(ch) { // test for quote
			return Token{code: LSTRING, literal: buf.String()}
		} else {
			escaping = !escaping && ch == '\\'
			if !escaping {
				// Add char to buffer if not escaping
				_, _ = buf.WriteRune(ch)
			}
		}
	}

	// This string did not end properly (with an enclosing quote).
	return Token{code: ILLEGAL, literal: buf.String() + string(ch)}
}

func (str TokenConsumerNumber) Test(ch rune) bool {
	return isDigit(ch)
}

// Consumes an entire number, fractional part included.
//
// The decimal point is the same rune that separates a reference path, so it is
// taken only when a digit follows it: "4.5" is one number, while the dot in
// "card.rarity" still belongs to the path and the dot in "4." still stands
// alone. Without the lookahead a fractional literal lexed as number, DOT,
// number and the parser refused the whole filter — so no money or score field
// could be compared against anything but a whole number.
func (str TokenConsumerNumber) Consume(s RuneReader) Token {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	seenPoint := false

	for {
		switch ch := s.peek(); {
		case isDigit(ch):
			_, _ = buf.WriteRune(s.read())

		case ch == '.' && !seenPoint && isDigit(s.peekAt(1)):
			seenPoint = true
			_, _ = buf.WriteRune(s.read())

		default:
			return Token{code: LNUMBER, literal: buf.String()}
		}
	}
}

// isLetter returns true if the rune is a letter.
func isLetter(ch rune) bool { return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') }

// isDigit returns true if the rune is a digit.
func isDigit(ch rune) bool { return ch >= '0' && ch <= '9' }

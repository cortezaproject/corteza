package ql

import (
	"bufio"
	"io"
)

// Lexer represents a lexical scanner.
type (
	RuneReader interface {
		read() rune
		unread()
	}

	TokenConsumers interface {
		Test(ch rune) bool
		Consume(s RuneReader) Token
	}

	Token struct {
		code    tokenCode
		literal string
		line    uint
		char    uint
	}

	Lexer struct {
		r         *bufio.Reader
		consumers []TokenConsumers
		line      uint
		char      uint
	}
)

// eof represents a marker rune for the end of the reader.
var eof = rune(0)

// NewLexer returns a new instance of Lexer.
func NewLexer(r io.Reader) *Lexer {
	return &Lexer{
		r: bufio.NewReader(r),
		consumers: []TokenConsumers{
			&TokenConsumerGeneric{token: WS, whitelist: CHAR_WHITELIST_WHITESPACE},
			// @todo ensure operator order (eg != is valid, =! is not)
			&TokenConsumerGeneric{token: OPERATOR, whitelist: CHAR_WHITELIST_OPERATORS},
			&TokenConsumerGeneric{token: COMMA, whitelist: ",", maxLength: 1},
			&TokenConsumerGeneric{token: PARENTHESIS_OPEN, whitelist: "(", maxLength: 1},
			&TokenConsumerGeneric{token: PARENTHESIS_CLOSE, whitelist: ")", maxLength: 1},
			&TokenConsumerString{},
			&TokenConsumerNumber{},
			&TokenConsumerIdent{},
		},
	}
}

// Scan returns the next token and literal value.
func (s *Lexer) Scan() Token {
	var ch = s.peek()

	if ch == '\n' {
		s.line++
		s.char = 0
	}

	if eof == ch {
		return Token{code: EOF, line: s.line, char: s.char}
	}

	for _, c := range s.consumers {
		if c.Test(ch) {
			t := c.Consume(s)
			t.line = s.line
			t.char = s.char
			return t
		}
	}

	return Token{code: ILLEGAL, literal: string(ch), line: s.line, char: s.char}
}

// read reads the next rune from the buffered reader.
// Returns the rune(0) if an error occurs (or io.EOF is returned).
func (s *Lexer) read() rune {
	ch, _, err := s.r.ReadRune()
	s.char++
	if err != nil {
		return eof
	}
	return ch
}

func (s *Lexer) peek() rune {
	bb, err := s.r.Peek(1)
	if err != nil || len(bb) == 0 {
		return eof
	}
	return rune(bb[0])
}

// unread places the previously read rune back on the reader.
func (s *Lexer) unread() { _ = s.r.UnreadRune() }

func (t Token) Is(cc ...tokenCode) bool {
	for _, c := range cc {
		if t.code == c {
			return true
		}
	}

	return false
}

package ql

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// A fractional literal used to lex as number, DOT, number, and the parser
// refused the whole filter — so no money or score field could be compared
// against anything but a whole number. The dot is the same rune that separates
// a reference path, so these pin both readings of it.
func TestLexerDecimalLiterals(t *testing.T) {
	lex := func(in string) []Token {
		l := NewLexer(strings.NewReader(in))
		var out []Token
		for {
			tk := l.Scan()
			if tk.code == EOF {
				return out
			}
			if tk.code != WS {
				out = append(out, tk)
			}
		}
	}

	t.Run("a fractional literal is one number", func(t *testing.T) {
		tt := lex("4.5")
		require.Len(t, tt, 1)
		require.Equal(t, LNUMBER, tt[0].code)
		require.Equal(t, "4.5", tt[0].literal)
	})

	t.Run("it survives a comparison, which is where it is used", func(t *testing.T) {
		tt := lex("overall_score > 4.5")
		require.Len(t, tt, 3)
		require.Equal(t, "overall_score", tt[0].literal)
		require.Equal(t, LNUMBER, tt[2].code)
		require.Equal(t, "4.5", tt[2].literal)
	})

	t.Run("a reference path is still a path", func(t *testing.T) {
		tt := lex("card.rarity")
		require.Len(t, tt, 3)
		require.Equal(t, IDENT, tt[0].code)
		require.Equal(t, DOT, tt[1].code)
		require.Equal(t, IDENT, tt[2].code)
	})

	t.Run("a trailing dot is not swallowed into the number", func(t *testing.T) {
		tt := lex("4.")
		require.Len(t, tt, 2)
		require.Equal(t, LNUMBER, tt[0].code)
		require.Equal(t, "4", tt[0].literal)
		require.Equal(t, DOT, tt[1].code)
	})

	t.Run("only the first dot joins the number", func(t *testing.T) {
		tt := lex("1.2.3")
		require.Equal(t, "1.2", tt[0].literal)
		require.Equal(t, DOT, tt[1].code)
	})

	t.Run("a whole number is unchanged", func(t *testing.T) {
		tt := lex("42")
		require.Len(t, tt, 1)
		require.Equal(t, "42", tt[0].literal)
	})
}

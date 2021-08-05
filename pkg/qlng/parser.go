package qlng

import (
	"fmt"
	"strings"
)

type (
	// Parser represents a parser.
	Parser struct {
		lexer  *Lexer
		tokbuf []Token

		OnIdent    IdentHandler
		OnFunction FunctionHandler

		// parenthesis level control
		level uint
	}

	IdentHandler    func(ident Ident) (Ident, error)
	FunctionHandler func(ident function) (parserNode, error)
)

// NewParser returns a new instance of Parser.
func NewParser() *Parser {
	p := &Parser{
		OnIdent:    func(ident Ident) (Ident, error) { return ident, nil },
		OnFunction: func(ident function) (parserNode, error) { return ident, nil },
	}

	return p
}

// Removes oldest token in the buffer, adds new one and returns 2nd oldest
func (p *Parser) nextToken() Token {
	var t Token
	for {
		t = p.lexer.Scan()
		if !t.Is(WS) {
			p.tokbuf = append(p.tokbuf[1:], t)
			return p.tokbuf[0]
		}
	}
}

func (p *Parser) peekToken(s int) Token {
	return p.tokbuf[s]
}

func (p *Parser) initLexer(s string) {
	p.lexer = NewLexer(strings.NewReader(s))
	p.tokbuf = make([]Token, 3)

	for c := 1; c < cap(p.tokbuf); c++ {
		// Fill the buffer
		p.nextToken()
	}
}

// Parse parses the given expression and returns the generated AST
func (p *Parser) Parse(s string) (*ASTNode, error) {
	p.initLexer(s)

	if set, err := p.parse(p.nextToken()); err != nil {
		return nil, err
	} else if len(set) == 1 {
		return set[0].ToAST(), set[0].Validate()
	} else {
		return set.ToAST(), set.Validate()
	}

}

// Peek ahead if there is an alias ident (<IDENT:AS>)
//
// @todo remove this; not used anywhere after report builder gets removed
func (p *Parser) peekIfAlias() bool {
	var f, s = p.peekToken(1), p.peekToken(2)
	return f.Is(IDENT) && strings.ToUpper(f.literal) == "AS" && s.Is(IDENT)
}

func (p *Parser) parse(t Token) (list parserNodes, err error) {
	goto checkToken

next:
	if p.peekToken(1).Is(COMMA) || p.peekIfAlias() {
		// Peek ahead and exit on comma
		return
	}
	if p.peekToken(1).Is(PARENTHESIS_CLOSE) {
		// Peek ahead and exit on closed parenthesis
		if p.level == 0 {
			return nil, fmt.Errorf("closing unopened parenthesis in expression")
		}
		p.level--
		return
	}

	t = p.nextToken()

checkToken:
	switch t.code {
	case EOF:
		break
	case WS:
		// Ignore ws... next token...
		goto next
	case ILLEGAL:
		return nil, fmt.Errorf("found an illegal token (%+v)", t)
	case IDENT:
		var ident parserNode
		if ident, err = p.parseIdent(t); err != nil {
			return nil, err
		} else {
			list = append(list, ident)
			goto next
		}
	case LNULL:
		list = append(list, lNull{})
		goto next
	case LBOOL:
		list = append(list, lBoolean{value: evalBool(t.literal)})
		goto next
	case OPERATOR:
		if len(list) > 0 {
			// Merge with previous operator node
			if prevOp, ok := list[len(list)-1].(operator); ok {
				list[len(list)-1] = operator{kind: prevOp.kind + " " + t.literal}
				goto next
			}
		}

		list = append(list, operator{kind: t.literal})
		goto next
	case KEYWORD:
		if keyword, err := p.parseKeyword(t); err != nil {
			return nil, err
		} else {
			list = append(list, keyword)
		}
		goto next
	case LNUMBER:
		list = append(list, lNumber{value: t.literal})
		goto next
	case LSTRING:
		list = append(list, lString{value: t.literal})
		goto next
	case PARENTHESIS_OPEN:
		depth := p.level
		p.level++
		if sub, err := p.parse(p.nextToken()); err != nil {
			return nil, err
		} else {
			list = append(list, sub)
		}

		// Allow parent level to continue parsing.
		// Example: ((A) AND (B))
		// +1 since PARENTHESIS_CLOSE decrease level
		if (p.level + 1) != depth {
			p.nextToken()
			goto next
		}
	default:
		return nil, fmt.Errorf("unexpected token while parsing expression (%v)", t)
	}

	return list, nil
}

func (p *Parser) parseIdent(t Token) (list parserNode, err error) {
	if p.peekToken(1).Is(PARENTHESIS_OPEN) {
		// Handle function calls: <IDENT><PARENTHESIS_OPEN>...
		f := function{name: t.literal}
		if f.arguments, err = p.parseSet(); err != nil {
			return nil, err
		} else {
			return p.OnFunction(f)
		}
	}

	i := Ident{Value: t.literal}

	if p.peekToken(1).Is(DOT) {
		p.nextToken()
		i.Value += "."
		l2 := p.nextToken()
		if l2.Is(IDENT) {
			i.Value += l2.literal
		}
	}

	return p.OnIdent(i)
}

func (p *Parser) parseSet() (list parserNodeSet, err error) {
	var expr parserNodes
	var parenthesisOpened = false

next:
	t := p.nextToken()

	if p.peekIfAlias() {
		return
	}

	switch t.code {
	case WS:
		goto next
	case PARENTHESIS_OPEN:
		p.level++
		parenthesisOpened = true
		goto next
	case LNUMBER:
		list = append(list, lNumber{value: t.literal})
		goto next
	case LSTRING:
		list = append(list, lString{value: t.literal})
		goto next
	case OPERATOR:
		list = append(list, operator{kind: t.literal})
		goto next
	case KEYWORD:
		if keyword, err := p.parseKeyword(t); err != nil {
			return nil, err
		} else {
			list = append(list, keyword)
		}
		goto next
	case IDENT:
		if p.peekToken(1).Is(OPERATOR) {
			// Looks like we have an expression ahead of us

			if parenthesisOpened {
				// Expression will find closing parenthesis and dec. the level
				// so, lets bump up the number
				p.level++
			}

			if expr, err = p.parse(t); err != nil {
				return
			}

			list = append(list, expr)
			goto next
		}

		var ident parserNode
		if ident, err = p.parseIdent(t); err != nil {
			return nil, err
		} else {
			list = append(list, ident)

			goto next
		}
	case COMMA:
		goto next
	case EOF:
		return
	case PARENTHESIS_CLOSE:
		// Peek ahead and exit on closed parenthesis
		if p.level == 0 {
			return nil, fmt.Errorf("closing unopened parenthesis in set")
		}
		p.level--
		return
	default:
		return nil, fmt.Errorf("unexpected token while parsing set (%v)", t)
	}
}

func (p *Parser) parseKeyword(t Token) (list parserNode, err error) {
	switch strings.ToUpper(t.literal) {
	case "INTERVAL":
		i := interval{value: p.nextToken().literal}
		u := p.nextToken()

		if u.code != IDENT {
			return nil, fmt.Errorf("expecting identifier, got %v", t)
		} else {
			switch strings.ToUpper(u.literal) {
			case "MICROSECOND", "SECOND", "MINUTE", "HOUR",
				"DAY", "WEEK", "MONTH", "QUARTER", "YEAR",
				"SECOND_MICROSECOND", "MINUTE_MICROSECOND", "MINUTE_SECOND", "HOUR_MICROSECOND", "HOUR_SECOND",
				"HOUR_MINUTE", "DAY_MICROSECOND", "DAY_SECOND", "DAY_MINUTE", "DAY_HOUR", "YEAR_MONTH":
				// All good
				break
			default:
				return nil, fmt.Errorf("expecting interval unit, got %v", u.literal)
			}
		}

		i.unit = u.literal

		return i, nil

	default:
		return keyword{keyword: t.literal}, nil
	}
}

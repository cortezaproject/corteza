package ql

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
	FunctionHandler func(ident Function) (Function, error)
)

// NewParser returns a new instance of Parser.
func NewParser() *Parser {
	p := &Parser{
		OnIdent:    func(ident Ident) (Ident, error) { return ident, nil },
		OnFunction: func(ident Function) (Function, error) { return ident, nil },
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

func (p *Parser) ParseSet(s string) (ASTNode, error) {
	p.initLexer(s)

	if set, err := p.parseSet(); err != nil {
		return nil, err
	} else {
		return set, set.Validate()
	}
}

func (p *Parser) ParseExpression(s string) (ASTNode, error) {
	p.initLexer(s)

	if set, err := p.parseExpr(p.nextToken()); err != nil {
		return nil, err
	} else if len(set) == 1 {
		return set[0], set[0].Validate()
	} else {
		return set, set.Validate()
	}

}

func (p *Parser) ParseColumns(s string) (columns Columns, err error) {
	p.initLexer(s)

	var t Token
	var c Column

next:
	t = p.nextToken()
	switch t.code {
	case COMMA:
		goto next
	case EOF:
		break
	case ILLEGAL:
		return nil, fmt.Errorf("found an illegal token (%+v)", t)
	default:
		if c, err = p.parseColumn(t); err != nil {
			return nil, err
		} else {
			columns = append(columns, c)
		}
		goto next
	}

	err = columns.Validate()
	return
}

func (p *Parser) parseColumn(t Token) (c Column, err error) {
	if c.Expr, err = p.parseExpr(t); err != nil {
		return
	}

	// Set alias move forward for 2 places
	if p.peekIfAlias() {
		c.Alias = p.peekToken(2).literal
		p.nextToken()
		p.nextToken()
	}

	return
}

// Peek ahead if there is an alias ident (<IDENT:AS>
func (p *Parser) peekIfAlias() bool {
	var f, s = p.peekToken(1), p.peekToken(2)
	return f.Is(IDENT) && strings.ToUpper(f.literal) == "AS" && s.Is(IDENT)
}

func (p *Parser) parseExpr(t Token) (list ASTNodes, err error) {
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
		var ident ASTNode
		if ident, err = p.parseIdent(t); err != nil {
			return nil, err
		} else {
			list = append(list, ident)
			goto next
		}
	case OPERATOR:
		if len(list) > 0 {
			// Merge with previous operator node
			if prevOp, ok := list[len(list)-1].(Operator); ok {
				list[len(list)-1] = Operator{Kind: prevOp.Kind + " " + t.literal}
				goto next
			}
		}

		list = append(list, Operator{Kind: t.literal})
		goto next
	case KEYWORD:
		if keyword, err := p.parseKeyword(t); err != nil {
			return nil, err
		} else {
			list = append(list, keyword)
		}
		goto next
	case NUMBER:
		list = append(list, Number{Value: t.literal})
		goto next
	case STRING:
		list = append(list, String{Value: t.literal})
		goto next
	case PARENTHESIS_OPEN:
		p.level++
		if sub, err := p.parseExpr(p.nextToken()); err != nil {
			return nil, err
		} else {
			list = append(list, sub)
		}
	default:
		return nil, fmt.Errorf("unexpected token while parsing expression (%v)", t)
	}

	return list, nil
}

func (p *Parser) parseIdent(t Token) (list ASTNode, err error) {
	if p.peekToken(1).Is(PARENTHESIS_OPEN) {
		// Handle function calls: <IDENT><PARENTHESIS_OPEN>...
		f := Function{Name: t.literal}
		if f.Arguments, err = p.parseSet(); err != nil {
			return nil, err
		} else {
			return p.OnFunction(f)
		}
	}

	return p.OnIdent(Ident{Value: t.literal})
}

func (p *Parser) parseSet() (list ASTSet, err error) {
	var expr ASTNodes
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
	case NUMBER:
		list = append(list, Number{Value: t.literal})
		goto next
	case STRING:
		list = append(list, String{Value: t.literal})
		goto next
	case OPERATOR:
		list = append(list, Operator{Kind: t.literal})
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

			if expr, err = p.parseExpr(t); err != nil {
				return
			}

			list = append(list, expr)
			goto next
		}

		var ident ASTNode
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

func (p *Parser) parseKeyword(t Token) (list ASTNode, err error) {
	switch strings.ToUpper(t.literal) {
	case "INTERVAL":
		i := Interval{Value: p.nextToken().literal}
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

		i.Unit = u.literal

		return i, nil

	default:
		return Keyword{Keyword: t.literal}, nil
	}
}

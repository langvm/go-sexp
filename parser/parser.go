// Copyright 2024 LangVM Project
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package parser

import (
	scanner "github.com/langvm/cee-scanner"
)

type Parser struct {
	scanner.Scanner

	ReachedEOF bool

	Token Token
}

func NewParser() Parser {
	p := Parser{
		Scanner: scanner.Scanner{
			Whitespaces: map[rune]int{
				' ':  1,
				'\t': 1,
				'\r': 1,
				'\n': 1,
			},
			Delimiters: map[rune]int{
				'(': 1,
				')': 1,
			},
		},
	}
	return p
}

func (p *Parser) Scan() error {
	begin, kind, format, litRunes, err := p.Scanner.ScanToken()
	switch err := err.(type) {
	case nil:
	case scanner.EOFError:
		if p.ReachedEOF {
			return err
		} else {
			p.ReachedEOF = true
			return nil
		}
	default:
		return err
	}

	lit := string(litRunes)

	switch kind {
	case scanner.MARK:
		switch lit {
		case "(":
			kind = LPAREN
		case ")":
			kind = RPAREN
		default:
			panic("impossible")
		}
	case scanner.IDENT:
		kind = IDENT
	case scanner.CHAR:
		kind = CHAR
	case scanner.STRING:
		kind = STRING
	case scanner.INT:
		kind = INT
	case scanner.FLOAT:
		kind = FLOAT
	case scanner.COMMENT:
		return p.Scan()
	default:
		panic("impossible")
	}

	p.Token = Token{
		PosRange: PosRange{begin, p.Position},
		Kind:     kind,
		Format:   format,
		Literal:  lit,
	}

	return nil
}

func (p *Parser) MatchTerm(term int) error {
	if p.Token.Kind != term {
		return UnexpectedToken{p.Token}
	}

	err := p.Scan()
	if err != nil {
		return err
	}

	return nil
}

func (p *Parser) ExpectLiteralValue() (LiteralValue, error) {
	t := p.Token

	err := p.Scan()
	if err != nil {
		return LiteralValue{}, err
	}

	return LiteralValue{t}, nil
}

func (p *Parser) ExpectIdent() (Ident, error) {
	if p.Token.Kind != IDENT {
		return Ident{}, UnexpectedToken{p.Token}
	}

	t := p.Token

	err := p.Scan()
	if err != nil {
		return Ident{}, err
	}

	return Ident{t}, nil
}

func (p *Parser) ExpectList() (List, error) {
	err := p.MatchTerm(LPAREN)
	if err != nil {
		return List{}, err
	}

	ident, err := p.ExpectIdent()
	if err != nil {
		return List{}, err
	}

	var elements []Node

	for p.Token.Kind != RPAREN {
		node, err := p.ExpectNode()
		if err != nil {
			return List{}, err
		}
		elements = append(elements, node)
	}

	err = p.MatchTerm(RPAREN)
	if err != nil {
		return List{}, err
	}

	return List{
		Prefix:   ident,
		Elements: elements,
	}, nil
}

func (p *Parser) ExpectNode() (Node, error) {
	switch p.Token.Kind {
	case LPAREN:
		return p.ExpectList()
	case IDENT:
		return p.ExpectIdent()
	case INT:
		fallthrough
	case FLOAT:
		fallthrough
	case STRING:
		fallthrough
	case CHAR:
		return p.ExpectLiteralValue()
	default:
		return nil, UnexpectedToken{p.Token}
	}
}

func Parse(buf []rune) (List, error) {
	p := Parser{
		Scanner: scanner.Scanner{
			BufferScanner: scanner.BufferScanner{
				Buffer: buf,
			},
			Whitespaces: map[rune]int{
				' ':  1,
				'\n': 1,
				'\t': 1,
				'\r': 1,
			},
			Delimiters: map[rune]int{
				'(': 1,
				')': 1,
			},
		},
	}
	err := p.Scan()
	if err != nil {
		return List{}, err
	}

	return p.ExpectList()
}
